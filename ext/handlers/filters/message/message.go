package message

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.Message) bool {
	return true
}

func FromUserID(id int64) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.From.Id == id
	}
}

func FromUsername(name string) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.From.Username == name
	}
}

func ChatUsername(name string) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.Chat.Username != "" && m.Chat.Username == name
	}
}

func ChatID(id int64) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.Chat.Id == id
	}
}

func ForwardFromUserID(id int64) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.ForwardFrom != nil && m.ForwardFrom.Id == id
	}
}

func ForwardFromChatID(id int64) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.ForwardFromChat != nil && m.ForwardFromChat.Id == id
	}
}

func Regex(p string) (filters.Message, error) {
	r, err := regexp.Compile(p)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}
	return func(m *gotgbot.Message) bool {
		if m.Text != "" {
			return r.MatchString(m.Text)
		} else if m.Caption != "" {
			return r.MatchString(m.Caption)
		}
		return false
	}, nil
}

func Reply(msg *gotgbot.Message) bool {
	return msg.ReplyToMessage != nil
}

func ChatType(t string) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.Chat.Type == t
	}
}

func Private(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "private"
}

func Group(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "group"
}

func Supergroup(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "supergroup"
}

func Channel(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "channel"
}

func Forwarded(msg *gotgbot.Message) bool {
	return msg.ForwardDate != 0
}

func Text(msg *gotgbot.Message) bool {
	return msg.Text != ""
}

func HasPrefix(prefix string) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return strings.HasPrefix(msg.Text, prefix) || strings.HasSuffix(msg.Caption, prefix)

	}
}

func HasSuffix(suffix string) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return strings.HasSuffix(msg.Text, suffix) || strings.HasSuffix(msg.Caption, suffix)
	}
}

func Contains(contains string) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return strings.Contains(msg.Text, contains) || strings.Contains(msg.Caption, contains)
	}
}

func Equal(eq string) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return msg.Text == eq || msg.Caption == eq
	}
}

func Caption(msg *gotgbot.Message) bool {
	return msg.Caption != ""
}

func Command(msg *gotgbot.Message) bool {
	if msg.Text != "" {
		return len(msg.Entities) > 0 && msg.Entities[0].Type == "bot_command" && msg.Entities[0].Offset == 0
	} else if msg.Caption != "" {
		return len(msg.CaptionEntities) > 0 && msg.CaptionEntities[0].Type == "bot_command" && msg.CaptionEntities[0].Offset == 0
	}
	return false
}

func Animation(msg *gotgbot.Message) bool {
	return msg.Animation != nil
}

func Audio(msg *gotgbot.Message) bool {
	return msg.Audio != nil
}

func Document(msg *gotgbot.Message) bool {
	return msg.Document != nil
}

func Photo(msg *gotgbot.Message) bool {
	return len(msg.Photo) > 0
}

func Sticker(msg *gotgbot.Message) bool {
	return msg.Sticker != nil
}

func Video(msg *gotgbot.Message) bool {
	return msg.Video != nil
}

func VideoNote(msg *gotgbot.Message) bool {
	return msg.VideoNote != nil
}

func Voice(msg *gotgbot.Message) bool {
	return msg.Voice != nil
}

func Contact(msg *gotgbot.Message) bool {
	return msg.Contact != nil
}

func Dice(msg *gotgbot.Message) bool {
	return msg.Dice != nil
}

func DiceValue(msg *gotgbot.Message, val int64) bool {
	return msg.Dice != nil && msg.Dice.Value == val
}

func Game(msg *gotgbot.Message) bool {
	return msg.Game != nil
}

func Poll(msg *gotgbot.Message) bool {
	return msg.Poll != nil
}

func Venue(msg *gotgbot.Message) bool {
	return msg.Venue != nil
}

func Location(msg *gotgbot.Message) bool {
	return msg.Location != nil
}

func NewChatMembers(msg *gotgbot.Message) bool {
	return msg.NewChatMembers != nil
}

func LeftChatMember(msg *gotgbot.Message) bool {
	return msg.LeftChatMember != nil
}

func PinnedMessage(msg *gotgbot.Message) bool {
	return msg.PinnedMessage != nil
}

func ViaBot(msg *gotgbot.Message) bool {
	return msg.ViaBot != nil
}

func Entities(m *gotgbot.Message) bool {
	return len(m.Entities) > 0
}

func Entity(entType string) filters.Message {
	return func(m *gotgbot.Message) bool {
		for _, ent := range m.Entities {
			if ent.Type == entType {
				return true
			}
		}
		return false
	}
}

func CaptionEntities(m *gotgbot.Message) bool {
	return len(m.CaptionEntities) > 0
}

func CaptionEntity(entType string) filters.Message {
	return func(m *gotgbot.Message) bool {
		for _, ent := range m.CaptionEntities {
			if ent.Type == entType {
				return true
			}
		}
		return false
	}
}

func Migrate(msg *gotgbot.Message) bool {
	return msg.MigrateFromChatId != 0 || msg.MigrateToChatId != 0
}

func MigrateFrom(msg *gotgbot.Message) bool {
	return msg.MigrateFromChatId != 0
}

func MigrateTo(msg *gotgbot.Message) bool {
	return msg.MigrateToChatId != 0
}

func ReplyMarkup(msg *gotgbot.Message) bool {
	return msg.ReplyMarkup != nil
}

func MediaGroup(msg *gotgbot.Message) bool {
	return msg.MediaGroupId != ""
}

func IsAutomaticForward(msg *gotgbot.Message) bool {
	return msg.IsAutomaticForward
}
