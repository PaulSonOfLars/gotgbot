package filters

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Message func(msg *gotgbot.Message) bool

func All(_ *gotgbot.Message) bool {
	return true
}

func Text(msg *gotgbot.Message) bool {
	return msg.Text != ""
}

func Caption(msg *gotgbot.Message) bool {
	return msg.Caption != ""
}

func Command(msg *gotgbot.Message) bool {
	return len(msg.Entities) > 0 && msg.Entities[0].Type == "bot_command" && msg.Entities[0].Offset == 0
}

func Reply(msg *gotgbot.Message) bool {
	return msg.ReplyToMessage != nil
}

func Audio(msg *gotgbot.Message) bool {
	return msg.Audio != nil
}

func Document(msg *gotgbot.Message) bool {
	return msg.Document != nil
}

func Photo(msg *gotgbot.Message) bool {
	return msg.Photo != nil
}

func Animation(msg *gotgbot.Message) bool {
	return msg.Animation != nil
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

func Location(msg *gotgbot.Message) bool {
	return msg.Location != nil
}

func Venue(msg *gotgbot.Message) bool {
	return msg.Venue != nil
}

func Forwarded(msg *gotgbot.Message) bool {
	return msg.ForwardDate != 0
}

func Game(msg *gotgbot.Message) bool {
	return msg.Game != nil
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

func Pin(msg *gotgbot.Message) bool {
	return msg.PinnedMessage != nil
}

func Dice(msg *gotgbot.Message) bool {
	return msg.Dice != nil
}

func DiceValue(msg *gotgbot.Message, val int64) bool {
	return msg.Dice != nil && msg.Dice.Value == val
}

func ViaBot(msg *gotgbot.Message) bool {
	return msg.ViaBot != nil
}

func Username(name string) Message {
	return func(m *gotgbot.Message) bool {
		return m.From.Username == name
	}
}

func Entity(entType string) Message {
	return func(m *gotgbot.Message) bool {
		for _, ent := range m.Entities {
			if ent.Type == entType {
				return true
			}
		}
		return false
	}
}

func CaptionEntity(entType string) Message {
	return func(m *gotgbot.Message) bool {
		for _, ent := range m.CaptionEntities {
			if ent.Type == entType {
				return true
			}
		}
		return false
	}
}

func MessageUserID(id int64) Message {
	return func(m *gotgbot.Message) bool {
		return m.From.Id == id
	}
}

func ChatUsername(name string) Message {
	return func(m *gotgbot.Message) bool {
		return m.Chat.Username != "" && m.Chat.Username == name
	}
}

func ChatID(id int64) Message {
	return func(m *gotgbot.Message) bool {
		return m.Chat.Id == id
	}
}

func NewChatMembers(msg *gotgbot.Message) bool {
	return msg.NewChatMembers != nil
}

func LeftChatMembers(msg *gotgbot.Message) bool {
	return msg.LeftChatMember != nil
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

func Poll(msg *gotgbot.Message) bool {
	return msg.Poll != nil
}

func Buttons(msg *gotgbot.Message) bool {
	return msg.ReplyMarkup != nil
}
