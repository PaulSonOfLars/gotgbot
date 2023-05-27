package ext

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	ErrMissingCertOrKeyFile = errors.New("missing certfile or keyfile")
	ErrExpectedEmptyServer  = errors.New("expected server to be nil")
)

// botData is an internal struct used by the updater to keep track of the necessary update channels for each bot.
type botData struct {
	// bot represents the bot for which this data is relevant.
	bot *gotgbot.Bot
	// updateChan represents the incoming updates channel.
	updateChan chan json.RawMessage
	// polling allows us to close the polling loop.
	polling chan struct{}
	// urlPath defines the incoming webhook URL path for this bot.
	urlPath string
}

type ErrorFunc func(error)

type Updater struct {
	// Dispatcher is where all the incoming updates are sent to be processed.
	Dispatcher *Dispatcher

	// UnhandledErrFunc provides more flexibility for dealing with previously unhandled errors, such as failures to get
	// updates (when long-polling), or failures to unmarshal.
	// If nil, the error goes to ErrorLog.
	UnhandledErrFunc ErrorFunc
	// ErrorLog specifies an optional logger for unexpected behavior from handlers.
	// If nil, logging is done via the log package's standard logger.
	ErrorLog *log.Logger

	// stopIdling is the channel that blocks the main thread from exiting, to keep the bots running.
	stopIdling chan struct{}
	// serveMux is where all our webhook paths are added for the server to use.
	serveMux *http.ServeMux
	// webhookServer is the server in charge of receiving all incoming webhook updates.
	webhookServer *http.Server
	// botMapping keeps track of the data required for each bot. The key is the bot token.
	// TODO: Add support for bot token migration.
	botMapping map[string]*botData
}

// UpdaterOpts defines various fields that can be changed to configure a new Updater.
type UpdaterOpts struct {
	// UnhandledErrFunc provides more flexibility for dealing with previously unhandled errors, such as failures to get
	// updates (when long-polling), or failures to unmarshal.
	// If nil, the error goes to ErrorLog.
	UnhandledErrFunc ErrorFunc
	// ErrorLog specifies an optional logger for unexpected behavior from handlers.
	// If nil, logging is done via the log package's standard logger.
	ErrorLog *log.Logger
	// The dispatcher instance to be used by the updater.
	Dispatcher *Dispatcher
}

// NewUpdater Creates a new Updater, as well as the necessary structures required for the associated Dispatcher.
func NewUpdater(opts *UpdaterOpts) *Updater {
	var unhandledErrFunc ErrorFunc
	var errLog *log.Logger

	// Default dispatcher, no special settings.
	dispatcher := NewDispatcher(nil)

	if opts != nil {
		if opts.Dispatcher != nil {
			dispatcher = opts.Dispatcher
		}

		unhandledErrFunc = opts.UnhandledErrFunc
		errLog = opts.ErrorLog
	}

	return &Updater{
		ErrorLog:         errLog,
		UnhandledErrFunc: unhandledErrFunc,
		Dispatcher:       dispatcher,
	}
}

