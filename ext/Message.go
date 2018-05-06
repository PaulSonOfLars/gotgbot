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

func (m Message) ReplyMessage(text string) (*Message, error) {
	return m.bot.ReplyMessage(m.Chat.Id, text, m.Message_id)
}

func (m Message) ReplyAudioStr(audio string) (*Message, error) {
	return m.bot.ReplyAudioStr(m.Chat.Id, audio, m.Message_id)
}

func (m Message) ReplyDocumentStr(document string) (*Message, error) {
	return m.bot.ReplyDocumentStr(m.Chat.Id, document, m.Message_id)
}

func (m Message) ReplyLocation(latitude float64, longitude float64) (*Message, error) {
	return m.bot.ReplyLocation(m.Chat.Id, latitude, longitude, m.Message_id)
}

func (m Message) ReplyPhotoStr(photo string) (*Message, error) {
	return m.bot.ReplyPhotoStr(m.Chat.Id, photo, m.Message_id)
}

func (m Message) ReplyStickerStr(sticker string) (*Message, error) {
	return m.bot.ReplyStickerStr(m.Chat.Id, sticker, m.Message_id)
}

func (m Message) ReplyVenue(latitude float64, longitude float64, title string, address string) (*Message, error) {
	return m.bot.ReplyVenue(m.Chat.Id, latitude, longitude, title, address, m.Message_id)
}

func (m Message) ReplyVideoStr(video string) (*Message, error) {
	return m.bot.ReplyVideoStr(m.Chat.Id, video, m.Message_id)
}

func (m Message) ReplyVideoNoteStr(videoNote string) (*Message, error) {
	return m.bot.ReplyVideoNoteStr(m.Chat.Id, videoNote, m.Message_id)
}

func (m Message) ReplyVoiceStr(voice string) (*Message, error) {
	return m.bot.ReplyVoiceStr(m.Chat.Id, voice, m.Message_id)
}

func (m Message) Delete() (bool, error) {
	return m.bot.DeleteMessage(m.Chat.Id, m.Message_id)
}

func (m Message) Forward(chatId int) (*Message, error) {
	return m.bot.ForwardMessage(chatId, m.Chat.Id, m.Message_id)
}