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
			var res []Update
			json.Unmarshal(r.Result, &res)
			if len(res) > 0 {
				offset = res[len(res)-1].UpdateId + 1
				v.Set("offset", strconv.Itoa(offset))
				if clean {
					continue
				}
			} else if len(res) == 0 { // TODO: this is unsustainable, and may eventually break on higher loads.
				clean = false
			}

			for _, upd := range res {
				if upd.Message != nil {
					upd.EffectiveMessage = u.Bot.NewMessage(upd.Message)
					upd.EffectiveChat = u.Bot.NewChat(upd.Message.Chat)
					upd.EffectiveUser = u.Bot.NewUser(upd.Message.From)

				} else if upd.EditedMessage != nil {
					upd.EffectiveMessage = u.Bot.NewMessage(upd.EditedMessage)
					upd.EffectiveChat = u.Bot.NewChat(upd.EditedMessage.Chat)
					upd.EffectiveUser = u.Bot.NewUser(upd.EditedMessage.From)

				} else if upd.ChannelPost != nil {
					upd.EffectiveMessage = u.Bot.NewMessage(upd.ChannelPost)
					upd.EffectiveChat = u.Bot.NewChat(upd.ChannelPost.Chat)

				} else if upd.EditedChannelPost != nil {
					upd.EffectiveMessage = u.Bot.NewMessage(upd.EditedChannelPost)
					upd.EffectiveChat = u.Bot.NewChat(upd.EditedChannelPost.Chat)

				} else if upd.InlineQuery != nil {
					upd.EffectiveMessage = u.Bot.NewMessage(upd.InlineQuery)
					upd.EffectiveUser = u.Bot.NewUser(upd.InlineQuery.From)

				} else if upd.CallbackQuery != nil && upd.CallbackQuery.Message != nil {
					upd.EffectiveChat = u.Bot.NewChat(upd.CallbackQuery.Message.Chat)
					upd.EffectiveUser = u.Bot.NewUser(upd.CallbackQuery.From)

				} else if upd.ChosenInlineResult != nil {
					upd.EffectiveUser = u.Bot.NewUser(upd.ChosenInlineResult.From)

				} else if upd.ShippingQuery != nil {
					upd.EffectiveUser = u.Bot.NewUser(upd.ShippingQuery.From)

				} else if upd.PreCheckoutQuery != nil {
					upd.EffectiveUser = u.Bot.NewUser(upd.PreCheckoutQuery.From)

				}

				u.updates <- upd
			}
		}
	}
}

func (u Updater) Idle() {
	for {
		time.Sleep(1 * time.Second)
	}

}

// TODO: finish handling updates on sigint
