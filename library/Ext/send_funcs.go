package Ext

import (
	"bot/library/Types"
	"encoding/json"
	"strconv"
	"log"
)

// TODO: Markdown and HTML - two different funcs?
func (b Bot) SendMessage(msg string, chat_id int) Types.Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["text"] = msg

	r := Get(b, "sendMessage", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}

func (b Bot) ForwardMessage(chat_id int, from_chat_id int, message_id int) Types.Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["from_chat_id"] = strconv.Itoa(from_chat_id)
	m["message_id"] = strconv.Itoa(message_id)

	r := Get(b, "forwardMessage", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}

func (b Bot) SendPhotoStr(chat_id int, photo string) Types.Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["photo"] = photo

	r := Get(b, "sendPhoto", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}

func (b Bot) SendAudioStr(chat_id int, audio string) Types.Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["audio"] = audio

	r := Get(b, "sendAudio", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}

func (b Bot) SendDocumentStr(chat_id int, document string) Types.Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["document"] = document

	r := Get(b, "sendDocument", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}


func (b Bot) SendVideoStr(chat_id int, video string) Types.Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["video"] = video

	r := Get(b, "sendVideo", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}

func (b Bot) SendVoiceStr(chat_id int, voice string) Types.Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["voice"] = voice

	r := Get(b, "sendVoice", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}

func (b Bot) SendVideoNoteStr(chat_id int, note string) Types.Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["video_note"] = note

	r := Get(b, "sendVideoNote", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}

func (b Bot) SendLocation(chat_id int, latitude float64, longitude float64) Types.Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["latitude"] = strconv.FormatFloat(latitude, 'f', -1, 64)
	m["longitude"] = strconv.FormatFloat(longitude, 'f', -1, 64)

	r := Get(b, "sendLocation", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}

func (b Bot) SendVenue(chat_id int, latitude float64, longitude float64, title string, address string) Types.Message {
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

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}

func (b Bot) SendContact(chat_id int, phone_number string, first_name string) Types.Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["phone_number"] = phone_number
	m["first_name"] = first_name

	r := Get(b, "sendContact", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
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
