package ext

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var ErrMissingCertOrKeyFile = errors.New("missing certfile or keyfile")

type Updater struct {
	Dispatcher *Dispatcher
	UpdateChan chan json.RawMessage
	ErrorLog   *log.Logger

	stopIdling chan bool
	running    chan bool
	server     *http.Server
}

var errorLog = log.New(os.Stderr, "ERROR", log.LstdFlags)

type UpdaterOpts struct {
	ErrorLog *log.Logger

	DispatcherOpts DispatcherOpts
}

// NewUpdater Creates a new Updater, as well as the necessary structures required for the associated Dispatcher.
func NewUpdater(opts *UpdaterOpts) Updater {
	errLog := errorLog
	var dispatcherOpts DispatcherOpts

	if opts != nil {
		if opts.ErrorLog != nil {
			errLog = opts.ErrorLog
		}

		dispatcherOpts = opts.DispatcherOpts
	}

	updateChan := make(chan json.RawMessage)
	return Updater{
		ErrorLog:   errLog,
		Dispatcher: NewDispatcher(updateChan, &dispatcherOpts),
		UpdateChan: updateChan,
	}
}

// PollingOpts represents the optional values to start long polling.
type PollingOpts struct {
	// DropPendingUpdates decides whether or not to drop "pending" updates; these are updates which were sent before
	// the bot was started.
	DropPendingUpdates bool
	// GetUpdatesOpts represents the opts passed to GetUpdates.
	// Note: It is recommended you edit the values here when running in production environments.
	// Changes might include:
	//    - Changing the "GetUpdatesOpts.AllowedUpates" to only refer to relevant updates
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

// StartPolling starts polling updates from telegram using the getUdpates long-polling method.
// See the PollingOpts for optional values to set in production environments.
func (u *Updater) StartPolling(b *gotgbot.Bot, opts *PollingOpts) error {
	// TODO: De-duplicate this code.
	// This logic is currently mostly duplicated over from the generated getUpdates code.
	// This is a performance improvement to avoid:
	// - needing to re-allocate new url.values structs.
	// - needing to convert the opt values to strings to pass to the values.
	// - unnecessary unmarshalling of the (possibly multiple) full Update structs.
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

	go u.Dispatcher.Start(b)
	go u.pollingLoop(b, reqOpts, dropPendingUpdates, v)

	return nil
}

func (u *Updater) pollingLoop(b *gotgbot.Bot, opts *gotgbot.RequestOpts, dropPendingUpdates bool, v map[string]string) {

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

		r, err := b.Post("getUpdates", v, nil, opts)
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
			u.UpdateChan <- temp
		}
	}
}

// Idle starts an infinite loop to avoid the program exciting while the background threads handle updates.
func (u *Updater) Idle() {
	u.stopIdling = make(chan bool)

	for {
		select {
		case <-u.stopIdling:
			return
		default:
			// continue as usual
		}
		time.Sleep(1 * time.Second)
	}
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

	close(u.UpdateChan)

	u.Dispatcher.Stop()

	if u.stopIdling != nil {
		// stop idling
		u.stopIdling <- false
		close(u.stopIdling)
	}
	return nil
}

// StartWebhook Starts the webhook server. The opts parameter allows for specifying TLS settings.
func (u *Updater) StartWebhook(b *gotgbot.Bot, opts WebhookOpts) error {
	var tls bool
	if opts.CertFile == "" && opts.KeyFile == "" {
		tls = false
	} else if opts.CertFile != "" && opts.KeyFile != "" {
		tls = true
	} else {
		return ErrMissingCertOrKeyFile
	}

	go u.Dispatcher.Start(b)

	mux := http.NewServeMux()
	mux.HandleFunc("/"+opts.URLPath, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		u.UpdateChan <- bytes
	})

	u.server = &http.Server{
		Addr:    opts.GetListenAddr(),
		Handler: mux,
	}

	go func() {
		var err error
		if tls {
			err = u.server.ListenAndServeTLS(opts.CertFile, opts.KeyFile)
		} else {
			err = u.server.ListenAndServe()
		}
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			panic("http server failed: " + err.Error())
		}
	}()

	return nil
}

type WebhookOpts struct {
	Listen  string
	Port    int
	URLPath string

	CertFile string
	KeyFile  string
}

// GetListenAddr returns the local listening address, including port.
func (w *WebhookOpts) GetListenAddr() string {
	if w.Listen == "" {
		w.Listen = "0.0.0.0"
	}
	if w.Port == 0 {
		w.Port = 443
	}
	return fmt.Sprintf("%s:%d", w.Listen, w.Port)
}

// GetWebhookURL returns the domain in the form domain/path.
// eg: example.com/super_secret_token
func (w *WebhookOpts) GetWebhookURL(domain string) string {
	return fmt.Sprintf("%s/%s", domain, w.URLPath)
}
