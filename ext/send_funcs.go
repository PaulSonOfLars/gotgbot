package ext

import (
	"io"
	"net/url"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/pkg/errors"
)

// TODO: Markdown and HTML - two different funcs?
func (b Bot) SendMessage(chatId int, text string) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	return newMsg.Send()
}

func (b Bot) ReplyText(chatId int, text string, replyToMessageId int) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	newMsg.ReplyToMessageId = replyToMessageId
	return newMsg.Send()
}

func (b Bot) ReplyHTML(chatId int, text string, replyToMessageId int) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	newMsg.ReplyToMessageId = replyToMessageId
	newMsg.ParseMode = parsemode.Html
	return newMsg.Send()
}

func (b Bot) ReplyMarkdown(chatId int, text string, replyToMessageId int) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	newMsg.ReplyToMessageId = replyToMessageId
	newMsg.ParseMode = parsemode.Markdown
	return newMsg.Send()
}

func (b Bot) ForwardMessage(chatId int, fromChatId int, messageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("from_chat_id", strconv.Itoa(fromChatId))
	v.Add("message_id", strconv.Itoa(messageId))

	r, err := Get(b, "forwardMessage", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to forwardMessage")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	return b.ParseMessage(r.Result), nil
}

func (b Bot) SendPhotoStr(chatId int, photoId string) (*Message, error) {
	return b.replyPhotoStr(chatId, photoId, "", 0)
}

func (b Bot) SendPhotoCaptionStr(chatId int, photoId string, caption string) (*Message, error) {
	return b.replyPhotoStr(chatId, photoId, caption, 0)
}

func (b Bot) ReplyPhotoStr(chatId int, photoId string, replyToMessageId int) (*Message, error) {
	return b.replyPhotoStr(chatId, photoId, "", replyToMessageId)
}

func (b Bot) ReplyPhotoCaptionStr(chatId int, photoId string, caption string, replyToMessageId int) (*Message, error) {
	return b.replyPhotoStr(chatId, photoId, caption, replyToMessageId)
}

func (b Bot) ReplyPhotoCaptionPath(chatId int, path string, caption string, replyToMessageId int) (*Message, error) {
	return b.replyPhotoPath(chatId, path, caption, replyToMessageId)
}

func (b Bot) ReplyPhotoCaptionReader(chatId int, reader io.Reader, caption string, replyToMessageId int) (*Message, error) {
	return b.replyPhotoReader(chatId, reader, caption, replyToMessageId)
}

func (b Bot) replyPhotoStr(chatId int, photo string, caption string, replyToMessageId int) (*Message, error) {
	photoMsg := b.NewSendablePhoto(chatId, caption)
	photoMsg.FileId = photo
	photoMsg.ReplyToMessageId = replyToMessageId
	return photoMsg.Send()
}

func (b Bot) replyPhotoPath(chatId int, path string, caption string, replyToMessageId int) (*Message, error) {
	photoMsg := b.NewSendablePhoto(chatId, caption)
	photoMsg.Path = path
	photoMsg.ReplyToMessageId = replyToMessageId
	return photoMsg.Send()
}

func (b Bot) replyPhotoReader(chatId int, reader io.Reader, caption string, replyToMessageId int) (*Message, error) {
	photoMsg := b.NewSendablePhoto(chatId, caption)
	photoMsg.Reader = reader
	photoMsg.ReplyToMessageId = replyToMessageId
	return photoMsg.Send()
}

func (b Bot) SendAudioStr(chatId int, audio string) (*Message, error) {
	return b.replyAudioStr(chatId, audio, "", 0)
}

func (b Bot) ReplyAudioStr(chatId int, audio string, replyToMessageId int) (*Message, error) {
	return b.replyAudioStr(chatId, audio, "", replyToMessageId)
}

func (b Bot) replyAudioStr(chatId int, audio string, caption string, replyToMessageId int) (*Message, error) {
	audioMsg := b.NewSendableAudio(chatId, caption)
	audioMsg.FileId = audio
	audioMsg.ReplyToMessageId = replyToMessageId
	return audioMsg.Send()
}

func (b Bot) replyAudioPath(chatId int, path string, caption string, replyToMessageId int) (*Message, error) {
	audioMsg := b.NewSendableAudio(chatId, caption)
	audioMsg.Path = path
	audioMsg.ReplyToMessageId = replyToMessageId
	return audioMsg.Send()
}

func (b Bot) replyAudioReader(chatId int, reader io.Reader, caption string, replyToMessageId int) (*Message, error) {
	audioMsg := b.NewSendableAudio(chatId, caption)
	audioMsg.Reader = reader
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

func (b Bot) ReplyDocumentPath(chatId int, path string, replyToMessageId int) (*Message, error) {
	return b.replyDocumentPath(chatId, path, "", replyToMessageId)
}

func (b Bot) ReplyDocumentCaptionPath(chatId int, path string, caption string, replyToMessageId int) (*Message, error) {
	return b.replyDocumentPath(chatId, path, caption, replyToMessageId)
}

func (b Bot) replyDocumentStr(chatId int, document string, caption string, replyToMessageId int) (*Message, error) {
	docMsg := b.NewSendableDocument(chatId, caption)
	docMsg.FileId = document
	docMsg.ReplyToMessageId = replyToMessageId
	return docMsg.Send()
}

func (b Bot) replyDocumentPath(chatId int, path string, caption string, replyToMessageId int) (*Message, error) {
	docMsg := b.NewSendableDocument(chatId, caption)
	docMsg.Path = path
	docMsg.ReplyToMessageId = replyToMessageId
	return docMsg.Send()
}

func (b Bot) replyDocumentReader(chatId int, reader io.Reader, caption string, replyToMessageId int) (*Message, error) {
	docMsg := b.NewSendableDocument(chatId, caption)
	docMsg.Reader = reader
	docMsg.ReplyToMessageId = replyToMessageId
	return docMsg.Send()
}

func (b Bot) SendVideoStr(chatId int, video string) (*Message, error) {
	return b.replyVideoStr(chatId, video, "", 0)
}

func (b Bot) ReplyVideoStr(chatId int, video string, replyToMessageId int) (*Message, error) {
	return b.replyVideoStr(chatId, video, "", replyToMessageId)
}

func (b Bot) replyVideoStr(chatId int, video string, caption string, replyToMessageId int) (*Message, error) {
	videoMsg := b.NewSendableVideo(chatId, caption)
	videoMsg.FileId = video
	videoMsg.ReplyToMessageId = replyToMessageId
	return videoMsg.Send()
}

func (b Bot) replyVideoPath(chatId int, path string, caption string, replyToMessageId int) (*Message, error) {
	videoMsg := b.NewSendableVideo(chatId, caption)
	videoMsg.Path = path
	videoMsg.ReplyToMessageId = replyToMessageId
	return videoMsg.Send()
}

func (b Bot) replyVideoReader(chatId int, reader io.Reader, caption string, replyToMessageId int) (*Message, error) {
	videoMsg := b.NewSendableVideo(chatId, caption)
	videoMsg.Reader = reader
	videoMsg.ReplyToMessageId = replyToMessageId
	return videoMsg.Send()
}

func (b Bot) SendVoiceStr(chatId int, voice string) (*Message, error) {
	return b.replyVoiceStr(chatId, voice, "", 0)
}

func (b Bot) ReplyVoiceStr(chatId int, voice string, replyToMessageId int) (*Message, error) {
	return b.replyVoiceStr(chatId, voice, "", replyToMessageId)
}

func (b Bot) replyVoiceStr(chatId int, voice string, caption string, replyToMessageId int) (*Message, error) {
	voiceMsg := b.NewSendableVoice(chatId, caption)
	voiceMsg.FileId = voice
	voiceMsg.ReplyToMessageId = replyToMessageId
	return voiceMsg.Send()
}

func (b Bot) replyVoicePath(chatId int, path string, caption string, replyToMessageId int) (*Message, error) {
	voiceMsg := b.NewSendableVoice(chatId, caption)
	voiceMsg.Path = path
	voiceMsg.ReplyToMessageId = replyToMessageId
	return voiceMsg.Send()
}

func (b Bot) replyVoiceReader(chatId int, reader io.Reader, caption string, replyToMessageId int) (*Message, error) {
	voiceMsg := b.NewSendableVoice(chatId, caption)
	voiceMsg.Reader = reader
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
	videoMsg.FileId = videoNote
	videoMsg.ReplyToMessageId = replyToMessageId
	return videoMsg.Send()
}

func (b Bot) replyVideoNotePath(chatId int, path string, replyToMessageId int) (*Message, error) {
	videoMsg := b.NewSendableVideoNote(chatId)
	videoMsg.Path = path
	videoMsg.ReplyToMessageId = replyToMessageId
	return videoMsg.Send()
}

func (b Bot) replyVideoNoteReader(chatId int, reader io.Reader, replyToMessageId int) (*Message, error) {
	videoMsg := b.NewSendableVideoNote(chatId)
	videoMsg.Reader = reader
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
	contactMsg.Action = action
	return contactMsg.Send()
}
