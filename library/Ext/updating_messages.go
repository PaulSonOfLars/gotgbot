package Ext

import (
	"strconv"
	"log"
	"encoding/json"
	"bot/library/Types"
	"net/url"
)

// TODO: Check return type
// TODO: set parsemode
func (b Bot) EditMessageText(chatId int, messageId int, text string) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))
	v.Add("text", text)

	r := Get(b, "editMessageText", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageText was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: Check return type
// TODO: set parsemode
func (b Bot) EditMessageTextInline(inlineMessageId string, text string) bool {
	v := url.Values{}
	v.Add("inline_message_id", inlineMessageId)
	v.Add("text", text)

	r := Get(b, "editMessageText", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageText was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: Check return type
func (b Bot) EditMessageCaption(chatId int, messageId int, caption string) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))
	v.Add("caption", caption)

	r := Get(b, "editMessageCaption", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageCaption was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: Check return type
func (b Bot) EditMessageCaptionInline(inlineMessageId string, caption string) bool {
	v := url.Values{}
	v.Add("inline_message_id", inlineMessageId)
	v.Add("caption", caption)

	r := Get(b, "editMessageCaption", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageCaptionInline was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: Check return
func (b Bot) EditMessageReplyMarkup(chatId int, messageId int, replyMarkup Types.InlineKeyboardMarkup) bool {
	markupStr, err := json.Marshal(replyMarkup)
	if err != nil {
		log.Println("could not unmarshal in editMessageReplyMarkup")
		log.Fatal(err)
	}
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))
	v.Add("reply_markup", string(markupStr))

	r := Get(b, "editMessageReplyMarkup", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageReplyMarkup was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: Check return
func (b Bot) EditMessageReplyMarkupInline(inlineMessageId string, replyMarkup Types.InlineKeyboardMarkup) bool {
	markupStr, err := json.Marshal(replyMarkup)
	if err != nil {
		log.Println("could not unmarshal in editMessageReplyMarkupInline")
		log.Fatal(err)
	}
	v := url.Values{}
	v.Add("inline_message_id", inlineMessageId)
	v.Add("reply_markup", string(markupStr))

	r := Get(b, "editMessageReplyMarkup", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageReplyMarkup was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: ensure not a private chat! cant delete in private chats.
func (b Bot) DeleteMessage(chatId int, messageId int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))

	r := Get(b, "deleteMessage", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for deleteMessage was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}
