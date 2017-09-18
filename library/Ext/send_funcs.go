package Ext

import (
	"strconv"
	"log"
	"net/url"
)

// TODO: Markdown and HTML - two different funcs?
func (b Bot) SendMessage(chat_id int, msg string) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("text", msg)

	r := Get(b, "sendMessage", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Printf("%+v\n", r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) ReplyMessage(chat_id int, msg string, reply_to_message_id int) Message {
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
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("photo", photo)

	r := Get(b, "sendPhoto", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) ReplyPhotoStr(chat_id int, photo string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("photo", photo)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendPhoto", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) SendAudioStr(chat_id int, audio string) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("audio", audio)

	r := Get(b, "sendAudio", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyAudioStr(chat_id int, audio string, reply_to_message_id int) Message {
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
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("document", document)

	r := Get(b, "sendDocument", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyDocumentStr(chat_id int, document string, reply_to_message_id int) Message {
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
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("video", video)

	r := Get(b, "sendVideo", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyVideoStr(chat_id int, video string, reply_to_message_id int) Message {
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
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("voice", voice)

	r := Get(b, "sendVoice", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyVoiceStr(chat_id int, voice string, reply_to_message_id int) Message {
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

func (b Bot) SendVideoNoteStr(chat_id int, note string) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("video_note", note)

	r := Get(b, "sendVideoNote", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyVideoNoteStr(chat_id int, video_note string, reply_to_message_id int) Message {
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
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))

	r := Get(b, "sendLocation", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyLocation(chat_id int, latitude float64, longitude float64, reply_to_message_id int) Message {
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
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	v.Add("title", title)
	v.Add("address", address)

	r := Get(b, "sendVenue", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyVenue(chat_id int, latitude float64, longitude float64, title string, address string, reply_to_message_id int) Message {
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
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("phone_number", phone_number)
	v.Add("first_name", first_name)

	r := Get(b, "sendContact", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyContact(chat_id int, phone_number string, first_name string, reply_to_message_id int) Message {
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
