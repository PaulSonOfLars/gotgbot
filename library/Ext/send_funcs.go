package Ext

import (
	"strconv"
	"log"
)

// TODO: Markdown and HTML - two different funcs?
func (b Bot) SendMessage(chat_id int, msg string) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["text"] = msg

	r := Get(b, "sendMessage", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Printf("%+v\n", r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) ReplyMessage(chat_id int, msg string, reply_to_message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["text"] = msg
	m["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)

	r := Get(b, "sendMessage", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Printf("%+v\n", r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) ForwardMessage(chat_id int, from_chat_id int, message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["from_chat_id"] = strconv.Itoa(from_chat_id)
	m["message_id"] = strconv.Itoa(message_id)

	r := Get(b, "forwardMessage", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)
}

// TODO: create InputFile version of all the Str's
func (b Bot) SendPhotoStr(chat_id int, photo string) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["photo"] = photo

	r := Get(b, "sendPhoto", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) ReplyPhotoStr(chat_id int, photo string, reply_to_message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["photo"] = photo
	m["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)

	r := Get(b, "sendPhoto", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)
}

func (b Bot) SendAudioStr(chat_id int, audio string) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["audio"] = audio

	r := Get(b, "sendAudio", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyAudioStr(chat_id int, audio string, reply_to_message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["audio"] = audio
	m["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)

	r := Get(b, "sendAudio", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)

}

func (b Bot) SendDocumentStr(chat_id int, document string) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["document"] = document

	r := Get(b, "sendDocument", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyDocumentStr(chat_id int, document string, reply_to_message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["document"] = document
	m["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)

	r := Get(b, "sendDocument", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}


func (b Bot) SendVideoStr(chat_id int, video string) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["video"] = video

	r := Get(b, "sendVideo", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyVideoStr(chat_id int, video string, reply_to_message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["video"] = video
	m["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)

	r := Get(b, "sendVideo", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}

func (b Bot) SendVoiceStr(chat_id int, voice string) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["voice"] = voice

	r := Get(b, "sendVoice", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyVoiceStr(chat_id int, voice string, reply_to_message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["voice"] = voice
	m["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)

	r := Get(b, "sendVoice", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}

func (b Bot) SendVideoNoteStr(chat_id int, note string) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["video_note"] = note

	r := Get(b, "sendVideoNote", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyVideoNoteStr(chat_id int, video_note string, reply_to_message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["video_note"] = video_note
	m["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)

	r := Get(b, "sendVideoNote", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}

func (b Bot) SendLocation(chat_id int, latitude float64, longitude float64) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["latitude"] = strconv.FormatFloat(latitude, 'f', -1, 64)
	m["longitude"] = strconv.FormatFloat(longitude, 'f', -1, 64)

	r := Get(b, "sendLocation", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyLocation(chat_id int, latitude float64, longitude float64, reply_to_message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["latitude"] = strconv.FormatFloat(latitude, 'f', -1, 64)
	m["longitude"] = strconv.FormatFloat(longitude, 'f', -1, 64)
	m["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)

	r := Get(b, "sendLocation", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}

func (b Bot) SendVenue(chat_id int, latitude float64, longitude float64, title string, address string) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["latitude"] = strconv.FormatFloat(latitude, 'f', -1, 64)
	m["longitude"] = strconv.FormatFloat(longitude, 'f', -1, 64)
	m["title"] = title
	m["address"] = address

	r := Get(b, "sendVenue", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyVenue(chat_id int, latitude float64, longitude float64, title string, address string, reply_to_message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["latitude"] = strconv.FormatFloat(latitude, 'f', -1, 64)
	m["longitude"] = strconv.FormatFloat(longitude, 'f', -1, 64)
	m["title"] = title
	m["address"] = address
	m["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)

	r := Get(b, "sendVenue", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}
	return b.ParseMessage(r.Result)

}

func (b Bot) SendContact(chat_id int, phone_number string, first_name string) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["phone_number"] = phone_number
	m["first_name"] = first_name

	r := Get(b, "sendContact", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) ReplyContact(chat_id int, phone_number string, first_name string, reply_to_message_id int) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["phone_number"] = phone_number
	m["first_name"] = first_name
	m["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)

	r := Get(b, "sendContact", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}

// TODO: r.OK or unmarshal??
func (b Bot) SendChatAction(chat_id int, phone_number string, first_name string) bool {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["phone_number"] = phone_number
	m["first_name"] = first_name

	r := Get(b, "sendChatAction", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return r.Ok
}
