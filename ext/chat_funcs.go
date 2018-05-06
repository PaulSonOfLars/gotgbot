package ext

import (
	"gotgbot/types"
	"encoding/json"
	"strconv"
	"net/url"
	"github.com/pkg/errors"
)

// TODO: r.OK or unmarshal??
func (b Bot) KickChatMember(chatId int, userId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r := Get(b, "kickChatMember", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	return r.Ok, nil
}

// TODO: r.OK or unmarshal??
func (b Bot) UnbanChatMember(chatId int, userId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r := Get(b, "unbanChatMember", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	return r.Ok, nil
}

// TODO: Add a nice way for all the methods
// TODO: r.OK or unmarshal??
func (b Bot) RestrictChatMember(chatId int, userId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r := Get(b, "restrictChatMember", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	return r.Ok, nil
}

// TODO: Add a nice way for all the methods
// TODO: r.OK or unmarshal??
func (b Bot) PromoteChatMember(chatId int, userId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r := Get(b, "promoteChatMember", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	return r.Ok, nil
}

func (b Bot) ExportChatLink(chatId int) (string, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "exportChatLink", v)
	if !r.Ok {
		return "", errors.New(r.Description)
	}

	var s string
	json.Unmarshal(r.Result, &s)

	return s, nil
}

// TODO: figure out InputFiles
// TODO: r.OK or unmarshal??
//func (b ext) SetChatPhoto(chatId int, photo types.InputFile) (bool, error) {
//	v := api_url.Values{}
//	v.Add("chat_id", strconv.Itoa(chatId))
//	v.Add("photo", photo)
//
//	r := Get(b, "setChatPhoto", v)
//	if !r.Ok {
//		return false, errors.New(r.Description)
//	}
//
//	var bb bool
//	json.Unmarshal(r.Result, &bb)
//
//	return bb, nil
//}

// TODO: r.OK or unmarshal??
func (b Bot) DeleteChatPhoto(chatId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "deleteChatPhoto", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: r.OK or unmarshal??
func (b Bot) SetChatTitle(chatId int, title string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("title", title)

	r := Get(b, "setChatTitle", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: r.OK or unmarshal??
func (b Bot) SetChatDescription(chatId int, description string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("description", description)

	r := Get(b, "setChatDescription", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: r.OK or unmarshal??
func (b Bot) PinChatMessage(chatId int, messageId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))

	r := Get(b, "pinChatMessage", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: r.OK or unmarshal??
func (b Bot) UnpinChatMessage(chatId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "unpinChatMessage", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: r.OK or unmarshal??
func (b Bot) LeaveChat(chatId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "leaveChat", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

func (b Bot) GetChat(chatId int) (*Chat, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "getChat", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	var c Chat
	json.Unmarshal(r.Result, &c)

	return &c, nil
}

func (b Bot) GetChatAdministrators(chatId int) ([]types.ChatMember, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "getChatAdministrators", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	var cm []types.ChatMember
	json.Unmarshal(r.Result, &cm)

	return cm, nil
}

func (b Bot) GetChatMembersCount(chatId int) (int, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "getChatMembersCount", v)
	if !r.Ok {
		return 0, errors.New(r.Description)
	}

	var c int
	json.Unmarshal(r.Result, &c)

	return c, nil
}

func (b Bot) GetChatMember(chatId int, userId int) (*types.ChatMember, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r := Get(b, "getChatMember", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	var cm types.ChatMember
	json.Unmarshal(r.Result, &cm)

	return &cm, nil
}

func (b Bot) SetChatStickerSet(chatId int, stickerSetName string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("sticker_set_name", stickerSetName)

	r := Get(b, "setChatStickerSet", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

func (b Bot) DeleteChatStickerSet(chatId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r := Get(b, "deleteChatStickerSet", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}