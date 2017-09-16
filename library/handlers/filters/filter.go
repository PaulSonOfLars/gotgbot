package filters

import (
	"bot/library"
	"strings"
)


func All(message *library.Message) bool {
	return true
}

func Text(message *library.Message) bool {
	return message.Text != ""
}


func Command(message *library.Message) bool {
	return message.Text != "" && strings.HasPrefix(message.Text, "/")
}


func Reply(message *library.Message) bool {
	return message.Reply_to_message != nil
}


func Audio(message *library.Message) bool {
	return message.Audio != nil
}

func Document(message *library.Message) bool {
	return message.Document != nil
}

func Photo(message *library.Message) bool {
	return message.Photo != nil
}

func Sticker(message *library.Message) bool {
	return message.Sticker != nil
}

func Video(message *library.Message) bool {
	return message.Video != nil
}

func Contact(message *library.Message) bool {
	return message.Contact != nil
}

func Location(message *library.Message) bool {
	return message.Location != nil
}

func Venue(message *library.Message) bool {
	return message.Venue != nil
}

func Forwarded(message *library.Message) bool {
	return message.Forward_date != 0
}

func Game(message *library.Message) bool {
	return message.Game != nil
}

//type Private struct {}
//
//func (cf Private) Eval(message library.Message) bool {
//	return message.Chat.Type == PRIVATE
//}

//type Group struct {}
//
//func (cf Group) Eval(message library.Message) bool {
//	return message.Chat.Type == GROUP or SUPERGROUP
//}

func Username(name string) func(message *library.Message) bool {
	return func(m *library.Message) bool {
		return m.From.Username == name
	}
}

func UserID(id int) func(message *library.Message) bool {
	return func(m *library.Message) bool {
		return m.From.Id == id
	}
}

func Chatusername(name string) func(message *library.Message) bool {
	return func(m *library.Message) bool {
		return m.Chat != nil && m.Chat.Username != "" && m.Chat.Username == name
	}
}

func ChatID(id int) func(message *library.Message) bool {
	return func(m *library.Message) bool {
		return m.Chat != nil && m.Chat.Id == id
	}
}
