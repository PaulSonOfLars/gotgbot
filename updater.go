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
	"github.com/sirupsen/logrus"

	"github.com/PaulSonOfLars/gotgbot/ext"
)

// Updater The main updater process. Receives incoming updates, then sends them to the dispatcher goroutine
// via an update channel for them to be handled.
type Updater struct {
	Bot          *ext.Bot
	Updates      chan *RawUpdate
	Dispatcher   *Dispatcher
	UpdateGetter *ext.TgBotGetter
}

// NewUpdater Creates a new updater struct, paired with the necessary dispatcher and bot structs.
func NewUpdater(token string) (*Updater, error) {
	u := &Updater{}
	user, err := ext.Bot{Token: token, Logger: logrus.New()}.GetMe()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create new updater")
	}
	u.Bot = &ext.Bot{
		Token:     token,
		Id:        user.Id,
		FirstName: user.FirstName,
		UserName:  user.Username,
		Logger:    logrus.New(),
	}
	u.Updates = make(chan *RawUpdate)
	u.Dispatcher = NewDispatcher(u.Bot, u.Updates)
	u.UpdateGetter = &ext.TgBotGetter{
		Client: &http.Client{
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
		r, err := u.UpdateGetter.Get(*u.Bot, "getUpdates", v)
		if err != nil {
			logrus.WithError(err).Error("unable to getUpdates")
			continue

		} else if !r.Ok {
			logrus.Errorf("getUpdates error: %v", r.Description)
			logrus.Errorf("Sleeping for 1 second...")
			time.Sleep(time.Second)
			continue

		} else if r.Result != nil {
			var rawUpdates []json.RawMessage
			if err := json.Unmarshal(r.Result, &rawUpdates); err != nil {
				logrus.WithError(err).Error("failed to unmarshal update while polling", r.Result)
				continue
			}
			if len(rawUpdates) > 0 {
				// parse last one here
				lastUpd, err := initUpdate(RawUpdate(rawUpdates[len(rawUpdates)-1]), *u.Bot)
				if err != nil {
					logrus.WithError(err).Error("failed to init update while polling", r.Result)
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

type Webhook struct {
	Serve     string // base url to where you listen
	ServePath string // path you listen to
	ServePort int    // port you listen on
	URL       string // where you set the webhook to send to
	// CertPath       string   // TODO
	MaxConnections int      // max connections; max 100, default 40
	AllowedUpdates []string // which updates to allow
}

func (w Webhook) GetListenUrl() string {
	if w.Serve == "" {
		w.Serve = "0.0.0.0"
	}
	if w.ServePort == 0 {
		w.ServePort = 443
	}
	return fmt.Sprintf("%s:%d", w.Serve, w.ServePort)
}

// StartWebhook Start the webhook server
func (u Updater) StartWebhook(webhook Webhook) {
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
			logrus.Fatal(errors.WithStack(err))
		}
	}()
}

// RemoveWebhook remove the webhook url from telegram servers
func (u Updater) RemoveWebhook() (bool, error) {
	r, err := ext.Get(*u.Bot, "deleteWebhook", nil)
	if err != nil {
		return false, errors.Wrapf(err, "failed to remove webhook")
	}
	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

// SetWebhook Set the webhook url for telegram to contact with updates
func (u Updater) SetWebhook(path string, webhook Webhook) (bool, error) {
	allowedUpdates := webhook.AllowedUpdates
	if allowedUpdates == nil {
		allowedUpdates = []string{}
	}
	allowed, err := json.Marshal(allowedUpdates)
	if err != nil {
		return false, errors.Wrap(err, "cannot marshal allowedUpdates")
	}

	v := url.Values{}
	v.Add("url", webhook.URL+"/"+path)
	// v.Add("certificate", ) // todo: add certificate support
	v.Add("max_connections", strconv.Itoa(webhook.MaxConnections))
	v.Add("allowed_updates", string(allowed))

	r, err := ext.Get(*u.Bot, "setWebhook", v)
	if err != nil {
		return false, errors.Wrap(err, "failed to set webhook")
	}

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

type WebhookInfo struct {
	URL                  string   `json:"url"`
	HasCustomCertificate bool     `json:"has_custom_certificate"`
	PendingUpdateCount   int      `json:"pending_update_count"`
	LastErrorDate        int      `json:"last_error_date"`
	LastErrorMessage     int      `json:"last_error_message"`
	MaxConnections       int      `json:"max_connections"`
	AllowedUpdates       []string `json:"allowed_updates"`
}

// GetWebhookInfo Get webhook info from telegram servers
func (u Updater) GetWebhookInfo() (*WebhookInfo, error) {
	r, err := ext.Get(*u.Bot, "getWebhookInfo", nil)
	if err != nil {
		return nil, err
	}

	var wh WebhookInfo
	return &wh, json.Unmarshal(r.Result, &wh)

}
