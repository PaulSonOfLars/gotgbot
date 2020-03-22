package handlers

import (
	"regexp"

	"github.com/pkg/errors"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type Regex struct {
	baseHandler
	AllowEdited  bool
	AllowChannel bool
	Match        string
	Response     func(b ext.Bot, u *gotgbot.Update) error
}

func NewRegex(match string, response func(b ext.Bot, u *gotgbot.Update) error) Regex {
	return Regex{
		baseHandler: baseHandler{
			Name: match,
		},
		Match:    match,
		Response: response,
	}
}

func (h Regex) HandleUpdate(u *gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.Response(*d.Bot, u)
}

func (h Regex) CheckUpdate(u *gotgbot.Update) (bool, error) {
	if u.EffectiveMessage == nil || u.CallbackQuery != nil {
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

	text := u.EffectiveMessage.Text
	if text == "" {
		text = u.EffectiveMessage.Caption
	}
	res, err := regexp.MatchString(h.Match, text)
	if err != nil {
		return false, errors.Wrapf(err, "Could not match regexp")
	}
	return res, nil
}
