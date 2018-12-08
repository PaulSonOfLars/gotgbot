package handlers

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type FilterFunc func(message *ext.Message) bool

type Message struct {
	baseHandler
	AllowEdited  bool
	AllowChannel bool
	Filter       FilterFunc
	Response     func(b ext.Bot, u *gotgbot.Update) error
}

func NewMessage(filterFunc FilterFunc, response func(b ext.Bot, u *gotgbot.Update) error) Message {
	return Message{
		baseHandler: baseHandler{
			Name: "unnamedMessageHandler",
		},
		AllowEdited:  false,
		AllowChannel: false,
		Filter:       filterFunc,
		Response:     response,
	}
}

func (h Message) HandleUpdate(u *gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.Response(d.Bot, u)
}

func (h Message) CheckUpdate(u *gotgbot.Update) (bool, error) {
	return (u.Message != nil ||
		(h.AllowEdited && u.EditedMessage != nil) ||
		(h.AllowChannel && u.ChannelPost != nil)) &&
		h.Filter(u.EffectiveMessage), nil
}
