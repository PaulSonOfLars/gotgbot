package gotgbot

import (
	"encoding/json"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"time"
	"net/http"
	"io/ioutil"
	"fmt"
)

type Updater struct {
	Bot        *ext.Bot
	updates    chan Update
	Dispatcher *Dispatcher
}

func NewUpdater(token string) (*Updater, error) {
	u := &Updater{}
	user, err := ext.Bot{Token: token}.GetMe()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create new updater")
	}
	u.Bot = &ext.Bot{
		Token:     token,
		Id:        user.Id,
		FirstName: user.FirstName,
		UserName:  user.Username,
	}
	u.updates = make(chan Update)
	u.Dispatcher = NewDispatcher(*u.Bot, u.updates)
	ok, err := u.RemoveWebhook() // just in case
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("failed to remove webhook")
	}
	return u, nil
}

func (u Updater) StartPolling() error {
	go u.Dispatcher.Start()
	go u.startPolling(false)
	return nil
}

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
		r, err := ext.Get(*u.Bot, "getUpdates", v)
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
			json.Unmarshal(r.Result, &rawUpdates)
			if len(rawUpdates) > 0 {
				// parse last one here
				lastUpd := initUpdate(rawUpdates[len(rawUpdates)-1], *u.Bot)
				offset = lastUpd.UpdateId + 1
				v.Set("offset", strconv.Itoa(offset))
				if clean {
					continue
				}
			} else if len(rawUpdates) == 0 { // TODO: this is unsustainable, and may eventually break on higher loads.
				clean = false
			}

			for _, updData := range rawUpdates {
				upd := initUpdate(updData, *u.Bot)
				u.updates <- upd
			}
		}
	}
}

// todo: move this into dispatcher update processor to updater CPU cycles
func initUpdate(data json.RawMessage, bot ext.Bot) Update {
	var upd Update
	json.Unmarshal(data, &upd)
	if upd.Message != nil {
		upd.EffectiveMessage = upd.Message
		upd.EffectiveChat = upd.Message.Chat
		upd.EffectiveUser = upd.Message.From

	} else if upd.EditedMessage != nil {
		upd.EffectiveMessage = upd.EditedMessage
		upd.EffectiveChat = upd.EditedMessage.Chat
		upd.EffectiveUser = upd.EditedMessage.From

	} else if upd.ChannelPost != nil {
		upd.EffectiveMessage = upd.ChannelPost
		upd.EffectiveChat = upd.ChannelPost.Chat

	} else if upd.EditedChannelPost != nil {
		upd.EffectiveMessage = upd.EditedChannelPost
		upd.EffectiveChat = upd.EditedChannelPost.Chat

	} else if upd.InlineQuery != nil {
		upd.EffectiveMessage = upd.InlineQuery
		upd.EffectiveUser = upd.InlineQuery.From

	} else if upd.CallbackQuery != nil && upd.CallbackQuery.Message != nil {
		upd.EffectiveMessage = upd.CallbackQuery.Message
		upd.EffectiveChat = upd.CallbackQuery.Message.Chat
		upd.EffectiveUser = upd.CallbackQuery.From

	} else if upd.ChosenInlineResult != nil {
		upd.EffectiveUser = upd.ChosenInlineResult.From

	} else if upd.ShippingQuery != nil {
		upd.EffectiveUser = upd.ShippingQuery.From

	} else if upd.PreCheckoutQuery != nil {
		upd.EffectiveUser = upd.PreCheckoutQuery.From
	}

	if upd.EffectiveMessage != nil {
		upd.EffectiveMessage.Bot = bot
		if upd.EffectiveMessage.ReplyToMessage != nil {
			upd.EffectiveMessage.ReplyToMessage.Bot = bot
			upd.EffectiveMessage.ReplyToMessage.From.Bot = bot
		}
	}
	if upd.EffectiveChat != nil {
		upd.EffectiveChat.Bot = bot
	}
	if upd.EffectiveUser != nil {
		upd.EffectiveUser.Bot = bot
	}
	return upd
}

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
	//CertPath       string   // TODO
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

func (u Updater) StartWebhook(webhook Webhook) {
	go u.Dispatcher.Start()
	http.HandleFunc("/"+webhook.ServePath, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		u.updates <- initUpdate(bytes, *u.Bot)
	})
	go func() {
		// todo: TLS when using certs
		err := http.ListenAndServe(webhook.GetListenUrl(), nil)
		if err != nil {
			logrus.Fatal(errors.WithStack(err))
		}
	}()
}

func (u Updater) RemoveWebhook() (bool, error) {
	r, err := ext.Get(*u.Bot, "deleteWebhook", nil)
	if err != nil {
		return false, errors.Wrapf(err, "failed to remove webhook")
	}
	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

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
	v.Add("url", webhook.URL + "/" + path)
	//v.Add("certificate", ) // todo: add certificate support
	v.Add("max_connections", strconv.Itoa(webhook.MaxConnections))
	v.Add("allowed_updates", string(allowed))

	r, err := ext.Get(*u.Bot, "setWebhook", v)
	if err != nil {
		return false, errors.Wrap(err, "failed to set webhook")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)
	return bb, nil
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

func (u Updater) GetWebhookInfo() (*WebhookInfo, error) {
	r, err := ext.Get(*u.Bot, "getWebhookInfo", nil)
	if err != nil {
		return nil, err
	}

	var wh WebhookInfo
	json.Unmarshal(r.Result, &wh)
	return &wh, nil

}
