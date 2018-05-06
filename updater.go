package gotgbot

import (
	"gotgbot/ext"
	"encoding/json"
	"log"
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
	for {
		r := ext.Get(u.bot, "getUpdates", v)
		if !r.Ok {
			log.Println(r)
			log.Fatal("You done goofed, API Res for getUpdates was not OK")
		}
		offset := 0
		if r.Result != nil {
			//fmt.Println(r)
			var res []Update
			json.Unmarshal(r.Result, &res)
			for _, upd := range res {
				switch upd.Message {
				case upd.Message:
					upd.Effective_message = u.bot.NewMessage(upd.Message)
				case upd.Edited_message:
					upd.Effective_message = u.bot.NewMessage(upd.Edited_message)
				case upd.Channel_post:
					upd.Effective_message = u.bot.NewMessage(upd.Channel_post)
				case upd.Edited_channel_post:
					upd.Effective_message = u.bot.NewMessage(upd.Edited_channel_post)
				case upd.Inline_query:
					upd.Effective_message = u.bot.NewMessage(upd.Inline_query)
				}
				u.updates <- upd
			}
			if len(res) > 0 {
				offset = res[len(res)-1].Update_id + 1
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
