package filters

import (
	"strings"
	"bot/library/Types"
)


func All(message *Types.Message) bool {
	return true
}

func Text(message *Types.Message) bool {
	return message.Text != ""
}


func Command(message *Types.Message) bool {
	return message.Text != "" && strings.HasPrefix(message.Text, "/")
}


func Reply(message *Types.Message) bool {
	return message.Reply_to_message != nil
}


func Audio(message *Types.Message) bool {
	return message.Audio != nil
}

func Document(message *Types.Message) bool {
	return message.Document != nil
}

func Photo(message *Types.Message) bool {
	return message.Photo != nil
}

func Sticker(message *Types.Message) bool {
	return message.Sticker != nil
}

func Video(message *Types.Message) bool {
	return message.Video != nil
}

func Contact(message *Types.Message) bool {
	return message.Contact != nil
}

func Location(message *Types.Message) bool {
	return message.Location != nil
}

func Venue(message *Types.Message) bool {
	return message.Venue != nil
}

func Forwarded(message *Types.Message) bool {
	return message.Forward_date != 0
}

func Game(message *Types.Message) bool {
	return message.Game != nil
}

//func Private(message *Types.Message) bool {
//	return message.Chat.Type == PRIVATE
//}
//
//func Group(message *Types.Message) bool {
//	return message.Chat.Type == GROUP or SUPERGROUP
//}

func Username(name string) func(message *Types.Message) bool {
	return func(m *Types.Message) bool {
		return m.From.Username == name
	}
}

func UserID(id int) func(message *Types.Message) bool {
	return func(m *Types.Message) bool {
		return m.From.Id == id
	}
}

func Chatusername(name string) func(message *Types.Message) bool {
	return func(m *Types.Message) bool {
		return m.Chat != nil && m.Chat.Username != "" && m.Chat.Username == name
	}
}

func ChatID(id int) func(message *Types.Message) bool {
	return func(m *Types.Message) bool {
		return m.Chat != nil && m.Chat.Id == id
	}
}
