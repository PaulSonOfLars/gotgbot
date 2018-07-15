package ext

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type ReplyMarkup interface {
	Marshal() ([]byte, error)
}

type ReplyKeyboardMarkup struct {
	Keyboard        *[][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool                `json:"resize_keyboard"`
	OneTimeKeyboard bool                `json:"one_time_keyboard"`
	Selective       bool                `json:"selective"`
}

func (rkm *ReplyKeyboardMarkup) Marshal() ([]byte, error) {
	if rkm == nil {
		rkm = &ReplyKeyboardMarkup{
			Keyboard:        &[][]KeyboardButton{},
			ResizeKeyboard:  false,
			OneTimeKeyboard: false,
			Selective:       false,
		}
	} else if rkm.Keyboard == nil {
		rkm.Keyboard = &[][]KeyboardButton{}
	}

	replyMarkup, err := json.Marshal(rkm)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal reply markup")
	}
	return replyMarkup, nil
}

type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

func (rkm *ReplyKeyboardRemove) Marshal() ([]byte, error) {
	if rkm == nil {
		rkm = &ReplyKeyboardRemove{
			RemoveKeyboard: false,
			Selective:      false,
		}
	}

	kbRemove, err := json.Marshal(rkm)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal keyboard removal")
	}
	return kbRemove, nil
}

type InlineKeyboardMarkup struct {
	InlineKeyboard *[][]InlineKeyboardButton `json:"inline_keyboard"`
}

func (rkm *InlineKeyboardMarkup) Marshal() ([]byte, error) {
	if rkm == nil {
		rkm = &InlineKeyboardMarkup{
			InlineKeyboard: &[][]InlineKeyboardButton{},
		}
	} else if rkm.InlineKeyboard == nil {
		rkm.InlineKeyboard = &[][]InlineKeyboardButton{}
	}

	inlineKBMarkup, err := json.Marshal(rkm)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal inline keyboard")
	}
	return inlineKBMarkup, nil
}

type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	Url                          string `json:"url"`
	CallbackData                 string `json:"callback_data"`
	SwitchInlineQuery            string `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat"`
	//Callback_game                    *CallbackGame
	Pay bool `json:"pay"`
}

type CallbackQuery struct {
	Id              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageId string   `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
	GameShortName   string   `json:"game_short_name"`
}

type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"selective"`
}

func (rkm *ForceReply) Marshal() ([]byte, error) {
	if rkm == nil {
		rkm = &ForceReply{
			ForceReply: false,
			Selective:  false,
		}
	}

	forceReply, err := json.Marshal(rkm)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal force reply")
	}
	return forceReply, nil
}
