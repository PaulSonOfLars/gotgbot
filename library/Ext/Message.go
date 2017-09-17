package Ext

import (
	"bot/library/Types"
	"encoding/json"
)

type Message struct {
	Types.Message // Contain all message structs
	bot Bot
}

func (b Bot) ParseMessage(message json.RawMessage) Message {
	var mess Message
	json.Unmarshal(message, &mess)
	mess.bot = b
	return mess
}

func (b Bot) NewMessage(message *Types.Message) *Message {
	return &Message{*message, b}

}

func (m Message) Reply_text(text string) Message {
	return m.bot.SendMessage(text, m.Chat.Id)
}

func (m Message) DeleteMessage() bool {
	return m.bot.DeleteMessage(m.Chat.Id, m.Message_id)
}