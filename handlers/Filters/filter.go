package Filters

import (
	"strings"
	"gotgbot/types"
)


func All(message *types.Message) bool {
	return true
}

func Text(message *types.Message) bool {
	return message.Text != ""
}


func Command(message *types.Message) bool {
	return message.Text != "" && strings.HasPrefix(message.Text, "/")
}


func Reply(message *types.Message) bool {
	return message.Reply_to_message != nil
}


func Audio(message *types.Message) bool {
	return message.Audio != nil
}

func Document(message *types.Message) bool {
	return message.Document != nil
}

func Photo(message *types.Message) bool {
	return message.Photo != nil
}

func Sticker(message *types.Message) bool {
	return message.Sticker != nil
}

func Video(message *types.Message) bool {
	return message.Video != nil
}

func Contact(message *types.Message) bool {
	return message.Contact != nil
}

func Location(message *types.Message) bool {
	return message.Location != nil
}

func Venue(message *types.Message) bool {
	return message.Venue != nil
}

func Forwarded(message *types.Message) bool {
	return message.Forward_date != 0
}

func Game(message *types.Message) bool {
	return message.Game != nil
}

//func Private(message *ext.Message) bool {
//	return message.Chat.Type == PRIVATE
//}
//
//func Group(message *ext.Message) bool {
//	return message.Chat.Type == GROUP or SUPERGROUP
//}

func Username(name string) func(message *types.Message) bool {
	return func(m *types.Message) bool {
		return m.From.Username == name
	}
}

func UserID(id int) func(message *types.Message) bool {
	return func(m *types.Message) bool {
		return m.From.Id == id
	}
}

func Chatusername(name string) func(message *types.Message) bool {
	return func(m *types.Message) bool {
		return m.Chat != nil && m.Chat.Username != "" && m.Chat.Username == name
	}
}

func ChatID(id int) func(message *types.Message) bool {
	return func(m *types.Message) bool {
		return m.Chat != nil && m.Chat.Id == id
	}
}
