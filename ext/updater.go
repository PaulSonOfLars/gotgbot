package ext

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var ErrMissingCertOrKeyFile = errors.New("missing certfile or keyfile")
var ErrExpectedEmptyServer = errors.New("expected server to be nil")

// botData is an internal struct that is used by the updater to keep track of the necessary update channels for each bot.
type botData struct {
	bot        *gotgbot.Bot
	updateChan chan json.RawMessage
	urlPath    string
}

type Updater struct {
	// Dispatcher defines how to handle incoming updates.
	Dispatcher *Dispatcher
	// ErrorLog defines how to log errors which occur in the updater.
	ErrorLog *log.Logger

	stopIdling chan bool
	running    chan bool
	server     *http.Server
	// map tokens to channels
	botMapping map[string]botData
	serveMux   *http.ServeMux
}

var errorLog = log.New(os.Stderr, "ERROR", log.LstdFlags)

type UpdaterOpts struct {
	ErrorLog *log.Logger
	// The dispatcher instance to be used by the updater.
	Dispatcher *Dispatcher
}

// NewUpdater Creates a new Updater, as well as the necessary structures required for the associated Dispatcher.
func NewUpdater(opts *UpdaterOpts) *Updater {
	errLog := errorLog
	// Default dispatcher
	dispatcher := NewDispatcher(nil)

	if opts != nil {
		if opts.ErrorLog != nil {
			errLog = opts.ErrorLog
		}
		if opts.Dispatcher != nil {
			dispatcher = opts.Dispatcher
		}
	}

	return &Updater{
		ErrorLog:   errLog,
		Dispatcher: dispatcher,
	}
}

// PollingOpts represents the optional values to start long polling.
type PollingOpts struct {
	// DropPendingUpdates decides  whether to drop "pending" updates; these are updates which were sent before
	// the bot was started.
	DropPendingUpdates bool
	// GetUpdatesOpts represents the opts passed to GetUpdates.
	// Note: It is recommended you edit the values here when running in production environments.
	// Changes might include:
	//    - Changing the "GetUpdatesOpts.AllowedUpdates" to only refer to relevant updates
	//    - Using a non-0 "GetUpdatesOpts.Timeout" value. This is how "long" telegram will hold the long-polling call
	//    while waiting for new messages. A value of 0 causes telegram to reply immediately, which will then cause
	//    your bot to immediately ask for more updates. While this can seem fine, it will eventually causing
	//    telegram to delay your requests when left running over longer periods. If you are seeing lots
	//    of "context deadline exceeded" errors on GetUpdates, this is likely the cause.
	//    Keep in mind that a timeout of 10 does not mean you only get updates every 10s; by the nature of
	//    long-polling, Telegram responds to your request as soon as new messages are available.
	//    When setting this, it is recommended you set your PollingOpts.Timeout value to be slightly bigger (eg, +1).
	GetUpdatesOpts gotgbot.GetUpdatesOpts
}

// StartPolling starts polling updates from telegram using getUpdates long-polling.
// See PollingOpts for optional values to set in production environments.
func (u *Updater) StartPolling(b *gotgbot.Bot, opts *PollingOpts) error {
	if u.botMapping == nil {
		u.botMapping = make(map[string]botData)
	}

	// This logic is currently mostly duplicated over from the generated getUpdates code.
	// This is a performance improvement to avoid:
	// - needing to re-allocate new url.values structs.
	// - needing to convert the 'opt' values to strings.
	// - unnecessary unmarshalling of multiple full Update structs.
	// Yes, this also makes me sad. :/
	v := map[string]string{}
	dropPendingUpdates := false
	var reqOpts *gotgbot.RequestOpts

	if opts != nil {
		dropPendingUpdates = opts.DropPendingUpdates
		if opts.GetUpdatesOpts.RequestOpts != nil {
			reqOpts = opts.GetUpdatesOpts.RequestOpts
		}

		v["offset"] = strconv.FormatInt(opts.GetUpdatesOpts.Offset, 10)
		v["limit"] = strconv.FormatInt(opts.GetUpdatesOpts.Limit, 10)
		v["timeout"] = strconv.FormatInt(opts.GetUpdatesOpts.Timeout, 10)
		if opts.GetUpdatesOpts.AllowedUpdates != nil {
			bs, err := json.Marshal(opts.GetUpdatesOpts.AllowedUpdates)
			if err != nil {
				return fmt.Errorf("failed to marshal field allowed_updates: %w", err)
			}
			v["allowed_updates"] = string(bs)
		}
	}

	updateChan := make(chan json.RawMessage)
	u.botMapping[b.GetToken()] = botData{
		bot:        b,
		updateChan: updateChan,
	}

	go u.Dispatcher.Start(b, updateChan)
	go u.pollingLoop(b, reqOpts, updateChan, dropPendingUpdates, v)

	return nil
}

