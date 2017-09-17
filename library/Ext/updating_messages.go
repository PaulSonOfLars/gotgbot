package Ext

import (
	"strconv"
	"log"
	"encoding/json"
	"bot/library/Types"
)

// TODO: Check return type
// TODO: set parsemode
func (b Bot) EditMessageText(chat_id int, message_id int, text string) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["message_id"] = strconv.Itoa(message_id)
	m["text"] = text

	r := Get(b, "editMessageText", m)
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
	m := make(map[string]string)
	m["inline_message_id"] = inline_message_id
	m["text"] = text

	r := Get(b, "editMessageText", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageText was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: Check return type
func (b Bot) EditMessageCaption(chat_id int, message_id int, caption string) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["message_id"] = strconv.Itoa(message_id)
	m["caption"] = caption

	r := Get(b, "editMessageCaption", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageCaption was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: Check return type
func (b Bot) EditMessageCaptionInline(inline_message_id string, caption string) bool {
	m := make(map[string]string)
	m["inline_message_id"] = inline_message_id
	m["caption"] = caption

	r := Get(b, "editMessageCaption", m)
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
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["message_id"] = strconv.Itoa(message_id)
	m["reply_markup"] = string(markup_str)

	r := Get(b, "editMessageReplyMarkup", m)
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
	m := make(map[string]string)
	m["inline_message_id"] = inline_message_id
	m["reply_markup"] = string(markup_str)

	r := Get(b, "editMessageReplyMarkup", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for editMessageReplyMarkup was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) DeleteMessage(chat_id int, message_id int) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["message_id"] = strconv.Itoa(message_id)

	r := Get(b, "deleteMessage", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for deleteMessage was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}
