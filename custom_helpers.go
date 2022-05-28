package gotgbot

import (
	"fmt"
	"strconv"
	"strings"
)

// GetLink is a helper method to easily get the message link (It will return an empty string in case of private or group chat type).
func (m Message) GetLink() string {
	if m.Chat.Type == "private" || m.Chat.Type == "group" {
		return ""
	}
	if m.Chat.Username != "" {
		return fmt.Sprintf("https://t.me/%s/%d", m.Chat.Username, m.MessageId)
	}
	// Message links use raw chatIds without the -100 prefix; this trims that prefix.
	rawChatId := strings.TrimPrefix(strconv.FormatInt(m.Chat.Id, 10), "-100")
	return fmt.Sprintf("https://t.me/c/%s/%d", rawChatId, m.MessageId)
}

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

// GetURL gets the URL the file can be downloaded from.
func (f File) GetURL(b *Bot) string {
	return fmt.Sprintf("%s/file/bot%s/%s", b.GetAPIURL(), b.GetToken(), f.FilePath)
}