func (u *Updater) pollingLoop(b *gotgbot.Bot, opts *gotgbot.RequestOpts, updateChan chan json.RawMessage, dropPendingUpdates bool, v map[string]string) {
	// if dropPendingUpdates, force the offset to -1
	if dropPendingUpdates {
		v["offset"] = "-1"
	}

	var offset int64

	u.running = make(chan bool)
	for {
		select {
		case <-u.running:
			// if anything comes in, stop.
			return
		default:
			// continue as usual
		}

		r, err := b.Request("getUpdates", v, nil, opts)
		if err != nil {
			u.ErrorLog.Println("failed to get updates; sleeping 1s: " + err.Error())
			time.Sleep(time.Second)
			continue

		} else if r == nil {
			dropPendingUpdates = false
			continue
		}

		var rawUpdates []json.RawMessage
		if err := json.Unmarshal(r, &rawUpdates); err != nil {
			u.ErrorLog.Println("failed to unmarshal updates: " + err.Error())
			continue
		}

		if len(rawUpdates) == 0 {
			dropPendingUpdates = false
			continue
		}

		var lastUpdate struct {
			UpdateId int64 `json:"update_id"`
		}

		if err := json.Unmarshal(rawUpdates[len(rawUpdates)-1], &lastUpdate); err != nil {
			u.ErrorLog.Println("failed to unmarshal last update: " + err.Error())
			continue
		}

		offset = lastUpdate.UpdateId + 1
		v["offset"] = strconv.FormatInt(offset, 10)
		if dropPendingUpdates {
			// Setting the offset to -1 gets just the last update; this should be skipped too.
			dropPendingUpdates = false
			continue
		}

		for _, updData := range rawUpdates {
			temp := updData // use new mem address to avoid loop conflicts
			updateChan <- temp
		}
	}
}

// Idle starts an infinite loop to avoid the program exciting while the background threads handle updates.
func (u *Updater) Idle() {
	// Create the idling channel
	u.stopIdling = make(chan bool)

	// Wait until some input is received from the idle channel, which will stop the idling.
	<-u.stopIdling
}

// Stop stops the current updater and dispatcher instances.
func (u *Updater) Stop() error {
	// if server, this is running on webhooks; shutdown the server
	if u.server != nil {
		err := u.server.Shutdown(context.Background())
		if err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
	}

	if u.running != nil {
		// stop the polling loop
		u.running <- false
		close(u.running)
	}

	// Close all the update channels
	for _, data := range u.botMapping {
		close(data.updateChan)
	}

	u.Dispatcher.Stop()

	if u.stopIdling != nil {
		// stop idling
		u.stopIdling <- false
		close(u.stopIdling)
	}
	return nil
}

// StartWebhook starts the webhook server for a single bot instance.
// This does NOT set the webhook on telegram - this should be done by the caller.
// The opts parameter allows for specifying various webhook settings.
func (u *Updater) StartWebhook(b *gotgbot.Bot, urlPath string, opts WebhookOpts) error {
	if u.server != nil {
		return ErrExpectedEmptyServer
	}

	u.AddWebhook(b, urlPath, opts)
	return u.StartServer(opts)
}

// AddWebhook prepares the webhook server to receive webhook updates for one bot, on a specific path.
func (u *Updater) AddWebhook(b *gotgbot.Bot, urlPath string, opts WebhookOpts) {
	if u.serveMux == nil {
		u.serveMux = http.NewServeMux()
	}
	if u.botMapping == nil {
		u.botMapping = make(map[string]botData)
	}

	updateChan := make(chan json.RawMessage)
	u.serveMux.HandleFunc("/"+urlPath, func(w http.ResponseWriter, r *http.Request) {
		if opts.SecretToken != "" && opts.SecretToken != r.Header.Get("X-Telegram-Bot-Api-Secret-Token") {
			// Drop any updates from invalid secret tokens.
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		bytes, _ := io.ReadAll(r.Body)
		updateChan <- bytes
	})

	u.botMapping[b.GetToken()] = botData{
		bot:        b,
		updateChan: updateChan,
		urlPath:    urlPath,
	}

	// Webhook has been added; relevant dispatcher should also be started.
	go u.Dispatcher.Start(b, updateChan)
}

// SetAllBotWebhooks sets all the webhooks for the bots that have been added to this updater via AddWebhook.
func (u *Updater) SetAllBotWebhooks(domain string, opts *gotgbot.SetWebhookOpts) error {
	for _, data := range u.botMapping {
		_, err := data.bot.SetWebhook(fmt.Sprintf("%s/%s", strings.TrimSuffix(domain, "/"), data.urlPath), opts)
		if err != nil {
			// Extract the botID, so we don't intentionally log the token
			botId := strings.Split(data.bot.GetToken(), ":")[0]
			return fmt.Errorf("failed to set webhook for %s: %w", botId, err)
		}
	}
	return nil
}

// StartServer starts the webhook server for all the bots added via AddWebhook.
// We recommend calling this BEFORE setting individual webhooks.
// The opts parameter allows for specifying TLS settings.
func (u *Updater) StartServer(opts WebhookOpts) error {
	if u.serveMux == nil {
		u.serveMux = http.NewServeMux()
	}

	var tls bool
	if opts.CertFile == "" && opts.KeyFile == "" {
		tls = false
	} else if opts.CertFile != "" && opts.KeyFile != "" {
		tls = true
	} else {
		return ErrMissingCertOrKeyFile
	}

	u.server = &http.Server{
		Addr:              opts.GetListenAddr(),
		Handler:           u.serveMux,
		ReadTimeout:       opts.ReadTimeout,
		ReadHeaderTimeout: opts.ReadHeaderTimeout,
	}

	go func() {
		var err error
		if tls {
			err = u.server.ListenAndServeTLS(opts.CertFile, opts.KeyFile)
		} else {
			err = u.server.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic("http server failed: " + err.Error())
		}
	}()

	return nil
}
