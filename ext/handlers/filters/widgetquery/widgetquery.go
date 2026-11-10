package widgetquery

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.WidgetQuery) bool {
	return true
}

func Prefix(prefix string) filters.WidgetQuery {
	return func(wq *gotgbot.WidgetQuery) bool {
		return strings.HasPrefix(wq.Data, prefix)
	}
}

func Suffix(suffix string) filters.WidgetQuery {
	return func(wq *gotgbot.WidgetQuery) bool {
		return strings.HasSuffix(wq.Data, suffix)
	}
}

func Equal(match string) filters.WidgetQuery {
	return func(wq *gotgbot.WidgetQuery) bool {
		return wq.Data == match
	}
}

func Contains(substr string) filters.WidgetQuery {
	return func(wq *gotgbot.WidgetQuery) bool {
		return strings.Contains(wq.Data, substr)
	}
}

func FromUserID(id int64) filters.WidgetQuery {
	return func(wq *gotgbot.WidgetQuery) bool {
		return wq.From.Id == id
	}
}

func ChatInstance(instance string) filters.WidgetQuery {
	return func(wq *gotgbot.WidgetQuery) bool {
		return wq.ChatInstance == instance
	}
}

func ChatType(chatType string) filters.WidgetQuery {
	return func(wq *gotgbot.WidgetQuery) bool {
		return wq.ChatType == chatType
	}
}

func WidgetID(widgetID string) filters.WidgetQuery {
	return func(wq *gotgbot.WidgetQuery) bool {
		return wq.WidgetId == widgetID
	}
}
