package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/pkg/errors"
)

func (b Bot) EditMessageText(chatId int, messageId int, text string) (bool, error) {
	return b.EditMessage(chatId, messageId, text, "")
}

func (b Bot) EditMessageHTML(chatId int, messageId int, text string) (bool, error) {
	return b.EditMessage(chatId, messageId, text, parsemode.Html)
}

func (b Bot) EditMessageMarkdown(chatId int, messageId int, text string) (bool, error) {
	return b.EditMessage(chatId, messageId, text, parsemode.Markdown)
}

func (b Bot) EditMessage(chatId int, messageId int, text string, parseMode string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))
	v.Add("text", text)
	v.Add("parse_mode", parseMode)

	return b.boolSender("editMessageText", v)
}

func (b Bot) EditMessageTextInline(messageId int, text string) (bool, error) {
	return b.EditMessageInline(messageId, text, "")
}

func (b Bot) EditMessageHTMLInline(messageId int, text string) (bool, error) {
	return b.EditMessageInline(messageId, text, parsemode.Html)
}

func (b Bot) EditMessageMarkdownInline(messageId int, text string) (bool, error) {
	return b.EditMessageInline(messageId, text, parsemode.Markdown)
}

func (b Bot) EditMessageInline(messageId int, text string, parseMode string) (bool, error) {
	v := url.Values{}
	v.Add("inline_message_id", strconv.Itoa(messageId))
	v.Add("text", text)
	v.Add("parse_mode", parseMode)

	return b.boolSender("editMessageText", v)
}

func (b Bot) EditMessageCaption(chatId int, messageId int, caption string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))
	v.Add("caption", caption)

	return b.boolSender("editMessageCaption", v)
}

func (b Bot) EditMessageCaptionInline(inlineMessageId string, caption string) (bool, error) {
	v := url.Values{}
	v.Add("inline_message_id", inlineMessageId)
	v.Add("caption", caption)

	return b.boolSender("editMessageCaption", v)

}

func (b Bot) EditMessageReplyMarkup(chatId int, messageId int, replyMarkup InlineKeyboardMarkup) (bool, error) {
	markupStr, err := json.Marshal(replyMarkup)
	if err != nil {
		return false, nil
	}
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))
	v.Add("reply_markup", string(markupStr))

	return b.boolSender("editMessageReplyMarkup", v)
}

func (b Bot) EditMessageReplyMarkupInline(inlineMessageId string, replyMarkup InlineKeyboardMarkup) (bool, error) {
	markupStr, err := json.Marshal(replyMarkup)
	if err != nil {
		return false, errors.Wrapf(err, "error editing inline markup reply")
	}
	v := url.Values{}
	v.Add("inline_message_id", inlineMessageId)
	v.Add("reply_markup", string(markupStr))

	return b.boolSender("editMessageReplyMarkup", v)
}

// TODO: ensure not a private chat! cant delete in private chats.
func (b Bot) DeleteMessage(chatId int, messageId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))

	return b.boolSender("deleteMessage", v)
}

func (b Bot) boolSender(meth string, v url.Values) (bb bool, err error) {
	r, err := Get(b, meth, v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to complete request for %s", meth)
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	return bb, json.Unmarshal(r.Result, &bb)
}
