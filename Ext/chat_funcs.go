package Ext

import (
	"log"
	"gotgbot/Types"
	"encoding/json"
	"strconv"
	"net/url"
)

// TODO: r.OK or unmarshal??
func (b Bot) KickChatMember(chatId int, userId int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r := Get(b, "kickChatMember", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for kickChatMember was not OK")
	}

	return r.Ok

}

// TODO: r.OK or unmarshal??
func (b Bot) UnbanChatMember(chatId int, userId int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r := Get(b, "unbanChatMember", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for unbanChatMember was not OK")
	}

	return r.Ok

}

// TODO: Add a nice way for all the methods
// TODO: r.OK or unmarshal??
func (b Bot) RestrictChatMember(chatId int, userId int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r := Get(b, "restrictChatMember", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for restrictChatMember was not OK")
	}

	return r.Ok

}

// TODO: Add a nice way for all the methods
// TODO: r.OK or unmarshal??
func (b Bot) PromoteChatMember(chatId int, userId int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r := Get(b, "promoteChatMember", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for promoteChatMember was not OK")
	}

	return r.Ok

}

func (b Bot) ExportChatLink(chatId int) string {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "exportChatLink", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for exportChatLink was not OK")
	}

	var s string
	json.Unmarshal(r.Result, &s)

	return s

}

// TODO: figure out InputFiles
// TODO: r.OK or unmarshal??
//func (b Ext) SetChatPhoto(chatId int, photo Types.InputFile) bool {
//	v := api_url.Values{}
//	v.Add("chat_id", strconv.Itoa(chatId))
//	v.Add("photo", photo)
//
//	r := Get(b, "setChatPhoto", v)
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
func (b Bot) DeleteChatPhoto(chatId int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "deleteChatPhoto", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for deleteChatPhoto was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb

}

// TODO: r.OK or unmarshal??
func (b Bot) SetChatTitle(chatId int, title string) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("title", title)

	r := Get(b, "setChatTitle", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for setChatTitle was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: r.OK or unmarshal??
func (b Bot) SetChatDescription(chatId int, description string) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("description", description)

	r := Get(b, "setChatDescription", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for setChatDescription was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: r.OK or unmarshal??
func (b Bot) PinChatMessage(chatId int, messageId int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))

	r := Get(b, "pinChatMessage", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for pinChatMessage was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: r.OK or unmarshal??
func (b Bot) UnpinChatMessage(chatId int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "unpinChatMessage", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for unpinChatMessage was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: r.OK or unmarshal??
func (b Bot) LeaveChat(chatId int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "leaveChat", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for leaveChat was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) GetChat(chatId int) Types.Chat {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "getChat", v)
	if !r.Ok {
		log.Println(r)
		log.Fatal("You done goofed, API Res for getChat was not OK")
	}

	var c Types.Chat
	json.Unmarshal(r.Result, &c)

	return c
}

func (b Bot) GetChatAdministrators(chatId int) []Types.ChatMember {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "getChatAdministrators", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChatAdministrators was not OK")
	}

	var cm []Types.ChatMember
	json.Unmarshal(r.Result, &cm)

	return cm
}

func (b Bot) GetChatMembersCount(chatId int) int {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "getChatMembersCount", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChatMembersCount was not OK")
	}

	var c int
	json.Unmarshal(r.Result, &c)

	return c
}

func (b Bot) GetChatMember(chatId int, userId int) Types.ChatMember {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r := Get(b, "getChatMember", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChatMember was not OK")
	}

	var cm Types.ChatMember
	json.Unmarshal(r.Result, &cm)

	return cm
}