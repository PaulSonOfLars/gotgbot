package ext

import (
	"encoding/json"
	"strconv"
	"net/url"
	"github.com/pkg/errors"
	"github.com/PaulSonOfLars/gotgbot/types"
)

// TODO: r.OK or unmarshal??
func (b Bot) KickChatMember(chatId int, userId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r, err := Get(b, "kickChatMember", v)
	if err != nil {
		return false, errors.Wrapf(err, "could not kickChatMember")
	}
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

	r, err := Get(b, "unbanChatMember", v)
	if err != nil {
		return false, errors.Wrapf(err, "could not unbanChatMember")
	}
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

	r, err := Get(b, "restrictChatMember", v)
	if err != nil {
		return false, errors.Wrapf(err, "could not restrictChatMember")
	}
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

	r, err := Get(b, "promoteChatMember", v)
	if err != nil {
		return false, errors.Wrapf(err, "could not promoteChatMember")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	return r.Ok, nil
}

func (b Bot) ExportChatLink(chatId int) (string, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r, err := Get(b, "exportChatLink", v)
	if err != nil {
		return "", errors.Wrapf(err, "unable to exportChatLink")
	}
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
//	r, err := Get(b, "setChatPhoto", v)
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

	r, err := Get(b, "deleteChatPhoto", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to deleteChatPhoto")
	}
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

	r, err := Get(b, "setChatTitle", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to setChatTitle")
	}
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

	r, err := Get(b, "setChatDescription", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to setChatDescription")
	}
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

	r, err := Get(b, "pinChatMessage", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to pinChatMessage")
	}
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

	r, err := Get(b, "unpinChatMessage", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to unpinChatMessage")
	}
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

	r, err := Get(b, "leaveChat", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to leaveChat")
	}
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

	r, err := Get(b, "getChat", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to getChat")
	}
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

	r, err := Get(b, "getChatAdministrators", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to getChatAdministrators")
	}
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

	r, err := Get(b, "getChatMembersCount", v)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to getChatMembersCount")
	}
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

	r, err := Get(b, "getChatMember", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to getChatMember")
	}
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

	r, err := Get(b, "setChatStickerSet", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to setChatStickerSet")
	}
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

	r, err := Get(b, "deleteChatStickerSet", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to deleteChatStickerSet")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}