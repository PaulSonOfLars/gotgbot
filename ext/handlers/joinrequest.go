package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type ChatJoinRequest struct {
	Response Response
}

func NewChatJoinRequest(r Response) ChatJoinRequest {
	return ChatJoinRequest{
		Response: r,
	}
}

func (m ChatJoinRequest) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
  return u.ChatJoinRequest != nil
}

func (m ChatJoinRequest) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return m.Response(b, ctx)
}

func (m ChatJoinRequest) Name() string {
	return fmt.Sprintf("chatjoinrequest_%p", m.Response)
}
