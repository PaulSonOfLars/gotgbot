package ext

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	ErrNotFound             = errors.New("not found")
	ErrEmptyPath            = errors.New("empty path")
)

type ErrorFunc func(error)

type Updater struct {
	// Dispatcher is where all the incoming updates are sent to be processed.
	// The Dispatcher runs in a separate goroutine, allowing for parallel update processing and dispatching.
	// Once the Updater has received an update, it sends it to the Dispatcher over a JSON channel.
	Dispatcher UpdateDispatcher

	// UnhandledErrFunc provides more flexibility for dealing with previously unhandled errors, such as failures to get
	// updates (when long-polling), or failures to unmarshal.
	// If nil, the error goes to ErrorLog.
	UnhandledErrFunc ErrorFunc
	// ErrorLog specifies an optional logger for unexpected behavior from handlers.
	// If nil, logging is done via the log package's standard logger.
	ErrorLog *log.Logger

	// stopIdling is the channel that blocks the main thread from exiting, to keep the bots running.
	stopIdling chan struct{}
	// webhookServer is the server in charge of receiving all incoming webhook updates.
	webhookServer *http.Server

	// botMapping keeps track of the data required for each bot, in a thread-safe manner.
	botMapping botMapping
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
}

// NewUpdater Creates a new Updater, as well as a Dispatcher and any optional updater configurations (via UpdaterOpts).
func NewUpdater(dispatcher UpdateDispatcher, opts *UpdaterOpts) *Updater {
	var unhandledErrFunc ErrorFunc
	var errLog *log.Logger

	if opts != nil {
		unhandledErrFunc = opts.UnhandledErrFunc
		errLog = opts.ErrorLog
	}

	return &Updater{
		Dispatcher:       dispatcher,
		UnhandledErrFunc: unhandledErrFunc,
		ErrorLog:         errLog,
		botMapping: botMapping{
			errFunc:  unhandledErrFunc,
			errorLog: errLog,
		},
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
	// This also implicitly enables webhook deletion.
	DropPendingUpdates bool
	// EnableWebhookDeletion deletes any existing webhooks to ensure that the updater works fine.
	EnableWebhookDeletion bool
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

// StartPolling starts polling updates from telegram using getUpdates long-polling.
// See PollingOpts for optional values to set in production environments.
func (u *Updater) StartPolling(b *gotgbot.Bot, opts *PollingOpts) error {
	// This logic is currently mostly duplicated over from the generated getUpdates code.
	// This is a performance improvement to avoid:
	//  - needing to re-allocate new url.values structs.
	//  - needing to convert the 'opt' values to strings.
	//  - unnecessary unmarshalling of multiple full Update structs.
	v := map[string]string{}
	var reqOpts *gotgbot.RequestOpts

	if opts != nil {
		if opts.EnableWebhookDeletion || opts.DropPendingUpdates {
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

	bData, err := u.botMapping.addBot(b, "", "")
	if err != nil {
		return fmt.Errorf("failed to add bot with long polling: %w", err)
	}

	go u.Dispatcher.Start(b, bData.updateChan)
	go u.pollingLoop(bData, reqOpts, v)

	return nil
}

func (u *Updater) pollingLoop(bData *botData, opts *gotgbot.RequestOpts, v map[string]string) {
	bData.updateWriterControl.Add(1)
	defer bData.updateWriterControl.Done()

	for {
		// Check if updater loop has been terminated.
		if bData.shouldStopUpdates() {
			return
		}

		// Manually craft the getUpdate calls to improve memory management, reduce json parsing overheads, and
		// unnecessary reallocation of url.Values in the polling loop.
		r, err := bData.bot.Request("getUpdates", v, nil, opts)
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
			bData.updateChan <- temp
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
//
// When using long polling, Stop() will wait for the getUpdates call to return, which may cause a delay due to the
// request timeout.
func (u *Updater) Stop() error {
	// Stop any running servers.
	if u.webhookServer != nil {
		err := u.webhookServer.Shutdown(context.Background())
		if err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
	}

	// Close all existing bot channels.
	u.StopAllBots()

	// Stop the dispatcher from processing any further updates.
	u.Dispatcher.Stop()

	// Finally, atop idling.
	if u.stopIdling != nil {
		close(u.stopIdling)
	}
	return nil
}

func (u *Updater) StopBot(token string) bool {
	bData, ok := u.botMapping.removeBot(token)
	if !ok {
		return false
	}

	bData.stop()
	return true
}

func (u *Updater) StopAllBots() {
	for _, bData := range u.botMapping.removeAllBots() {
		bData.stop()
	}
}

// StartWebhook starts the webhook server for a single bot instance.
// This does NOT set the webhook on telegram - this should be done by the caller.
// The opts parameter allows for specifying various webhook settings.
func (u *Updater) StartWebhook(b *gotgbot.Bot, urlPath string, opts WebhookOpts) error {
	if u.webhookServer != nil {
		return ErrExpectedEmptyServer
	}

	err := u.AddWebhook(b, urlPath, &AddWebhookOpts{SecretToken: opts.SecretToken})
	if err != nil {
		return fmt.Errorf("failed to add webhook: %w", err)
	}

	return u.StartServer(opts)
}

// AddWebhookOpts stores any optional parameters for the Updater.AddWebhook method.
type AddWebhookOpts struct {
	// The secret token to be used to validate webhook authenticity.
	SecretToken string
}

// AddWebhook prepares the webhook server to receive webhook updates for one bot, on a specific path.
func (u *Updater) AddWebhook(b *gotgbot.Bot, urlPath string, opts *AddWebhookOpts) error {
	// We expect webhooks to use unique URL paths; otherwise, we wouldnt be able to differentiate them from polling, or
	// from each other.
	if urlPath == "" {
		return fmt.Errorf("expected a non-empty url path: %w", ErrEmptyPath)
	}

	secretToken := ""
	if opts != nil {
		secretToken = opts.SecretToken
	}

	bData, err := u.botMapping.addBot(b, urlPath, secretToken)
	if err != nil {
		return fmt.Errorf("failed to add webhook for bot: %w", err)
	}

	// Webhook has been added; relevant dispatcher should also be started.
	go u.Dispatcher.Start(b, bData.updateChan)
	return nil
}

// SetAllBotWebhooks sets all the webhooks for the bots that have been added to this updater via AddWebhook.
func (u *Updater) SetAllBotWebhooks(domain string, opts *gotgbot.SetWebhookOpts) error {
	for _, data := range u.botMapping.getBots() {
		_, err := data.bot.SetWebhook(strings.Join([]string{strings.TrimSuffix(domain, "/"), data.urlPath}, "/"), opts)
		if err != nil {
			// Extract the botID, so we don't intentionally log the token
			botId := strings.Split(data.bot.Token, ":")[0]
			return fmt.Errorf("failed to set webhook for %s: %w", botId, err)
		}
	}
	return nil
}

// GetHandlerFunc returns the http.HandlerFunc responsible for processing incoming webhook updates.
// It is provided to allow for an alternative to the StartServer method using a user-defined http server.
func (u *Updater) GetHandlerFunc(pathPrefix string) http.HandlerFunc {
	return u.botMapping.getHandlerFunc(pathPrefix)
}

// StartServer starts the webhook server for all the bots added via AddWebhook.
// It is recommended to call this BEFORE calling setWebhooks.
// The opts parameter allows for specifying TLS settings.
func (u *Updater) StartServer(opts WebhookOpts) error {
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

	if u.webhookServer != nil {
		return ErrExpectedEmptyServer
	}

	u.webhookServer = &http.Server{
		Handler:           u.GetHandlerFunc("/"),
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
