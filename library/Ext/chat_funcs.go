package Ext

import (
	"log"
	"bot/library/Types"
	"encoding/json"
	"strconv"
	"net/url"
)

// TODO: r.OK or unmarshal??
func (b Bot) KickChatMember(chat_id int, user_id int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("user_id", strconv.Itoa(user_id))

	r := Get(b, "kickChatMember", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for kickChatMember was not OK")
	}

	return r.Ok

}

// TODO: r.OK or unmarshal??
func (b Bot) UnbanChatMember(chat_id int, user_id int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("user_id", strconv.Itoa(user_id))

	r := Get(b, "unbanChatMember", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for unbanChatMember was not OK")
	}

	return r.Ok

}

// TODO: Add a nice way for all the methods
// TODO: r.OK or unmarshal??
func (b Bot) RestrictChatMember(chat_id int, user_id int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("user_id", strconv.Itoa(user_id))

	r := Get(b, "restrictChatMember", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for restrictChatMember was not OK")
	}

	return r.Ok

}

// TODO: Add a nice way for all the methods
// TODO: r.OK or unmarshal??
func (b Bot) PromoteChatMember(chat_id int, user_id int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("user_id", strconv.Itoa(user_id))

	r := Get(b, "promoteChatMember", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for promoteChatMember was not OK")
	}

	return r.Ok

}

func (b Bot) ExportChatLink(chat_id int) string {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))

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
//func (b Ext) SetChatPhoto(chat_id int, photo Types.InputFile) bool {
//	v := api_url.Values{}
//	v.Add("chat_id", strconv.Itoa(chat_id))
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
func (b Bot) DeleteChatPhoto(chat_id int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))

	r := Get(b, "deleteChatPhoto", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for deleteChatPhoto was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb

}

// TODO: r.OK or unmarshal??
func (b Bot) SetChatTitle(chat_id int, title string) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
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
func (b Bot) SetChatDescription(chat_id int, description string) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
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
func (b Bot) PinChatMessage(chat_id int, message_id int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("message_id", strconv.Itoa(message_id))

	r := Get(b, "pinChatMessage", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for pinChatMessage was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: r.OK or unmarshal??
func (b Bot) UnpinChatMessage(chat_id int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))

	r := Get(b, "unpinChatMessage", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for unpinChatMessage was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: r.OK or unmarshal??
func (b Bot) LeaveChat(chat_id int) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))

	r := Get(b, "leaveChat", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for leaveChat was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) GetChat(chat_id int) Types.Chat {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))

	r := Get(b, "getChat", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChat was not OK")
	}

	var c Types.Chat
	json.Unmarshal(r.Result, &c)

	return c
}

func (b Bot) GetChatAdministrators(chat_id int) []Types.ChatMember {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))

	r := Get(b, "getChatAdministrators", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChatAdministrators was not OK")
	}

	var cm []Types.ChatMember
	json.Unmarshal(r.Result, &cm)

	return cm
}

func (b Bot) GetChatMembersCount(chat_id int) int {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))

	r := Get(b, "getChatMembersCount", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChatMembersCount was not OK")
	}

	var c int
	json.Unmarshal(r.Result, &c)

	return c
}

func (b Bot) GetChatMember(chat_id int, user_id int) Types.ChatMember {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("user_id", strconv.Itoa(user_id))

	r := Get(b, "getChatMember", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getChatMember was not OK")
	}

	var cm Types.ChatMember
	json.Unmarshal(r.Result, &cm)

	return cm
}