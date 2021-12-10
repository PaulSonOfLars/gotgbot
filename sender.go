package gotgbot

// Sender is a merge of the User and SenderChat fields of a message, to provide easier interaction with
// message senders from the telegram API.
type Sender struct {
	User               *User
	Chat               *Chat
	IsAutomaticForward bool
	ChatId             int64
}

// GetSender populates the relevant fields of a Sender struct given a message.
func (m Message) GetSender() *Sender {
	return &Sender{
		User:               m.From,
		Chat:               m.SenderChat,
		IsAutomaticForward: m.IsAutomaticForward,
		ChatId:             m.Chat.Id,
	}
}

// Id determines the real sender ID.
func (s Sender) Id() int64 {
	if s.Chat != nil {
		return s.Chat.Id
	}
	if s.User != nil {
		return s.User.Id
	}
	return 0
}

// Username determines the real sender username.
func (s Sender) Username() string {
	if s.Chat != nil {
		return s.Chat.Username
	}
	if s.User != nil {
		return s.User.Username
	}
	return ""
}

// Name determines the real name of the sender.
// This is:
// - chat.Title for a chat.
// - user.Firstname + user.Lastname for a user (the full name).
func (s Sender) Name() string {
	if s.Chat != nil {
		return s.Chat.Title
	}
	if s.User != nil {
		if s.User.LastName == "" {
			return s.User.FirstName
		}
		return s.User.FirstName + " " + s.User.LastName
	}
	return ""
}

// IsUser returns true if the Sender is a user (including bot).
func (s Sender) IsUser() bool {
	return s.Chat == nil && s.User != nil
}

// IsBot returns true if the Sender is a bot.
// Returns false if the user is a bot setup by telegram for backwards compatibility with
// the sender_chat fields.
func (s Sender) IsBot() bool {
	return s.Chat == nil && s.User != nil && s.User.IsBot
}

// IsAnonymousAdmin returns true if the Sender is an anonymous admin.
func (s Sender) IsAnonymousAdmin() bool {
	return s.Chat != nil && s.Chat.Id == s.ChatId
}

// IsAnonymousChannel returns true if the Sender is an anonymous channel.
func (s Sender) IsAnonymousChannel() bool {
	return s.Chat != nil && s.Chat.Id != s.ChatId && !s.IsAutomaticForward
}

// IsLinkedChannel returns true if the Sender is a linked channel.
func (s Sender) IsLinkedChannel() bool {
	return s.Chat != nil && s.Chat.Id != s.ChatId && s.IsAutomaticForward
}
