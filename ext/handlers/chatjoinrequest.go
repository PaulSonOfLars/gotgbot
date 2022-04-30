package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type ChatJoinRequest struct {
	Filter   filters.ChatJoinRequest
	Response Response
}

func NewChatJoinRequest(f filters.ChatJoinRequest, r Response) ChatJoinRequest {
	return ChatJoinRequest{
		Filter:   f,
		Response: r,
	}
}

func (r ChatJoinRequest) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
	if u.ChatJoinRequest == nil {
		return false
	}
	return r.Filter == nil || r.Filter(u.ChatJoinRequest)
}

func (r ChatJoinRequest) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return r.Response(b, ctx)
}

func (r ChatJoinRequest) Name() string {
	return fmt.Sprintf("chatjoinrequest_%p", r.Response)
}
