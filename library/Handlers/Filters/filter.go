package Filters

import (
	"strings"
	"bot/library/Ext"
)


func All(message *Ext.Message) bool {
	return true
}

func Text(message *Ext.Message) bool {
	return message.Text != ""
}


func Command(message *Ext.Message) bool {
	return message.Text != "" && strings.HasPrefix(message.Text, "/")
}


func Reply(message *Ext.Message) bool {
	return message.Reply_to_message != nil
}


func Audio(message *Ext.Message) bool {
	return message.Audio != nil
}

func Document(message *Ext.Message) bool {
	return message.Document != nil
}

func Photo(message *Ext.Message) bool {
	return message.Photo != nil
}

func Sticker(message *Ext.Message) bool {
	return message.Sticker != nil
}

func Video(message *Ext.Message) bool {
	return message.Video != nil
}

func Contact(message *Ext.Message) bool {
	return message.Contact != nil
}

func Location(message *Ext.Message) bool {
	return message.Location != nil
}

func Venue(message *Ext.Message) bool {
	return message.Venue != nil
}

func Forwarded(message *Ext.Message) bool {
	return message.Forward_date != 0
}

func Game(message *Ext.Message) bool {
	return message.Game != nil
}

//func Private(message *Ext.Message) bool {
//	return message.Chat.Type == PRIVATE
//}
//
//func Group(message *Ext.Message) bool {
//	return message.Chat.Type == GROUP or SUPERGROUP
//}

func Username(name string) func(message *Ext.Message) bool {
	return func(m *Ext.Message) bool {
		return m.From.Username == name
	}
}

func UserID(id int) func(message *Ext.Message) bool {
	return func(m *Ext.Message) bool {
		return m.From.Id == id
	}
}

func Chatusername(name string) func(message *Ext.Message) bool {
	return func(m *Ext.Message) bool {
		return m.Chat != nil && m.Chat.Username != "" && m.Chat.Username == name
	}
}

func ChatID(id int) func(message *Ext.Message) bool {
	return func(m *Ext.Message) bool {
		return m.Chat != nil && m.Chat.Id == id
	}
}
