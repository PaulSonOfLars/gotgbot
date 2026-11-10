package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

// WidgetQuery represents a handler for Telegram WidgetQuery updates.
type WidgetQuery struct {
	Filter   filters.WidgetQuery
	Response Response
}

// NewWidgetQuery creates a new WidgetQuery handler with the given filter and response.
func NewWidgetQuery(filter filters.WidgetQuery, r Response) WidgetQuery {
	return WidgetQuery{
		Filter:   filter,
		Response: r,
	}
}

func (w WidgetQuery) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return w.Response(b, ctx)
}

func (w WidgetQuery) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.WidgetQuery == nil {
		return false
	}
	return w.Filter == nil || w.Filter(ctx.WidgetQuery)
}

func (w WidgetQuery) Name() string {
	return fmt.Sprintf("widget_query_%p", w.Response)
}
