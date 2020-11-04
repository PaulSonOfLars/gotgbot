package ext

import (
	"encoding/json"
	"net/url"
	"strconv"
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
	return b.UnbanChatMemberOnlyIfBanned(chatId, userId, false)
}

func (b Bot) UnbanChatMemberOnlyIfBanned(chatId int, userId int, onlyIfBanned bool) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))
	v.Add("only_if_banned", strconv.FormatBool(onlyIfBanned))

	r, err := b.Get("unbanChatMember", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
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
	unRestrict.Permissions.CanChangeInfo = &temp
	unRestrict.Permissions.CanInviteUsers = &temp
	unRestrict.Permissions.CanPinMessages = &temp
	return unRestrict.Send()
}

func (b Bot) PromoteChatMember(chatId int, userId int) (bool, error) {
	promote := b.NewSendablePromoteChatMember(chatId, userId)
	return promote.Send()
}

func (b Bot) SetChatAdministratorCustomTitle(chatId int, userId int, customTitle string) (bool, error) {
	setTitle := b.NewSendableSetChatAdministratorCustomTitle(chatId, userId, customTitle)
	return setTitle.Send()
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

	r, err := b.Get("exportChatInviteLink", v)
	if err != nil {
		return "", err
	}

	var s string
	return s, json.Unmarshal(r, &s)
}

func (b Bot) SetChatPhoto(chatId int, photo InputFile) (bool, error) {
	setChatPhoto := b.NewSendableSetChatPhoto(chatId)
	setChatPhoto.Photo = photo
	return setChatPhoto.Send()
}

func (b Bot) DeleteChatPhoto(chatId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r, err := b.Get("deleteChatPhoto", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

func (b Bot) SetChatTitle(chatId int, title string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("title", title)

	r, err := b.Get("setChatTitle", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

func (b Bot) SetChatDescription(chatId int, description string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("description", description)

	r, err := b.Get("setChatDescription", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
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

	r, err := b.Get("unpinChatMessage", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

func (b Bot) LeaveChat(chatId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r, err := b.Get("leaveChat", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

func (b Bot) GetChat(chatId int) (*Chat, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r, err := b.Get("getChat", v)
	if err != nil {
		return nil, err
	}

	c := Chat{Bot: b}
	return &c, json.Unmarshal(r, &c)
}

func (b Bot) GetChatAdministrators(chatId int) ([]ChatMember, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r, err := b.Get("getChatAdministrators", v)
	if err != nil {
		return nil, err
	}

	var cm []ChatMember
	return cm, json.Unmarshal(r, &cm)
}

func (b Bot) GetChatMembersCount(chatId int) (int, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r, err := b.Get("getChatMembersCount", v)
	if err != nil {
		return 0, err
	}

	var c int
	return c, json.Unmarshal(r, &c)
}

func (b Bot) GetChatMember(chatId int, userId int) (*ChatMember, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("user_id", strconv.Itoa(userId))

	r, err := b.Get("getChatMember", v)
	if err != nil {
		return nil, err
	}

	var cm ChatMember
	return &cm, json.Unmarshal(r, &cm)
}

func (b Bot) SetChatStickerSet(chatId int, stickerSetName string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("sticker_set_name", stickerSetName)

	r, err := b.Get("setChatStickerSet", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

func (b Bot) DeleteChatStickerSet(chatId int) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))

	r, err := b.Get("deleteChatStickerSet", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}
