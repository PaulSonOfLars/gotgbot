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
	return h.Response(*d.Bot, u)
}

func (h Message) CheckUpdate(u *gotgbot.Update) (bool, error) {
	if u.EffectiveMessage == nil || u.CallbackQuery != nil { // don't trigger on callback query messages
		return false, nil
	}
	// if no edits and message is edited
	if !h.AllowEdited && u.EditedMessage != nil {
		return false, nil
	}
	// if no channel and message is channel message
	if !h.AllowChannel && u.ChannelPost != nil {
		return false, nil
	}
	// if no channel, no edits, and post is edited
	if !h.AllowChannel && !h.AllowEdited && u.EditedChannelPost != nil {
		return false, nil
	}

	return h.Filter(u.EffectiveMessage), nil
}
