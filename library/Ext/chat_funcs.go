package Ext

import (
	"log"
	"bot/library/Types"
	"encoding/json"
	"strconv"
)

// TODO: r.OK or unmarshal??
func (b Bot) KickChatMember(chat_id int, user_id int) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["user_id"] = strconv.Itoa(user_id)

	r := Get(b, "kickChatMember", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for kickChatMember was not OK")
	}

	return r.Ok

}

// TODO: r.OK or unmarshal??
func (b Bot) UnbanChatMember(chat_id int, user_id int) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["user_id"] = strconv.Itoa(user_id)

	r := Get(b, "unbanChatMember", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for unbanChatMember was not OK")
	}

	return r.Ok

}

// TODO: Add a nice way for all the methods
// TODO: r.OK or unmarshal??
func (b Bot) RestrictChatMember(chat_id int, user_id int) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["user_id"] = strconv.Itoa(user_id)

	r := Get(b, "restrictChatMember", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for restrictChatMember was not OK")
	}

	return r.Ok

}

// TODO: Add a nice way for all the methods
// TODO: r.OK or unmarshal??
func (b Bot) PromoteChatMember(chat_id int, user_id int) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["user_id"] = strconv.Itoa(user_id)

	r := Get(b, "promoteChatMember", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for promoteChatMember was not OK")
	}

	return r.Ok

}

func (b Bot) ExportChatLink(chat_id int) string {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)

	r := Get(b, "exportChatLink", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for exportChatLink was not OK")
	}

	var s string
	json.Unmarshal(r.Result, &s)

	return s

}

// TODO: figure out InputFiles
// TODO: r.OK or unmarshal??
//func (b Ext) SetChatPhoto(chat_id int, photo Types.InputFile) bool {
//	m := make(map[string]string)
//	m["chat_id"] = strconv.Itoa(chat_id)
//	m["photo"] = photo
//
//	r := Get(b, "setChatPhoto", m)
//	if !r.Ok {
//		log.Fatal("You done goofed, API Res for setChatPhoto was not OK")
//	}
//
//	var bb bool
//	json.Unmarshal(r.Result, &bb)
//
//	return bb
//
//}

// TODO: r.OK or unmarshal??
func (b Bot) DeleteChatPhoto(chat_id int) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)

	r := Get(b, "deleteChatPhoto", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for deleteChatPhoto was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb

}

// TODO: r.OK or unmarshal??
func (b Bot) SetChatTitle(chat_id int, title string) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["title"] = title

	r := Get(b, "setChatTitle", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for setChatTitle was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: r.OK or unmarshal??
func (b Bot) SetChatDescription(chat_id int, description string) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["description"] = description

	r := Get(b, "setChatDescription", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for setChatDescription was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: r.OK or unmarshal??
func (b Bot) PinChatMessage(chat_id int, message_id int) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["message_id"] = strconv.Itoa(message_id)

	r := Get(b, "pinChatMessage", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for pinChatMessage was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: r.OK or unmarshal??
func (b Bot) UnpinChatMessage(chat_id int) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)

	r := Get(b, "unpinChatMessage", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for unpinChatMessage was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: r.OK or unmarshal??
func (b Bot) LeaveChat(chat_id int) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)

	r := Get(b, "leaveChat", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for leaveChat was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) GetChat(chat_id int) Types.Chat {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)

	r := Get(b, "getChat", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChat was not OK")
	}

	var c Types.Chat
	json.Unmarshal(r.Result, &c)

	return c
}

func (b Bot) GetChatAdministrators(chat_id int) []Types.ChatMember {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)

	r := Get(b, "getChatAdministrators", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChatAdministrators was not OK")
	}

	var cm []Types.ChatMember
	json.Unmarshal(r.Result, &cm)

	return cm
}

func (b Bot) GetChatMembersCount(chat_id int) int {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)

	r := Get(b, "getChatMembersCount", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChatMembersCount was not OK")
	}

	var c int
	json.Unmarshal(r.Result, &c)

	return c
}

func (b Bot) GetChatMember(chat_id int, user_id int) Types.ChatMember {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["user_id"] = strconv.Itoa(user_id)

	r := Get(b, "getChatMember", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChatMember was not OK")
	}

	var cm Types.ChatMember
	json.Unmarshal(r.Result, &cm)

	return cm
}