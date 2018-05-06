package ext

import (
	"net/url"
	"strconv"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
)

const (
	Markdown = "Markdown"
	Html     = "HTML"
)

func (b Bot) NewSendableMessage(chatId int, text string) *sendableTextMessage {
	return &sendableTextMessage{bot: b, ChatId: chatId, Text: text}
}

func (b Bot) NewSendablePhoto(chatId int, caption string) *sendablePhoto {
	return &sendablePhoto{bot: b, ChatId: chatId, Caption: caption}
}

func (b Bot) NewSendableAudio(chatId int, caption string) *sendableAudio {
	return &sendableAudio{bot: b, ChatId: chatId, Caption: caption}
}

func (b Bot) NewSendableDocument(chatId int, caption string) *sendableDocument {
	return &sendableDocument{bot: b, ChatId: chatId, Caption: caption}
}

func (b Bot) NewSendableVideo(chatId int, caption string) *sendableVideo {
	return &sendableVideo{bot: b, ChatId: chatId, Caption: caption}
}

func (b Bot) NewSendableVoice(chatId int, caption string) *sendableVoice {
	return &sendableVoice{bot: b, ChatId: chatId, Caption: caption}
}

func (b Bot) NewSendableVideoNote(chatId int) *sendableVideoNote {
	return &sendableVideoNote{bot: b, ChatId: chatId}
}

func (b Bot) NewSendableMediaGroup(chatId int) *sendableMediaGroup {
	return &sendableMediaGroup{bot: b, ChatId: chatId}
}

func (b Bot) NewSendableLocation(chatId int) *sendableLocation {
	return &sendableLocation{bot: b, ChatId: chatId}
}

func (b Bot) NewSendableVenue(chatId int) *sendableVenue {
	return &sendableVenue{bot: b, ChatId: chatId}
}

func (b Bot) NewSendableContact(chatId int) *sendableContact {
	return &sendableContact{bot: b, ChatId: chatId}
}

func (b Bot) NewSendableChatAction(chatId int) *sendableChatAction {
	return &sendableChatAction{bot: b, ChatId: chatId}
}

type sendableTextMessage struct {
	bot                 Bot
	ChatId              int
	Text                string
	ParseMode           string
	DisableWebPreview   bool
	DisableNotification bool
	ReplyToMessageId    int
	// replyMarkup
}

func (msg *sendableTextMessage) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("text", msg.Text)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_web_page_preview", strconv.FormatBool(msg.DisableWebPreview))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendMessage", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendablePhoto struct {
	bot         Bot
	ChatId      int
	PhotoString string
	//photoFile
	Caption             string
	ParseMode           string
	DisableNotification bool
	ReplyToMessageId    int
	//replyMarkup
}

func (msg *sendablePhoto) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	if msg.PhotoString == "" {
		v.Add("photo", msg.PhotoString)

	} else {
		// TODO figure out type here
		log.Println("TODO: implement photofiles")
	}

	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendPhoto", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableAudio struct {
	bot         Bot
	ChatId      int
	AudioString string
	//audioFile
	Caption             string
	ParseMode           string
	Duration            int
	Performer           string
	Title               string
	DisableNotification bool
	ReplyToMessageId    int
	//replyMarkup
}

func (msg *sendableAudio) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	if msg.AudioString == "" {
		v.Add("audio", msg.AudioString)
	} else {
		// TODO figure out type here
		log.Println("TODO: implement audiofiles")
	}

	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("performer", msg.Performer)
	v.Add("title", msg.Title)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendAudio", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableDocument struct {
	bot       Bot
	ChatId    int
	DocString string
	//docFile
	Caption             string
	ParseMode           string
	DisableNotification bool
	ReplyToMessageId    int
	//replyMarkup
}

