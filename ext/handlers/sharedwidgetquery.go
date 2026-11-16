package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

// SharedWidgetQuery represents a handler for Telegram SharedWidgetQuery updates.
type SharedWidgetQuery struct {
	Filter   filters.SharedWidgetQuery
	Response Response
}

// NewSharedWidgetQuery creates a new SharedWidgetQuery handler with the given filter and response.
func NewSharedWidgetQuery(filter filters.SharedWidgetQuery, r Response) SharedWidgetQuery {
	return SharedWidgetQuery{
		Filter:   filter,
		Response: r,
	}
}

func (w SharedWidgetQuery) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	if w.Response == nil {
		return nil
	}
	return w.Response(b, ctx)
}

func (w SharedWidgetQuery) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx == nil || ctx.SharedWidgetQuery == nil {
		return false
	}
	return w.Filter == nil || w.Filter(ctx.SharedWidgetQuery)
}

func (w SharedWidgetQuery) Name() string {
	return fmt.Sprintf("shared_widget_query_%p", w.Response)
}
