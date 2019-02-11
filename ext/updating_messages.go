package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/pkg/errors"
)

func (b Bot) EditMessageText(chatId int, messageId int, text string) (*Message, error) {
	return b.EditMessageMarkup(chatId, messageId, text, "", nil)
}

func (b Bot) EditMessageHTML(chatId int, messageId int, text string) (*Message, error) {
	return b.EditMessageMarkup(chatId, messageId, text, parsemode.Html, nil)
}

func (b Bot) EditMessageMarkdown(chatId int, messageId int, text string) (*Message, error) {
	return b.EditMessageMarkup(chatId, messageId, text, parsemode.Markdown, nil)
}

func (b Bot) EditMessage(chatId int, messageId int, text string, parseMode string) (*Message, error) {
	return b.EditMessageMarkup(chatId, messageId, text, parseMode, nil)
}

func (b Bot) EditMessageMarkup(chatId int, messageId int, text string, parseMode string, markup ReplyMarkup) (*Message, error) {
	msg := b.NewSendableEditMessageText(chatId, messageId, text)
	msg.ParseMode = parseMode
	msg.ReplyMarkup = markup
	return msg.Send()
}

func (b Bot) EditMessageTextInline(inlineMessageId string, text string) (*Message, error) {
	return b.EditMessageInline(inlineMessageId, text, "")
}

func (b Bot) EditMessageHTMLInline(inlineMessageId string, text string) (*Message, error) {
	return b.EditMessageInline(inlineMessageId, text, parsemode.Html)
}

func (b Bot) EditMessageMarkdownInline(inlineMessageId string, text string) (*Message, error) {
	return b.EditMessageInline(inlineMessageId, text, parsemode.Markdown)
}

func (b Bot) EditMessageInline(inlineMessageId string, text string, parseMode string) (*Message, error) {
	msg := b.NewSendableEditMessageText(0, 0, text)
	msg.InlineMessageId = inlineMessageId
	msg.ParseMode = parseMode
	return msg.Send()
}

func (b Bot) EditMessageCaption(chatId int, messageId int, caption string) (*Message, error) {
	msg := b.NewSendableEditMessageCaption(chatId, messageId, caption)
	return msg.Send()
}

func (b Bot) EditMessageCaptionInline(inlineMessageId string, caption string) (*Message, error) {
	msg := b.NewSendableEditMessageCaption(0, 0, caption)
	msg.InlineMessageId = inlineMessageId
	return msg.Send()
}

func (b Bot) EditMessageReplyMarkup(chatId int, messageId int, replyMarkup InlineKeyboardMarkup) (*Message, error) {
	msg := b.NewSendableEditMessageReplyMarkup(chatId, messageId, &replyMarkup)
	return msg.Send()
}

func (b Bot) EditMessageReplyMarkupInline(inlineMessageId string, replyMarkup InlineKeyboardMarkup) (*Message, error) {
	msg := b.NewSendableEditMessageReplyMarkup(0, 0, &replyMarkup)
	msg.InlineMessageId = inlineMessageId
	return msg.Send()
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
