package Ext

import (
	"strconv"
	"log"
	"net/url"
)

// TODO: Markdown and HTML - two different funcs?
func (b Bot) SendMessage(chat_id int, msg string) Message {
	return b.replyMessage(chat_id, msg, 0)
}

func (b Bot) ReplyMessage(chat_id int, msg string, reply_to_message_id int) Message {
	return b.replyMessage(chat_id, msg, reply_to_message_id)
}

func (b Bot) replyMessage(chat_id int, msg string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("text", msg)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendMessage", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Printf("%+v\n", r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) ForwardMessage(chat_id int, from_chat_id int, message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("from_chat_id", strconv.Itoa(from_chat_id))
	v.Add("message_id", strconv.Itoa(message_id))

	r := Get(b, "forwardMessage", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)
}

// TODO: create InputFile version of all the Str's
func (b Bot) SendPhotoStr(chat_id int, photo string) Message {
	return b.replyPhotoStr(chat_id, photo, "", 0)
}

func (b Bot) ReplyPhotoStr(chat_id int, photo string, reply_to_message_id int) Message {
	return b.replyPhotoStr(chat_id, photo, "", reply_to_message_id)
}

func (b Bot) replyPhotoStr(chat_id int, photo string, caption string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("photo", photo)
	v.Add("caption", caption)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendPhoto", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) SendAudioStr(chat_id int, audio string) Message {
	return b.replyAudioStr(chat_id, audio, 0)
}

func (b Bot) ReplyAudioStr(chat_id int, audio string, reply_to_message_id int) Message {
	return b.replyAudioStr(chat_id, audio, reply_to_message_id)
}

func (b Bot) replyAudioStr(chat_id int, audio string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("audio", audio)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendAudio", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) SendDocumentStr(chat_id int, document string) Message {
	return b.replyDocumentStr(chat_id, document, 0)
}

func (b Bot) ReplyDocumentStr(chat_id int, document string, reply_to_message_id int) Message {
	return b.replyDocumentStr(chat_id, document, reply_to_message_id)
}

func (b Bot) replyDocumentStr(chat_id int, document string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("document", document)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendDocument", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)
}

func (b Bot) SendVideoStr(chat_id int, video string) Message {
	return b.replyVideoStr(chat_id, video, 0)
}

func (b Bot) ReplyVideoStr(chat_id int, video string, reply_to_message_id int) Message {
	return b.replyVideoStr(chat_id, video, reply_to_message_id)
}

func (b Bot) replyVideoStr(chat_id int, video string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("video", video)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendVideo", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)
}

func (b Bot) SendVoiceStr(chat_id int, voice string) Message {
	return b.replyVoiceStr(chat_id, voice, 0)
}

func (b Bot) ReplyVoiceStr(chat_id int, voice string, reply_to_message_id int) Message {
	return b.replyVoiceStr(chat_id, voice, reply_to_message_id)
}

func (b Bot) replyVoiceStr(chat_id int, voice string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("voice", voice)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendVoice", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)
}

func (b Bot) SendVideoNoteStr(chat_id int, video_note string) Message {
	return b.replyVideoNoteStr(chat_id, video_note, 0)
}

func (b Bot) ReplyVideoNoteStr(chat_id int, video_note string, reply_to_message_id int) Message {
	return b.replyVideoNoteStr(chat_id, video_note, reply_to_message_id)
}

func (b Bot) replyVideoNoteStr(chat_id int, video_note string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("video_note", video_note)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendVideoNote", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)
}

func (b Bot) SendLocation(chat_id int, latitude float64, longitude float64) Message {
	return b.replyLocation(chat_id, latitude, longitude, 0)
}

func (b Bot) ReplyLocation(chat_id int, latitude float64, longitude float64, reply_to_message_id int) Message {
	return b.replyLocation(chat_id, latitude, longitude, reply_to_message_id)
}
func (b Bot) replyLocation(chat_id int, latitude float64, longitude float64, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendLocation", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)
}

func (b Bot) SendVenue(chat_id int, latitude float64, longitude float64, title string, address string) Message {
	return b.replyVenue(chat_id, latitude, longitude, title, address, 0)

}

func (b Bot) ReplyVenue(chat_id int, latitude float64, longitude float64, title string, address string, reply_to_message_id int) Message {
	return b.replyVenue(chat_id, latitude, longitude, title, address, reply_to_message_id)
}
func (b Bot) replyVenue(chat_id int, latitude float64, longitude float64, title string, address string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	v.Add("title", title)
	v.Add("address", address)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendVenue", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) SendContact(chat_id int, phone_number string, first_name string) Message {
	return b.replyContact(chat_id, phone_number, first_name, 0)
}

func (b Bot) ReplyContact(chat_id int, phone_number string, first_name string, reply_to_message_id int) Message {
	return b.replyContact(chat_id, phone_number, first_name, reply_to_message_id)
}

func (b Bot) replyContact(chat_id int, phone_number string, first_name string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("phone_number", phone_number)
	v.Add("first_name", first_name)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendContact", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)
}

// TODO: r.OK or unmarshal??
func (b Bot) SendChatAction(chat_id int, phone_number string, first_name string) bool {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("phone_number", phone_number)
	v.Add("first_name", first_name)

	r := Get(b, "sendChatAction", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return r.Ok
}
