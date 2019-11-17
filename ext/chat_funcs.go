package ext

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

func (b Bot) KickChatMember(chatId int, userId int) (bool, error) {
	kick := b.NewSendableKickChatMember(chatId, userId)
	return kick.Send()
}

func (b Bot) KickChatMemberUntil(chatId int, userId int, untilDate int64) (bool, error) {
	kick := b.NewSendableKickChatMember(chatId, userId)
	kick.UntilDate = untilDate
	return kick.Send()
}

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

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

func (b Bot) RestrictChatMember(chatId int, userId int) (bool, error) {
	restrict := b.NewSendableRestrictChatMember(chatId, userId)
	return restrict.Send()
}

func (b Bot) UnRestrictChatMember(chatId int, userId int) (bool, error) {
	unRestrict := b.NewSendableRestrictChatMember(chatId, userId)
	temp := true
	unRestrict.Permissions.CanSendMessages = &temp
	unRestrict.Permissions.CanSendMediaMessages = &temp
	unRestrict.Permissions.CanSendPolls = &temp
	unRestrict.Permissions.CanSendOtherMessages = &temp
	unRestrict.Permissions.CanAddWebPagePreviews = &temp
	unRestrict.Permissions.CanInviteUsers = &temp
	return unRestrict.Send()
}

func (b Bot) PromoteChatMember(chatId int, userId int) (bool, error) {
	promote := b.NewSendablePromoteChatMember(chatId, userId)
	return promote.Send()
}

func (b Bot) DemoteChatMember(chatId int, userId int) (bool, error) {
	demote := b.NewSendablePromoteChatMember(chatId, userId)
	demote.CanChangeInfo = false
	demote.CanPostMessages = false
	demote.CanEditMessages = false
	demote.CanDeleteMessages = false
	demote.CanInviteUsers = false
	demote.CanRestrictMembers = false
	demote.CanPinMessages = false
	demote.CanPromoteMembers = false
	return demote.Send()
}

func (b Bot) SetChatPermissions(chatId int, perms ChatPermissions) (bool, error) {
	setChatPerms := b.NewSendableSetChatPermissions(chatId, perms)
	return setChatPerms.Send()
}

func (b Bot) ExportChatInviteLink(chatId int) (string, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r, err := Get(b, "exportChatInviteLink", v)
	if err != nil {
		return "", errors.Wrapf(err, "unable to exportChatInviteLink")
	}
	if !r.Ok {
		return "", errors.New(r.Description)
	}

	var s string
	return s, json.Unmarshal(r.Result, &s)
}

func (b Bot) SetChatPhotoStr(chatId int, photoId string) (bool, error) {
	setChatPhoto := b.NewSendableSetChatPhoto(chatId)
	setChatPhoto.FileId = photoId
	return setChatPhoto.Send()
}

func (b Bot) SetChatPhotoPath(chatId int, path string) (bool, error) {
	setChatPhoto := b.NewSendableSetChatPhoto(chatId)
	setChatPhoto.Path = path
	return setChatPhoto.Send()
}

func (b Bot) SetChatPhotoReader(chatId int, reader io.Reader) (bool, error) {
	setChatPhoto := b.NewSendableSetChatPhoto(chatId)
	setChatPhoto.Reader = reader
	return setChatPhoto.Send()
}

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
	return bb, json.Unmarshal(r.Result, &bb)
}

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
	return bb, json.Unmarshal(r.Result, &bb)
}

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
	return bb, json.Unmarshal(r.Result, &bb)
}

func (b Bot) PinChatMessage(chatId int, messageId int) (bool, error) {
	pin := b.NewSendablePinChatMessage(chatId, messageId)
	return pin.Send()
}

func (b Bot) PinChatMessageQuiet(chatId int, messageId int) (bool, error) {
	pin := b.NewSendablePinChatMessage(chatId, messageId)
	pin.DisableNotification = true
	return pin.Send()
}

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
	return bb, json.Unmarshal(r.Result, &bb)
}

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
	return bb, json.Unmarshal(r.Result, &bb)
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
	c.Bot = b
	return &c, json.Unmarshal(r.Result, &c)
}

func (b Bot) GetChatAdministrators(chatId int) ([]ChatMember, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r, err := Get(b, "getChatAdministrators", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to getChatAdministrators")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	var cm []ChatMember
	return cm, json.Unmarshal(r.Result, &cm)
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
	return c, json.Unmarshal(r.Result, &c)
}

func (b Bot) GetChatMember(chatId int, userId int) (*ChatMember, error) {
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

	var cm ChatMember
	return &cm, json.Unmarshal(r.Result, &cm)
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
	return bb, json.Unmarshal(r.Result, &bb)
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
	return bb, json.Unmarshal(r.Result, &bb)
}