func (u *Updater) logf(format string, args ...interface{}) {
	if u.ErrorLog != nil {
		u.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// PollingOpts represents the optional values to start long polling.
type PollingOpts struct {
	// DropPendingUpdates toggles whether to drop updates which were sent before the bot was started.
	// Note: the DropPendingUpdates and DisableWebhookDeletion fields are mutually exclusive; both cannot be enabled at
	//  the same time.
	DropPendingUpdates bool
	// DisableWebhookDeletion stops the updater from deleting webhooks before polling starts.
	// By default, the updater deletes any existing webhooks, to ensure that getUpdates works as expected.
	// Note: the DropPendingUpdates and DisableWebhookDeletion fields are mutually exclusive; both cannot be enabled at
	//  the same time.
	DisableWebhookDeletion bool
	// GetUpdatesOpts represents the opts passed to GetUpdates.
	// Note: It is recommended you edit the values here when running in production environments.
	// Suggestions include:
	//    - Changing the "GetUpdatesOpts.AllowedUpdates" to only refer to updates relevant to your bot's functionality.
	//    - Using a non-0 "GetUpdatesOpts.Timeout" value. This is how "long" telegram will hold the long-polling call
	//    while waiting for new messages. A value of 0 causes telegram to reply immediately, which will then cause
	//    your bot to immediately ask for more updates. While this can seem fine, it will eventually causing
	//    telegram to delay your requests when left running over longer periods. If you are seeing lots
	//    of "context deadline exceeded" errors on GetUpdates, this is likely the cause.
	//    Keep in mind that a timeout of 10 does not mean you only get updates every 10s; by the nature of
	//    long-polling, Telegram responds to your request as soon as new messages are available.
	//    When setting this, it is recommended you set your PollingOpts.Timeout value to be slightly bigger (eg, +1).
	GetUpdatesOpts *gotgbot.GetUpdatesOpts
}

var ErrInvalidPendingUpdateAndWebhookConfiguration = errors.New("pending updates cannot be dropped without webhook deletion; consider setting the offset to -1 instead of enabling DropPendingUpdates")

// StartPolling starts polling updates from telegram using getUpdates long-polling.
// See PollingOpts for optional values to set in production environments.
func (u *Updater) StartPolling(b *gotgbot.Bot, opts *PollingOpts) error {
	if u.botMapping == nil {
		u.botMapping = make(map[string]*botData)
	}

	// This logic is currently mostly duplicated over from the generated getUpdates code.
	// This is a performance improvement to avoid:
	//  - needing to re-allocate new url.values structs.
	//  - needing to convert the 'opt' values to strings.
	//  - unnecessary unmarshalling of multiple full Update structs.
	v := map[string]string{}
	var reqOpts *gotgbot.RequestOpts

	if opts != nil {
		if opts.DropPendingUpdates && opts.DisableWebhookDeletion {
			// opts.DropPendingUpdates depends on webhook deletion being enabled.
			// If webhook deletions are disabled, the user should consider setting the offset to -1 instead.
			// This is not perfect, but achieves nearly the same effect.
			return ErrInvalidPendingUpdateAndWebhookConfiguration
		}

		if !opts.DisableWebhookDeletion {
			// For polling to work, we want to make sure we don't have an existing webhook.
			// Extra perk - we can also use this to drop pending updates!
			_, err := b.DeleteWebhook(&gotgbot.DeleteWebhookOpts{
				DropPendingUpdates: opts.DropPendingUpdates,
				RequestOpts:        reqOpts,
			})
			if err != nil {
				return fmt.Errorf("failed to delete webhook: %w", err)
			}
		}

		if updateOpts := opts.GetUpdatesOpts; updateOpts != nil {
			if updateOpts.RequestOpts != nil {
				reqOpts = updateOpts.RequestOpts
			}

			if updateOpts.Offset != 0 {
				v["offset"] = strconv.FormatInt(updateOpts.Offset, 10)
			}
			if updateOpts.Limit != 0 {
				v["limit"] = strconv.FormatInt(updateOpts.Limit, 10)
			}
			if updateOpts.Timeout != 0 {
				v["timeout"] = strconv.FormatInt(updateOpts.Timeout, 10)
			}
			if updateOpts.AllowedUpdates != nil {
				bs, err := json.Marshal(updateOpts.AllowedUpdates)
				if err != nil {
					return fmt.Errorf("failed to marshal field allowed_updates: %w", err)
				}
				v["allowed_updates"] = string(bs)
			}
		}
	}

	updateChan := make(chan json.RawMessage)
	pollChan := make(chan struct{})

	u.botMapping[b.GetToken()] = &botData{
		bot:        b,
		updateChan: updateChan,
		polling:    pollChan,
	}

	go u.Dispatcher.Start(b, updateChan)
	go u.pollingLoop(b, reqOpts, pollChan, updateChan, v)

	return nil
}

func (u *Updater) pollingLoop(b *gotgbot.Bot, opts *gotgbot.RequestOpts, stopPolling <-chan struct{}, updateChan chan<- json.RawMessage, v map[string]string) {
	for {
		select {
		case <-stopPolling:
			// if anything comes in, stop polling.
			return
		default:
			// otherwise, continue as usual
		}

		// Manually craft the getUpdate calls to improve memory management, reduce json parsing overheads, and
		// unnecessary reallocation of url.Values in the polling loop.
		r, err := b.Request("getUpdates", v, nil, opts)
		if err != nil {
			if u.UnhandledErrFunc != nil {
				u.UnhandledErrFunc(err)
			} else {
				u.logf("Failed to get updates; sleeping 1s: %s", err.Error())
				time.Sleep(time.Second)
			}
			continue

		} else if len(r) == 0 {
			continue
		}

		var rawUpdates []json.RawMessage
		if err := json.Unmarshal(r, &rawUpdates); err != nil {
			if u.UnhandledErrFunc != nil {
				u.UnhandledErrFunc(err)
			} else {
				u.logf("Failed to unmarshal updates: %s", err.Error())
			}
			continue
		}

		if len(rawUpdates) == 0 {
			continue
		}

		var lastUpdate struct {
			UpdateId int64 `json:"update_id"`
		}

		// Only unmarshal the last update, so we can get the next update ID.
		if err := json.Unmarshal(rawUpdates[len(rawUpdates)-1], &lastUpdate); err != nil {
			if u.UnhandledErrFunc != nil {
				u.UnhandledErrFunc(err)
			} else {
				u.logf("Failed to unmarshal last update: %s", err.Error())
			}
			continue
		}

		v["offset"] = strconv.FormatInt(lastUpdate.UpdateId+1, 10)
		for _, updData := range rawUpdates {
			temp := updData // use new mem address to avoid loop conflicts
			updateChan <- temp
		}
	}
}

// Idle starts an infinite loop to avoid the program exciting while the background threads handle updates.
func (u *Updater) Idle() {
	// Create the idling channel
	u.stopIdling = make(chan struct{})

	// Wait until some input is received from the idle channel, which will stop the idling.
	<-u.stopIdling
}

// Stop stops the current updater and dispatcher instances.
func (u *Updater) Stop() error {
	// Stop any running servers.
	if u.webhookServer != nil {
		err := u.webhookServer.Shutdown(context.Background())
		if err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
	}

	// Close all the update channels and polling loops
	for _, data := range u.botMapping {
		// Close polling loops first, to ensure any updates currently being polled have the time to be sent to the
		// updateChan.
		if data.polling != nil {
			close(data.polling)
		}

		// Then, close the updates channel.
		close(data.updateChan)
	}

	// Stop the dispatcher from processing any further updates.
	u.Dispatcher.Stop()

	// Finally, atop idling.
	if u.stopIdling != nil {
		close(u.stopIdling)
	}
	return nil
}

// StartWebhook starts the webhook server for a single bot instance.
// This does NOT set the webhook on telegram - this should be done by the caller.
// The opts parameter allows for specifying various webhook settings.
func (u *Updater) StartWebhook(b *gotgbot.Bot, urlPath string, opts WebhookOpts) error {
	if u.webhookServer != nil {
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
		u.botMapping = make(map[string]*botData)
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

	u.botMapping[b.GetToken()] = &botData{
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
		_, err := data.bot.SetWebhook(strings.Join([]string{strings.TrimSuffix(domain, "/"), data.urlPath}, "/"), opts)
		if err != nil {
			// Extract the botID, so we don't intentionally log the token
			botId := strings.Split(data.bot.GetToken(), ":")[0]
			return fmt.Errorf("failed to set webhook for %s: %w", botId, err)
		}
	}
	return nil
}

// StartServer starts the webhook server for all the bots added via AddWebhook.
// It is recommended to call this BEFORE calling setWebhooks.
// The opts parameter allows for specifying TLS settings.
func (u *Updater) StartServer(opts WebhookOpts) error {
	if u.serveMux == nil {
		u.serveMux = http.NewServeMux()
	}

	var tls bool
	switch {
	case opts.CertFile == "" && opts.KeyFile == "":
		tls = false
	case opts.CertFile != "" && opts.KeyFile != "":
		tls = true
	default:
		return ErrMissingCertOrKeyFile
	}

	ln, err := net.Listen(opts.GetListenNet(), opts.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s:%s: %w", opts.ListenNet, opts.ListenAddr, err)
	}

	u.webhookServer = &http.Server{
		Handler:           u.serveMux,
		ReadTimeout:       opts.ReadTimeout,
		ReadHeaderTimeout: opts.ReadHeaderTimeout,
	}

	go func() {
		if tls {
			err = u.webhookServer.ServeTLS(ln, opts.CertFile, opts.KeyFile)
		} else {
			err = u.webhookServer.Serve(ln)
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic("http server failed: " + err.Error())
		}
	}()

	return nil
}
