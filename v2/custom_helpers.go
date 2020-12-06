package gotgbot

func (m Message) Reply(b *Bot, text string, opts SendMessageOpts) (*Message, error) {
	if opts.ReplyToMessageId == 0 {
		opts.ReplyToMessageId = m.MessageId
	}
	return b.SendMessage(m.Chat.Id, text, opts)
}

func (c Chat) SendMessage(b *Bot, text string, opts SendMessageOpts) (*Message, error) {
	return b.SendMessage(c.Id, text, opts)
}

func (c Chat) Unban(b *Bot, userId int64, opts UnbanChatMemberOpts) (bool, error) {
	return b.UnbanChatMember(c.Id, userId, opts)
}

func (c Chat) Promote(b *Bot, userId int64, opts PromoteChatMemberOpts) (bool, error) {
	return b.PromoteChatMember(c.Id, userId, opts)
}
