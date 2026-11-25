package sharedwidgetquery

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.SharedWidgetQuery) bool {
	return true
}

func Prefix(prefix string) filters.SharedWidgetQuery {
	return func(swq *gotgbot.SharedWidgetQuery) bool {
		return swq != nil && strings.HasPrefix(swq.Data, prefix)
	}
}

func Suffix(suffix string) filters.SharedWidgetQuery {
	return func(swq *gotgbot.SharedWidgetQuery) bool {
		return swq != nil && strings.HasSuffix(swq.Data, suffix)
	}
}

func Equal(match string) filters.SharedWidgetQuery {
	return func(swq *gotgbot.SharedWidgetQuery) bool {
		return swq != nil && swq.Data == match
	}
}

func Contains(substr string) filters.SharedWidgetQuery {
	return func(swq *gotgbot.SharedWidgetQuery) bool {
		return swq != nil && strings.Contains(swq.Data, substr)
	}
}

func FromUserID(id int64) filters.SharedWidgetQuery {
	return func(swq *gotgbot.SharedWidgetQuery) bool {
		return swq != nil && swq.From.Id == id
	}
}

func WidgetID(widgetID string) filters.SharedWidgetQuery {
	return func(swq *gotgbot.SharedWidgetQuery) bool {
		return swq != nil && swq.WidgetId == widgetID
	}
}
