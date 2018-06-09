package gotgbot

import (
	"encoding/json"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"time"
)

type Updater struct {
	Bot        ext.Bot
	updates    chan Update
	Dispatcher *Dispatcher
}

func NewUpdater(token string) *Updater {
	u := &Updater{}
	u.Bot = ext.Bot{Token: token}
	u.updates = make(chan Update)
	u.Dispatcher = NewDispatcher(u.Bot, u.updates)
	return u
}

func (u Updater) StartPolling() error {
	if _, err := u.Bot.GetMe(); err != nil {
		return err
	}
	go u.Dispatcher.Start()
	go u.startPolling(false)
	return nil
}

func (u Updater) StartCleanPolling() error {
	if _, err := u.Bot.GetMe(); err != nil {
		return err
	}
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
		r, err := ext.Get(u.Bot, "getUpdates", v)
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
				lastUpd := initUpdate(rawUpdates[len(rawUpdates)-1], u.Bot)
				offset = lastUpd.UpdateId + 1
				v.Set("offset", strconv.Itoa(offset))
				if clean {
					continue
				}
			} else if len(rawUpdates) == 0 { // TODO: this is unsustainable, and may eventually break on higher loads.
				clean = false
			}

			for _, updData := range rawUpdates {
				upd := initUpdate(updData, u.Bot)
				u.updates <- upd
			}
		}
	}
}

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
