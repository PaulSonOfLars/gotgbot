package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type sendableKickChatMember struct {
	bot       Bot `json:"-"`
	ChatId    int
	UserId    int
	UntilDate int64
}

func (b Bot) NewSendableKickChatMember(chatId int, userId int) *sendableKickChatMember {
	return &sendableKickChatMember{
		bot:    b,
		ChatId: chatId,
		UserId: userId,
	}
}

func (kcm *sendableKickChatMember) Send() (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(kcm.ChatId))
	v.Add("user_id", strconv.Itoa(kcm.UserId))
	v.Add("until_date", strconv.FormatInt(kcm.UntilDate, 10))

	r, err := Get(kcm.bot, "kickChatMember", v)
	if err != nil {
		return false, errors.Wrapf(err, "could not kickChatMember")
	}

	if !r.Ok {
		return false, errors.New(r.Description)
	}
	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

type sendableRestrictChatMember struct {
	bot         Bot `json:"-"`
	ChatId      int
	UserId      int
	Permissions ChatPermissions
	UntilDate   int64
}

func (b Bot) NewSendableRestrictChatMember(chatId int, userId int) *sendableRestrictChatMember {
	temp := false
	return &sendableRestrictChatMember{
		bot:       b,
		ChatId:    chatId,
		UserId:    userId,
		UntilDate: 0,
		Permissions: ChatPermissions{
			CanSendMessages:       &temp,
			CanSendMediaMessages:  &temp,
			CanSendPolls:          &temp,
			CanSendOtherMessages:  &temp,
			CanAddWebPagePreviews: &temp,
			CanChangeInfo:         &temp,
			CanInviteUsers:        &temp,
			CanPinMessages:        &temp,
		},
	}
}

func (rcm *sendableRestrictChatMember) Send() (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(rcm.ChatId))
	v.Add("user_id", strconv.Itoa(rcm.UserId))
	v.Add("until_date", strconv.FormatInt(rcm.UntilDate, 10))

	perms, err := json.Marshal(&rcm.Permissions)
	if err != nil {
		return false, errors.Wrapf(err, "could not marshal permissions")
	}

	v.Add("permissions", string(perms))

	r, err := Get(rcm.bot, "restrictChatMember", v)
	if err != nil {
		return false, errors.Wrapf(err, "could not restrictChatMember")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

type sendablePromoteChatMember struct {
	bot                Bot `json:"-"`
	ChatId             int
	UserId             int
	CanChangeInfo      bool
	CanPostMessages    bool
	CanEditMessages    bool
	CanDeleteMessages  bool
	CanInviteUsers     bool
	CanRestrictMembers bool
	CanPinMessages     bool
	CanPromoteMembers  bool
}

// note: set all as true for promotion by default
func (b Bot) NewSendablePromoteChatMember(chatId int, userId int) *sendablePromoteChatMember {
	return &sendablePromoteChatMember{
		bot:                b,
		ChatId:             chatId,
		UserId:             userId,
		CanChangeInfo:      true,
		CanPostMessages:    true,
		CanEditMessages:    true,
		CanDeleteMessages:  true,
		CanInviteUsers:     true,
		CanRestrictMembers: true,
		CanPinMessages:     true,
		CanPromoteMembers:  true,
	}
}

func (rcm *sendablePromoteChatMember) Send() (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(rcm.ChatId))
	v.Add("user_id", strconv.Itoa(rcm.UserId))
	v.Add("can_change_info", strconv.FormatBool(rcm.CanChangeInfo))
	v.Add("can_post_messages", strconv.FormatBool(rcm.CanPostMessages))
	v.Add("can_edit_messages", strconv.FormatBool(rcm.CanEditMessages))
	v.Add("can_delete_messages", strconv.FormatBool(rcm.CanDeleteMessages))
	v.Add("can_invite_users", strconv.FormatBool(rcm.CanInviteUsers))
	v.Add("can_restrict_members", strconv.FormatBool(rcm.CanRestrictMembers))
	v.Add("can_pin_messages", strconv.FormatBool(rcm.CanPinMessages))
	v.Add("can_promote_members", strconv.FormatBool(rcm.CanPromoteMembers))

	r, err := Get(rcm.bot, "promoteChatMember", v)
	if err != nil {
		return false, errors.Wrapf(err, "could not promoteChatMember")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

type sendableSetChatPermissions struct {
	bot         Bot `json:"-"`
	ChatId      int
	Permissions ChatPermissions
}

// note: set all as true for promotion by default
func (b Bot) NewSendableSetChatPermissions(chatId int, perms ChatPermissions) *sendableSetChatPermissions {
	return &sendableSetChatPermissions{
		bot:         b,
		ChatId:      chatId,
		Permissions: perms,
	}
}

func (scp *sendableSetChatPermissions) Send() (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(scp.ChatId))

	perms, err := json.Marshal(&scp.Permissions)
	if err != nil {
		return false, errors.Wrapf(err, "could not marshal permissions")
	}

	v.Add("permissions", string(perms))

	r, err := Get(scp.bot, "setChatPermissions", v)
	if err != nil {
		return false, errors.Wrapf(err, "could not setChatPermissions")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

type sendablePinChatMessage struct {
	bot                 Bot `json:"-"`
	ChatId              int
	MessageId           int
	DisableNotification bool
}

func (b Bot) NewSendablePinChatMessage(chatId int, messageId int) *sendablePinChatMessage {
	return &sendablePinChatMessage{
		bot:                 b,
		ChatId:              chatId,
		MessageId:           messageId,
		DisableNotification: false,
	}
}

func (pcm *sendablePinChatMessage) Send() (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(pcm.ChatId))
	v.Add("message_id", strconv.Itoa(pcm.MessageId))
	v.Add("disable_notification", strconv.FormatBool(pcm.DisableNotification))

	r, err := Get(pcm.bot, "pinChatMessage", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to pinChatMessage")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

type sendableSetChatPhoto struct {
	bot    Bot `json:"-"`
	ChatId int
	file
}

func (b Bot) NewSendableSetChatPhoto(chatId int) *sendableSetChatPhoto {
	return &sendableSetChatPhoto{bot: b, ChatId: chatId}
}

func (scp *sendableSetChatPhoto) Send() (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(scp.ChatId))

	r, err := scp.bot.sendFile(scp.file, "photo", "setChatPhoto", v)

	if err != nil {
		return false, errors.Wrapf(err, "unable to setChatPhoto")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}
	var newMsg bool
	return newMsg, json.Unmarshal(r.Result, newMsg)
}
