package ext

import (
	"strconv"
	"net/url"
	"github.com/pkg/errors"
)

// TODO: Markdown and HTML - two different funcs?
func (b Bot) SendMessage(chatId int, text string) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	return newMsg.Send()
}

func (b Bot) ReplyMessage(chatId int, text string, replyToMessageId int) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	newMsg.ReplyToMessageId = replyToMessageId
	return newMsg.Send()
}

func (b Bot) ForwardMessage(chatId int, fromChatId int, messageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("from_chat_id", strconv.Itoa(fromChatId))
	v.Add("message_id", strconv.Itoa(messageId))

	r := Get(b, "forwardMessage", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	return b.ParseMessage(r.Result), nil
}

// TODO: create InputFile version of all the Str's
func (b Bot) SendPhotoStr(chatId int, photo string) (*Message, error) {
	return b.replyPhotoStr(chatId, photo, "", 0)
}

func (b Bot) SendPhotoCaptionStr(chatId int, photo string, caption string) (*Message, error) {
	return b.replyPhotoStr(chatId, photo, caption, 0)
}

func (b Bot) ReplyPhotoStr(chatId int, photo string, replyToMessageId int) (*Message, error) {
	return b.replyPhotoStr(chatId, photo, "", replyToMessageId)
}

func (b Bot) ReplyPhotoCaptionStr(chatId int, photo string, caption string, replyToMessageId int) (*Message, error) {
	return b.replyPhotoStr(chatId, photo, caption, replyToMessageId)
}

func (b Bot) replyPhotoStr(chatId int, photo string, caption string, replyToMessageId int) (*Message, error) {
	photoMsg := b.NewSendablePhoto(chatId, caption)
	photoMsg.PhotoString = photo
	photoMsg.ReplyToMessageId = replyToMessageId
	return photoMsg.Send()
}

func (b Bot) SendAudioStr(chatId int, audio string) (*Message, error) {
	return b.replyAudioStr(chatId, audio, 0)
}

func (b Bot) ReplyAudioStr(chatId int, audio string, replyToMessageId int) (*Message, error) {
	return b.replyAudioStr(chatId, audio, replyToMessageId)
}

func (b Bot) replyAudioStr(chatId int, audio string, replyToMessageId int) (*Message, error) {
	audioMsg := b.NewSendableAudio(chatId, "")
	audioMsg.AudioString = audio
	audioMsg.ReplyToMessageId = replyToMessageId
	return audioMsg.Send()
}

func (b Bot) SendDocumentStr(chatId int, document string) (*Message, error) {
	return b.replyDocumentStr(chatId, document, "", 0)
}

func (b Bot) SendDocumentCaptionStr(chatId int, document string, caption string) (*Message, error) {
	return b.replyDocumentStr(chatId, document, caption, 0)
}

func (b Bot) ReplyDocumentStr(chatId int, document string, replyToMessageId int) (*Message, error) {
	return b.replyDocumentStr(chatId, document, "", replyToMessageId)
}

func (b Bot) ReplyDocumentCaptionStr(chatId int, document string, caption string, replyToMessageId int) (*Message, error) {
	return b.replyDocumentStr(chatId, document, caption, replyToMessageId)
}

func (b Bot) replyDocumentStr(chatId int, document string, caption string, replyToMessageId int) (*Message, error) {
	docMsg := b.NewSendableDocument(chatId, caption)
	docMsg.DocString = document
	docMsg.ReplyToMessageId = replyToMessageId
	return docMsg.Send()
}

func (b Bot) SendVideoStr(chatId int, video string) (*Message, error) {
	return b.replyVideoStr(chatId, video, 0)
}

func (b Bot) ReplyVideoStr(chatId int, video string, replyToMessageId int) (*Message, error) {
	return b.replyVideoStr(chatId, video, replyToMessageId)
}

func (b Bot) replyVideoStr(chatId int, video string, replyToMessageId int) (*Message, error) {
	videoMsg := b.NewSendableVideo(chatId, "")
	videoMsg.VideoString = video
	videoMsg.ReplyToMessageId = replyToMessageId
	return videoMsg.Send()
}

func (b Bot) SendVoiceStr(chatId int, voice string) (*Message, error) {
	return b.replyVoiceStr(chatId, voice, 0)
}

func (b Bot) ReplyVoiceStr(chatId int, voice string, replyToMessageId int) (*Message, error) {
	return b.replyVoiceStr(chatId, voice, replyToMessageId)
}

func (b Bot) replyVoiceStr(chatId int, voice string, replyToMessageId int) (*Message, error) {
	voiceMsg := b.NewSendableVoice(chatId, "")
	voiceMsg.VoiceString = voice
	voiceMsg.ReplyToMessageId = replyToMessageId
	return voiceMsg.Send()
}

func (b Bot) SendVideoNoteStr(chatId int, videoNote string) (*Message, error) {
	return b.replyVideoNoteStr(chatId, videoNote, 0)
}

func (b Bot) ReplyVideoNoteStr(chatId int, videoNote string, replyToMessageId int) (*Message, error) {
	return b.replyVideoNoteStr(chatId, videoNote, replyToMessageId)
}

func (b Bot) replyVideoNoteStr(chatId int, videoNote string, replyToMessageId int) (*Message, error) {
	videoMsg := b.NewSendableVideoNote(chatId)
	videoMsg.VideoNoteString = videoNote
	videoMsg.ReplyToMessageId = replyToMessageId
	return videoMsg.Send()
}

func (b Bot) SendLocation(chatId int, latitude float64, longitude float64) (*Message, error) {
	return b.replyLocation(chatId, latitude, longitude, 0)
}

func (b Bot) ReplyLocation(chatId int, latitude float64, longitude float64, replyToMessageId int) (*Message, error) {
	return b.replyLocation(chatId, latitude, longitude, replyToMessageId)
}

func (b Bot) replyLocation(chatId int, latitude float64, longitude float64, replyToMessageId int) (*Message, error) {
	locationMsg := b.NewSendableLocation(chatId)
	locationMsg.Latitude = latitude
	locationMsg.Longitude = longitude
	locationMsg.ReplyToMessageId = replyToMessageId
	return locationMsg.Send()
}

func (b Bot) SendVenue(chatId int, latitude float64, longitude float64, title string, address string) (*Message, error) {
	return b.replyVenue(chatId, latitude, longitude, title, address, 0)
}

func (b Bot) ReplyVenue(chatId int, latitude float64, longitude float64, title string, address string, replyToMessageId int) (*Message, error) {
	return b.replyVenue(chatId, latitude, longitude, title, address, replyToMessageId)
}

func (b Bot) replyVenue(chatId int, latitude float64, longitude float64, title string, address string, replyToMessageId int) (*Message, error) {
	venueMsg := b.NewSendableVenue(chatId)
	venueMsg.Latitude = latitude
	venueMsg.Longitude = longitude
	venueMsg.Title = title
	venueMsg.Address = address
	venueMsg.ReplyToMessageId = replyToMessageId
	return venueMsg.Send()
}

func (b Bot) SendContact(chatId int, phoneNumber string, firstName string) (*Message, error) {
	return b.replyContact(chatId, phoneNumber, firstName, 0)
}

func (b Bot) ReplyContact(chatId int, phoneNumber string, firstName string, replyToMessageId int) (*Message, error) {
	return b.replyContact(chatId, phoneNumber, firstName, replyToMessageId)
}

func (b Bot) replyContact(chatId int, phoneNumber string, firstName string, replyToMessageId int) (*Message, error) {
	contactMsg := b.NewSendableContact(chatId)
	contactMsg.PhoneNumber = phoneNumber
	contactMsg.FirstName = firstName
	contactMsg.ReplyToMessageId = replyToMessageId
	return contactMsg.Send()
}

func (b Bot) SendChatAction(chatId int, action string) (bool, error) {
	contactMsg := b.NewSendableChatAction(chatId)
	contactMsg.action = action
	return contactMsg.Send()
}
