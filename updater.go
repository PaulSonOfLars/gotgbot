package gotgbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/PaulSonOfLars/gotgbot/ext"
)

// Updater The main updater process. Receives incoming updates, then sends them to the dispatcher goroutine
// via an update channel for them to be handled.
type Updater struct {
	Bot          *ext.Bot
	Updates      chan *RawUpdate
	Dispatcher   *Dispatcher
	UpdateGetter ext.Requester
}

// NewUpdater Creates a new updater struct, paired with the necessary dispatcher and bot structs.
func NewUpdater(l *zap.Logger, token string) (*Updater, error) {
	u := &Updater{}
	bot, err := ext.NewBot(l, token)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get bot")
	}

	u.Bot = bot
	u.Updates = make(chan *RawUpdate)
	u.Dispatcher = NewDispatcher(u.Bot, u.Updates)
	u.UpdateGetter = ext.Requester{
		Client: http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Second * 5,
		},
		ApiUrl: ext.ApiUrl,
	}
	ok, err := u.RemoveWebhook() // just in case
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("failed to remove webhook")
	}
	return u, nil
}

// StartPolling Starts the polling logic
func (u Updater) StartPolling() error {
	go u.Dispatcher.Start()
	go u.startPolling(false)
	return nil
}

// StartCleanPolling Starts clean polling (ignoring stale updates)
func (u Updater) StartCleanPolling() error {
	go u.Dispatcher.Start()
	go u.startPolling(true)
	return nil
}

func (u Updater) startPolling(clean bool) {
	v := url.Values{}
	v.Add("offset", strconv.Itoa(0))
	v.Add("timeout", strconv.Itoa(0))
	offset := 0
	for {
		// Note: use updateGetter.Get instead of u.Bot.Get to use the updater timeout instead of bot command timeout
		r, err := u.UpdateGetter.Get(u.Bot.Logger, u.Bot.Token, "getUpdates", v)
		if err != nil {
			u.Bot.Logger.Errorw("unable to getUpdates", zap.Error(err))
			u.Bot.Logger.Error("Sleeping for 1 second...")
			time.Sleep(time.Second)
			continue

		} else if r != nil {
			var rawUpdates []json.RawMessage
			if err := json.Unmarshal(r, &rawUpdates); err != nil {
				u.Bot.Logger.Errorw("failed to unmarshal update while polling",
					zap.Field{
						Key:    "result",
						Type:   zapcore.StringType,
						String: string(r)},
					zap.Error(err))
				continue
			}
			if len(rawUpdates) > 0 {
				// parse last one here
				lastUpd, err := initUpdate(RawUpdate(rawUpdates[len(rawUpdates)-1]), *u.Bot)
				if err != nil {
					u.Bot.Logger.Errorw("failed to init update while polling",
						zap.Field{
							Key:    "result",
							Type:   zapcore.StringType,
							String: string(r)},
						zap.Error(err))
					continue
				}
				offset = lastUpd.UpdateId + 1
				v.Set("offset", strconv.Itoa(offset))
				if clean {
					continue
				}
			} else if len(rawUpdates) == 0 { // TODO: this is unsustainable, and may eventually break on higher loads.
				clean = false
			}

			for _, updData := range rawUpdates {
				temp := RawUpdate(updData) // necessary to avoid memory stuff from loops
				u.Updates <- &temp
			}
		}
	}
}

// Idle sets the main thread to idle, allowing the background processes to run as expected (dispatcher and update handlers)
func (u Updater) Idle() {
	for {
		time.Sleep(1 * time.Second)
	}
}

// TODO: finish handling updates on sigint

// StartWebhook Start the webhook server
func (u Updater) StartWebhook(webhook ext.Webhook) {
	go u.Dispatcher.Start()
	http.HandleFunc("/"+webhook.ServePath, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		temp := RawUpdate(bytes)
		u.Updates <- &temp
	})
	go func() {
		// todo: TLS when using certs
		err := http.ListenAndServe(webhook.GetListenUrl(), nil)
		if err != nil {
			u.Bot.Logger.Errorw("failed to init update while polling", zap.Error(errors.WithStack(err)))
		}
	}()
}

// RemoveWebhook remove the webhook url from telegram servers
func (u Updater) RemoveWebhook() (bool, error) {
	return u.Bot.DeleteWebhook()
}

// SetWebhook Set the webhook url for telegram to contact with updates
func (u Updater) SetWebhook(path string, webhook ext.Webhook) (bool, error) {
	return u.Bot.SetWebhook(path, webhook)
}

// GetWebhookInfo Get webhook info from telegram
func (u Updater) GetWebhookInfo() (*ext.WebhookInfo, error) {
	return u.Bot.GetWebhookInfo()
}
