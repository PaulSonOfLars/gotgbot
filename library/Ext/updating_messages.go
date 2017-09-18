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
func (b Bot) EditMessageText(chat_id int, message_id int, text string) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("message_id", strconv.Itoa(message_id))
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
func (b Bot) EditMessageTextInline(inline_message_id string, text string) bool {
	v := url.Values{}
	v.Add("inline_message_id", inline_message_id)
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
func (b Bot) EditMessageCaption(chat_id int, message_id int, caption string) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("message_id", strconv.Itoa(message_id))
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
func (b Bot) EditMessageCaptionInline(inline_message_id string, caption string) bool {
	v := url.Values{}
	v.Add("inline_message_id", inline_message_id)
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
func (b Bot) EditMessageReplyMarkup(chat_id int, message_id int, reply_markup Types.InlineKeyboardMarkup) bool {
	markup_str, err := json.Marshal(reply_markup)
	if err != nil {
		log.Println("could not unmarshal in editMessageReplyMarkup")
		log.Fatal(err)
	}
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("message_id", strconv.Itoa(message_id))
	v.Add("reply_markup", string(markup_str))

	r := Get(b, "editMessageReplyMarkup", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageReplyMarkup was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: Check return
func (b Bot) EditMessageReplyMarkupInline(inline_message_id string, reply_markup Types.InlineKeyboardMarkup) bool {
	markup_str, err := json.Marshal(reply_markup)
	if err != nil {
		log.Println("could not unmarshal in editMessageReplyMarkupInline")
		log.Fatal(err)
	}
	v := url.Values{}
	v.Add("inline_message_id", inline_message_id)
	v.Add("reply_markup", string(markup_str))

	r := Get(b, "editMessageReplyMarkup", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageReplyMarkup was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: ensure not a private chat! cant delete in private chats.
func (b Bot) DeleteMessage(chat_id int, message_id int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("message_id", strconv.Itoa(message_id))

	r := Get(b, "deleteMessage", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for deleteMessage was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}
