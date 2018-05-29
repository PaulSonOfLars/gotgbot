package ext

import (
	"net/url"
	"strconv"
	"github.com/pkg/errors"
	"encoding/json"
)

type sendableKickChatMember struct {
	bot       Bot
	ChatId    int
	UserId    int
	UntilDate int
}

func (b Bot) NewSendableKickChatMember(chatId int, userId int) *sendableKickChatMember {
	return &sendableKickChatMember{
		bot:       b,
		ChatId:    chatId,
		UserId:    userId,
	}
}

func (kcm *sendableKickChatMember) Send() (bool, error){
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(kcm.ChatId))
	v.Add("user_id", strconv.Itoa(kcm.UserId))
	v.Add("until_date", strconv.Itoa(kcm.UntilDate))

	r, err := Get(kcm.bot, "kickChatMember", v)
	if err != nil {
		return false, errors.Wrapf(err, "could not kickChatMember")
	}

	if !r.Ok {
		return false, errors.New(r.Description)
	}
	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

type sendableRestrictChatMember struct {
	bot                   Bot
	ChatId                int
	UserId                int
	UntilDate             int
	CanSendMessages       bool
	CanSendMediaMessages  bool
	CanSendOtherMessages  bool
	CanAddWebPagePreviews bool
}

func (b Bot) NewSendableRestrictChatMember(chatId int, userId int) *sendableRestrictChatMember {
	return &sendableRestrictChatMember{
		bot:                   b,
		ChatId:                chatId,
		UserId:                userId,
		UntilDate:             0,
		CanSendMessages:       false,
		CanSendMediaMessages:  false,
		CanSendOtherMessages:  false,
		CanAddWebPagePreviews: false,
	}
}

func (rcm *sendableRestrictChatMember) Send() (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(rcm.ChatId))
	v.Add("user_id", strconv.Itoa(rcm.UserId))
	v.Add("until_date", strconv.Itoa(rcm.UntilDate))
	v.Add("can_send_messages", strconv.FormatBool(rcm.CanSendMessages))
	v.Add("can_send_media_messages", strconv.FormatBool(rcm.CanSendMediaMessages))
	v.Add("can_send_other_messages", strconv.FormatBool(rcm.CanSendOtherMessages))
	v.Add("can_add_web_page_previews", strconv.FormatBool(rcm.CanAddWebPagePreviews))

	r, err := Get(rcm.bot, "restrictChatMember", v)
	if err != nil {
		return false, errors.Wrapf(err, "could not restrictChatMember")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

type sendablePromoteChatMember struct {
	bot                Bot
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
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

type sendablePinChatMessage struct {
	bot                 Bot
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
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}