func (msg *sendableDocument) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	if msg.DocString == "" {
		v.Add("document", msg.DocString)
	} else {
		// TODO figure out type here
		log.Println("TODO: implement docfiles")
	}
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendDocument", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableVideo struct {
	bot         Bot
	ChatId      int
	VideoString string
	//videoFile
	Duration            int
	Width               int
	Height              int
	Caption             string
	ParseMode           string
	SupportsStreaming   bool
	DisableNotification bool
	ReplyToMessageId    int
	//replyMarkup
}

func (msg *sendableVideo) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	if msg.VideoString == "" {
		v.Add("video", msg.VideoString)
	} else {
		// TODO figure out type here
		log.Println("TODO: implement videofiles")
	}
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("width", strconv.Itoa(msg.Width))
	v.Add("height", strconv.Itoa(msg.Height))
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("supports_streaming", strconv.FormatBool(msg.SupportsStreaming))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendMessage", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableVoice struct {
	bot         Bot
	ChatId      int
	VoiceString string
	//voiceFile
	Caption             string
	ParseMode           string
	Duration            int
	DisableNotification bool
	ReplyToMessageId    int
	//replyMarkup
}

func (msg *sendableVoice) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	if msg.VoiceString == "" {
		v.Add("voice", msg.VoiceString)
	} else {
		// TODO figure out type here
		log.Println("TODO: implement voicefiles")
	}
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendMessage", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableVideoNote struct {
	bot             Bot
	ChatId          int
	VideoNoteString string
	//videoNoteFile
	Duration            int
	Length              int
	DisableNotification bool
	ReplyToMessageId    int
	//replyMarkup
}

func (msg *sendableVideoNote) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	if msg.VideoNoteString == "" {
		v.Add("video_note", msg.VideoNoteString)
	} else {
		// TODO figure out type here
		log.Println("TODO: implement videonotefiles")
	}
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("length", strconv.Itoa(msg.Length))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendMessage", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableMediaGroup struct {
	bot    Bot
	ChatId int
	//media
	DisableNotification bool
	ReplyToMessageId    int
	//replyMarkup
}

func (msg *sendableMediaGroup) Send() (*Message, error) {
	log.Println("TODO: media groups") // TODO
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	//v.Add("media")
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendMessage", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableLocation struct {
	bot                 Bot
	ChatId              int
	Latitude            float64
	Longitude           float64
	LivePeriod          int
	DisableNotification bool
	ReplyToMessageId    int
	//replyMarkup
}

func (msg *sendableLocation) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("latitude", strconv.FormatFloat(msg.Latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(msg.Longitude, 'f', -1, 64))
	v.Add("live_period", strconv.Itoa(msg.LivePeriod))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendMessage", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

//TODO: edit live location
//TODO: stop live location

type sendableVenue struct {
	bot                 Bot
	ChatId              int
	Latitude            float64
	Longitude           float64
	Title               string
	Address             string
	FoursquareId        string
	DisableNotification bool
	ReplyToMessageId    int
	//replyMarkup
}

func (msg *sendableVenue) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("latitude", strconv.FormatFloat(msg.Latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(msg.Longitude, 'f', -1, 64))
	v.Add("title", msg.Title)
	v.Add("address", msg.Address)
	v.Add("foursquare_id", msg.FoursquareId)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendMessage", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableContact struct {
	bot                 Bot
	ChatId              int
	PhoneNumber         string
	FirstName           string
	LastName            string
	DisableNotification bool
	ReplyToMessageId    int
	//replyMarkup
}

func (msg *sendableContact) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("phone_number", msg.PhoneNumber)
	v.Add("first_name", msg.FirstName)
	v.Add("last_name", msg.LastName)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	// v.Add("reply_markup", "")

	r := Get(msg.bot, "sendMessage", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableChatAction struct {
	bot    Bot
	ChatId int
	action string
}

func (msg *sendableChatAction) Send() (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("action", msg.action)

	r := Get(msg.bot, "sendChatAction", v)

	if !r.Ok {
		return false, errors.New(r.Description)
	}
	var newMsg bool
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}
