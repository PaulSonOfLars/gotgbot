package ext

import (
	"strconv"
	"log"
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
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("photo", photo)
	v.Add("caption", caption)
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r := Get(b, "sendPhoto", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return b.ParseMessage(r.Result), nil
}

func (b Bot) SendAudioStr(chatId int, audio string) (*Message, error) {
	return b.replyAudioStr(chatId, audio, 0)
}

func (b Bot) ReplyAudioStr(chatId int, audio string, replyToMessageId int) (*Message, error) {
	return b.replyAudioStr(chatId, audio, replyToMessageId)
}

func (b Bot) replyAudioStr(chatId int, audio string, replyToMessageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("audio", audio)
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r := Get(b, "sendAudio", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result), nil
}

func (b Bot) SendDocumentStr(chatId int, document string) (*Message, error) {
	return b.replyDocumentStr(chatId, document,"", 0)
}

func (b Bot) SendDocumentCaptionStr(chatId int, photo string, caption string) (*Message, error) {
	return b.replyDocumentStr(chatId, photo, caption, 0)
}

func (b Bot) ReplyDocumentStr(chatId int, document string, replyToMessageId int) (*Message, error) {
	return b.replyDocumentStr(chatId, document, "", replyToMessageId)
}

func (b Bot) ReplyDocumentCaptionStr(chatId int, photo string, caption string, replyToMessageId int) (*Message, error) {
	return b.replyDocumentStr(chatId, photo, caption, replyToMessageId)
}

func (b Bot) replyDocumentStr(chatId int, document string, caption string, replyToMessageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("document", document)
	v.Add("caption", caption)
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r := Get(b, "sendDocument", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return b.ParseMessage(r.Result), nil
}

func (b Bot) SendVideoStr(chatId int, video string) (*Message, error) {
	return b.replyVideoStr(chatId, video, 0)
}

func (b Bot) ReplyVideoStr(chatId int, video string, replyToMessageId int) (*Message, error) {
	return b.replyVideoStr(chatId, video, replyToMessageId)
}

func (b Bot) replyVideoStr(chatId int, video string, replyToMessageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("video", video)
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r := Get(b, "sendVideo", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return b.ParseMessage(r.Result), nil
}

func (b Bot) SendVoiceStr(chatId int, voice string) (*Message, error) {
	return b.replyVoiceStr(chatId, voice, 0)
}

func (b Bot) ReplyVoiceStr(chatId int, voice string, replyToMessageId int) (*Message, error) {
	return b.replyVoiceStr(chatId, voice, replyToMessageId)
}

func (b Bot) replyVoiceStr(chatId int, voice string, replyToMessageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("voice", voice)
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r := Get(b, "sendVoice", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return b.ParseMessage(r.Result), nil
}

func (b Bot) SendVideoNoteStr(chatId int, videoNote string) (*Message, error) {
	return b.replyVideoNoteStr(chatId, videoNote, 0)
}

func (b Bot) ReplyVideoNoteStr(chatId int, videoNote string, replyToMessageId int) (*Message, error) {
	return b.replyVideoNoteStr(chatId, videoNote, replyToMessageId)
}

func (b Bot) replyVideoNoteStr(chatId int, videoNote string, replyToMessageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("video_note", videoNote)
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r := Get(b, "sendVideoNote", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return b.ParseMessage(r.Result), nil
}

func (b Bot) SendLocation(chatId int, latitude float64, longitude float64) (*Message, error) {
	return b.replyLocation(chatId, latitude, longitude, 0)
}

func (b Bot) ReplyLocation(chatId int, latitude float64, longitude float64, replyToMessageId int) (*Message, error) {
	return b.replyLocation(chatId, latitude, longitude, replyToMessageId)
}
func (b Bot) replyLocation(chatId int, latitude float64, longitude float64, replyToMessageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r := Get(b, "sendLocation", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return b.ParseMessage(r.Result), nil
}

func (b Bot) SendVenue(chatId int, latitude float64, longitude float64, title string, address string) (*Message, error) {
	return b.replyVenue(chatId, latitude, longitude, title, address, 0)

}

func (b Bot) ReplyVenue(chatId int, latitude float64, longitude float64, title string, address string, replyToMessageId int) (*Message, error) {
	return b.replyVenue(chatId, latitude, longitude, title, address, replyToMessageId)
}
func (b Bot) replyVenue(chatId int, latitude float64, longitude float64, title string, address string, replyToMessageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	v.Add("title", title)
	v.Add("address", address)
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r := Get(b, "sendVenue", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	return b.ParseMessage(r.Result), nil
}

func (b Bot) SendContact(chatId int, phoneNumber string, firstName string) (*Message, error) {
	return b.replyContact(chatId, phoneNumber, firstName, 0)
}

func (b Bot) ReplyContact(chatId int, phoneNumber string, firstName string, replyToMessageId int) (*Message, error) {
	return b.replyContact(chatId, phoneNumber, firstName, replyToMessageId)
}

func (b Bot) replyContact(chatId int, phoneNumber string, firstName string, replyToMessageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("phone_number", phoneNumber)
	v.Add("first_name", firstName)
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r := Get(b, "sendContact", v)
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return b.ParseMessage(r.Result), nil
}

// TODO: r.OK or unmarshal??
func (b Bot) SendChatAction(chatId int, phoneNumber string, firstName string) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("phone_number", phoneNumber)
	v.Add("first_name", firstName)

	r := Get(b, "sendChatAction", v)
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	return r.Ok, nil
}
