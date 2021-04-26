package gotgbot

// Reply is a helper function to easily call Bot.SendMessage as a reply to an existing message.
func (m Message) Reply(b *Bot, text string, opts *SendMessageOpts) (*Message, error) {
	if opts == nil {
		opts = &SendMessageOpts{}
	}

	if opts.ReplyToMessageId == 0 {
		opts.ReplyToMessageId = m.MessageId
	}

	return b.SendMessage(m.Chat.Id, text, opts)
}

// SendMessage is a helper function to easily call Bot.SendMessage in a chat.
func (c Chat) SendMessage(b *Bot, text string, opts *SendMessageOpts) (*Message, error) {
	return b.SendMessage(c.Id, text, opts)
}

// Unban is a helper function to easily call Bot.UnbanChatMember in a chat.
func (c Chat) Unban(b *Bot, userId int64, opts *UnbanChatMemberOpts) (bool, error) {
	return b.UnbanChatMember(c.Id, userId, opts)
}

// Promote is a helper function to easily call Bot.PromoteChatMember in a chat.
func (c Chat) Promote(b *Bot, userId int64, opts *PromoteChatMemberOpts) (bool, error) {
	return b.PromoteChatMember(c.Id, userId, opts)
}
