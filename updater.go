package gotgbot

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gotgbot/ext"
	"net/url"
	"strconv"
	"time"
)

type Updater struct {
	bot        ext.Bot
	updates    chan Update
	Dispatcher Dispatcher
}

func NewUpdater(token string) Updater {
	u := Updater{}
	u.bot = ext.Bot{Token: token}
	u.updates = make(chan Update)
	u.Dispatcher = NewDispatcher(u.bot, u.updates)
	return u
}

func (u Updater) StartPolling() {
	go u.Dispatcher.Start()
	go u.startPolling()
}

func (u Updater) startPolling() {
	v := url.Values{}
	v.Add("offset", strconv.Itoa(0))
	v.Add("timeout", strconv.Itoa(0))
	offset := 0
	for {
		r, err := ext.Get(u.bot, "getUpdates", v)
		if err != nil {
			logrus.WithError(err).Error("unable to getUpdates")
			continue

		} else if !r.Ok {
			logrus.Errorf("getUpdates error: %v", r.Description)
			logrus.Errorf("Sleeping for 1 second...")
			time.Sleep(time.Second)
			continue

		} else if r.Result != nil {
			//fmt.Println(r)
			var res []Update
			json.Unmarshal(r.Result, &res)
			for _, upd := range res {
				if upd.Message != nil {
					upd.EffectiveMessage = u.bot.NewMessage(upd.Message)
					upd.EffectiveChat = u.bot.NewChat(upd.Message.Chat)
					upd.EffectiveUser = u.bot.NewUser(upd.Message.From)

				} else if upd.EditedMessage != nil {
					upd.EffectiveMessage = u.bot.NewMessage(upd.EditedMessage)
					upd.EffectiveChat = u.bot.NewChat(upd.EditedMessage.Chat)
					upd.EffectiveUser = u.bot.NewUser(upd.EditedMessage.From)

				} else if upd.ChannelPost != nil {
					upd.EffectiveMessage = u.bot.NewMessage(upd.ChannelPost)
					upd.EffectiveChat = u.bot.NewChat(upd.ChannelPost.Chat)

				} else if upd.EditedChannelPost != nil {
					upd.EffectiveMessage = u.bot.NewMessage(upd.EditedChannelPost)
					upd.EffectiveChat = u.bot.NewChat(upd.EditedChannelPost.Chat)

				} else if upd.InlineQuery != nil {
					upd.EffectiveMessage = u.bot.NewMessage(upd.InlineQuery)
					upd.EffectiveUser = u.bot.NewUser(upd.InlineQuery.From)

				} else if upd.CallbackQuery != nil && upd.CallbackQuery.Message != nil {
					upd.EffectiveChat = u.bot.NewChat(upd.CallbackQuery.Message.Chat)
					upd.EffectiveUser = u.bot.NewUser(upd.CallbackQuery.From)

				} else if upd.ChosenInlineResult != nil {
					upd.EffectiveUser = u.bot.NewUser(upd.ChosenInlineResult.From)

				} else if upd.ShippingQuery != nil {
					upd.EffectiveUser = u.bot.NewUser(upd.ShippingQuery.From)

				} else if upd.PreCheckoutQuery != nil {
					upd.EffectiveUser = u.bot.NewUser(upd.PreCheckoutQuery.From)

				}

				u.updates <- upd
			}
			if len(res) > 0 {
				offset = res[len(res)-1].UpdateId + 1
			}
		}

		v.Set("offset", strconv.Itoa(offset))
	}
}

func (u Updater) Idle() {
	for {
		time.Sleep(1 * time.Second)
	}

}

// TODO: finish handling updates on sigint
