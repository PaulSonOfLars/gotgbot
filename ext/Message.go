package ext

import (
	"gotgbot/types"
	"encoding/json"
)

type Message struct {
	types.Message // Contain all message structs
	bot Bot
}

func (b Bot) Message(chatId int, text string) Message {
	return Message{bot: b}
}

func (b Bot) ParseMessage(message json.RawMessage) *Message {
	var mess Message
	json.Unmarshal(message, &mess)
	mess.bot = b
	return &mess
}

func (b Bot) NewMessage(message *types.Message) *Message {
	return &Message{*message, b}
}

func (b Bot) NewChat(chat *types.Chat) *Chat{
	return &Chat{*chat, b}
}

func (b Bot) NewUser(user *types.User) *User{
	return &User{*user, b}
}

func (m Message) ReplyMessage(text string) (*Message, error) {
	return m.bot.ReplyMessage(m.Chat.Id, text, m.MessageId)
}

func (m Message) ReplyAudioStr(audio string) (*Message, error) {
	return m.bot.ReplyAudioStr(m.Chat.Id, audio, m.MessageId)
}

func (m Message) ReplyDocumentStr(document string) (*Message, error) {
	return m.bot.ReplyDocumentStr(m.Chat.Id, document, m.MessageId)
}

func (m Message) ReplyLocation(latitude float64, longitude float64) (*Message, error) {
	return m.bot.ReplyLocation(m.Chat.Id, latitude, longitude, m.MessageId)
}

func (m Message) ReplyPhotoStr(photo string) (*Message, error) {
	return m.bot.ReplyPhotoStr(m.Chat.Id, photo, m.MessageId)
}

func (m Message) ReplyStickerStr(sticker string) (*Message, error) {
	return m.bot.ReplyStickerStr(m.Chat.Id, sticker, m.MessageId)
}

func (m Message) ReplyVenue(latitude float64, longitude float64, title string, address string) (*Message, error) {
	return m.bot.ReplyVenue(m.Chat.Id, latitude, longitude, title, address, m.MessageId)
}

func (m Message) ReplyVideoStr(video string) (*Message, error) {
	return m.bot.ReplyVideoStr(m.Chat.Id, video, m.MessageId)
}

func (m Message) ReplyVideoNoteStr(videoNote string) (*Message, error) {
	return m.bot.ReplyVideoNoteStr(m.Chat.Id, videoNote, m.MessageId)
}

func (m Message) ReplyVoiceStr(voice string) (*Message, error) {
	return m.bot.ReplyVoiceStr(m.Chat.Id, voice, m.MessageId)
}

func (m Message) Delete() (bool, error) {
	return m.bot.DeleteMessage(m.Chat.Id, m.MessageId)
}

func (m Message) Forward(chatId int) (*Message, error) {
	return m.bot.ForwardMessage(chatId, m.Chat.Id, m.MessageId)
}