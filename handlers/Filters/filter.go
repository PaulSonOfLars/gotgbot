package Filters

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/ext"
)


func All(message *ext.Message) bool {
	return true
}

func Text(message *ext.Message) bool {
	return message.Text != ""
}

func Command(message *ext.Message) bool {
	return message.Text != "" && strings.HasPrefix(message.Text, "/")
}

func Reply(message *ext.Message) bool {
	return message.ReplyToMessage != nil
}

func Audio(message *ext.Message) bool {
	return message.Audio != nil
}

func Document(message *ext.Message) bool {
	return message.Document != nil
}

func Photo(message *ext.Message) bool {
	return message.Photo != nil
}

func Animation(message *ext.Message) bool {
	return message.Animation != nil
}

func Sticker(message *ext.Message) bool {
	return message.Sticker != nil
}

func Video(message *ext.Message) bool {
	return message.Video != nil
}

func VideoNote(message *ext.Message) bool {
	return message.VideoNote != nil
}

func Voice(message *ext.Message) bool {
	return message.Voice != nil
}

func Contact(message *ext.Message) bool {
	return message.Contact != nil
}

func Location(message *ext.Message) bool {
	return message.Location != nil
}

func Venue(message *ext.Message) bool {
	return message.Venue != nil
}

func Forwarded(message *ext.Message) bool {
	return message.ForwardDate != 0
}

func Game(message *ext.Message) bool {
	return message.Game != nil
}

func Private(message *ext.Message) bool {
	return message.Chat.Type == "private"
}

func Group(message *ext.Message) bool {
	return message.Chat.Type == "group" || message.Chat.Type == "supergroup"
}

func Username(name string) func(message *ext.Message) bool {
	return func(m *ext.Message) bool {
		return m.From.Username == name
	}
}

func Entity(entType string) func(message *ext.Message) bool {
	return func(m *ext.Message) bool {
		for _, ent := range m.Entities {
			if ent.Type == entType {
				return true
			}
		}
		return false
	}
}

func CaptionEntity(entType string) func(message *ext.Message) bool {
	return func(m *ext.Message) bool {
		for _, ent := range m.CaptionEntities {
			if ent.Type == entType {
				return true
			}
		}
		return false
	}
}

func UserID(id int) func(message *ext.Message) bool {
	return func(m *ext.Message) bool {
		return m.From.Id == id
	}
}

func Chatusername(name string) func(message *ext.Message) bool {
	return func(m *ext.Message) bool {
		return m.Chat != nil && m.Chat.Username != "" && m.Chat.Username == name
	}
}

func ChatID(id int) func(message *ext.Message) bool {
	return func(m *ext.Message) bool {
		return m.Chat != nil && m.Chat.Id == id
	}
}

func NewChatMembers() func(message *ext.Message) bool {
	return func(m *ext.Message) bool {
		return m.NewChatMembers != nil
	}
}

func LeftChatMembers() func(message *ext.Message) bool {
	return func(m *ext.Message) bool {
		return m.LeftChatMember != nil
	}
}

func Migrate() func(message *ext.Message) bool {
	return func(m *ext.Message) bool {
		return m.MigrateFromChatId != 0 || m.MigrateToChatId != 0
	}
}

func StartsWith(prefix string) func(message *ext.Message) bool {
	return func(m *ext.Message) bool {
		return (m.Text != "" && strings.HasPrefix(m.Text, prefix)) || (m.Caption != "" && strings.HasPrefix(m.Caption, prefix))
	}
}
