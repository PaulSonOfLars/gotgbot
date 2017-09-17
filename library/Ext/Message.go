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

func (m Message) reply_text(text string) Message {
	return m.bot.SendMessage(text, m.Chat.Id)
}