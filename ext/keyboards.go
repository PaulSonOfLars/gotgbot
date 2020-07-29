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

type KeyboardButtonPollType struct {
	Type []string `json:"type"`
}

type KeyboardButton struct {
	Text            string                  `json:"text"`
	RequestContact  bool                    `json:"request_contact"`
	RequestLocation bool                    `json:"request_location"`
	RequestPoll     *KeyboardButtonPollType `json:"request_poll"`
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

type LoginUrl struct {
	Url                string `json:"url"`
	ForwardText        string `json:"forward_text"`
	BotUsername        string `json:"bot_username"`
	RequestWriteAccess bool   `json:"request_write_access"`
}

type InlineKeyboardButton struct {
	Text                         string    `json:"text"`
	Url                          string    `json:"url"`
	LoginUrl                     *LoginUrl `json:"login_url"`
	CallbackData                 string    `json:"callback_data"`
	SwitchInlineQuery            string    `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string    `json:"switch_inline_query_current_chat"`
	// Callback_game                    *CallbackGame
	Pay bool `json:"pay"`
}

type CallbackQuery struct {
	Bot             Bot      `json:"-"`
	Id              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageId string   `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
	GameShortName   string   `json:"game_short_name"`
}

func (cq *CallbackQuery) AnswerCallbackQuery() (bool, error) {
	return cq.Bot.AnswerCallbackQuery(cq.Id)
}

func (cq *CallbackQuery) AnswerCallbackQueryText(text string, alert bool) (bool, error) {
	return cq.Bot.AnswerCallbackQueryText(cq.Id, text, alert)
}

func (cq *CallbackQuery) AnswerCallbackQueryURL(url string) (bool, error) {
	return cq.Bot.AnswerCallbackQueryURL(cq.Id, url)
}

type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"selective"`
}

func (fr *ForceReply) Marshal() ([]byte, error) {
	if fr == nil {
		fr = &ForceReply{
			ForceReply: false,
			Selective:  false,
		}
	}

	forceReply, err := json.Marshal(fr)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal force reply")
	}
	return forceReply, nil
}
