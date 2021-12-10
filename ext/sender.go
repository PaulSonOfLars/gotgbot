package ext

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Sender struct {
	User *gotgbot.User
	Chat *gotgbot.Chat

	IsAutomaticForward bool
	ChatId             int64
}

func GetSender(msg *gotgbot.Message) *Sender {
	return &Sender{
		User:               msg.From,
		Chat:               msg.SenderChat,
		IsAutomaticForward: msg.IsAutomaticForward,
		ChatId:             msg.Chat.Id,
	}
}

func (s Sender) Id() int64 {
	if s.Chat != nil {
		return s.Chat.Id
	}
	if s.User != nil {
		return s.User.Id
	}
	return 0
}

func (s Sender) Username() string {
	if s.Chat != nil {
		return s.Chat.Username
	}
	if s.User != nil {
		return s.User.Username
	}
	return ""
}

func (s Sender) IsUser() bool {
	return s.Chat == nil && s.User != nil
}

func (s Sender) IsBot() bool {
	return s.Chat == nil && s.User != nil && s.User.IsBot
}

func (s Sender) IsAnonymousAdmin() bool {
	return s.Chat != nil && s.Chat.Id == s.ChatId
}

func (s Sender) IsAnonymousChannel() bool {
	return s.Chat != nil && s.Chat.Id != s.ChatId && !s.IsAutomaticForward
}

func (s Sender) IsLinkedChannel() bool {
	return s.Chat != nil && s.Chat.Id != s.ChatId && s.IsAutomaticForward
}
