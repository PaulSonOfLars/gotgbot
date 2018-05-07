package ext

import (
	"strconv"
	"encoding/json"
	"gotgbot/types"
	"net/url"
	"github.com/pkg/errors"
)

// TODO: Check return type
// TODO: set parsemode
func (b Bot) EditMessageText(chatId int, messageId int, text string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))
	v.Add("text", text)

	r, err := Get(b, "editMessageText", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to editMessageText")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: Check return type
// TODO: set parsemode
func (b Bot) EditMessageTextInline(inlineMessageId string, text string) (bool, error) {
	v := url.Values{}
	v.Add("inline_message_id", inlineMessageId)
	v.Add("text", text)

	r, err := Get(b, "editMessageText", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to editMessageText")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: Check return type
func (b Bot) EditMessageCaption(chatId int, messageId int, caption string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))
	v.Add("caption", caption)

	r, err := Get(b, "editMessageCaption", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to editMessageCaption")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: Check return type
func (b Bot) EditMessageCaptionInline(inlineMessageId string, caption string) (bool, error) {
	v := url.Values{}
	v.Add("inline_message_id", inlineMessageId)
	v.Add("caption", caption)

	r, err := Get(b, "editMessageCaption", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to editMessageCaption")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: Check return
func (b Bot) EditMessageReplyMarkup(chatId int, messageId int, replyMarkup types.InlineKeyboardMarkup) (bool, error) {
	markupStr, err := json.Marshal(replyMarkup)
	if err != nil {
		return false, nil
	}
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))
	v.Add("reply_markup", string(markupStr))

	r, err := Get(b, "editMessageReplyMarkup", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to edit message")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: Check return
func (b Bot) EditMessageReplyMarkupInline(inlineMessageId string, replyMarkup types.InlineKeyboardMarkup) (bool, error) {
	markupStr, err := json.Marshal(replyMarkup)
	if err != nil {
		return false, errors.Wrapf(err, "error editing inline markup reply")
	}
	v := url.Values{}
	v.Add("inline_message_id", inlineMessageId)
	v.Add("reply_markup", string(markupStr))

	r, err := Get(b, "editMessageReplyMarkup", v)
	if err != nil {
		return false, errors.New(r.Description)
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: ensure not a private chat! cant delete in private chats.
func (b Bot) DeleteMessage(chatId int, messageId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))

	r, err := Get(b, "deleteMessage", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to deleteMessage")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}
