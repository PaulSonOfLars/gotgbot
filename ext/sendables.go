package ext

import (
	"net/url"
	"strconv"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"io"
	"os"
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

type file struct {
	Name   string
	FileId string
	Path   string
	Reader io.Reader
	URL    string
}
type sendableTextMessage struct {
	bot                 Bot
	ChatId              int
	Text                string
	ParseMode           string
	DisableWebPreview   bool
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendableTextMessage) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("text", msg.Text)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_web_page_preview", strconv.FormatBool(msg.DisableWebPreview))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := Get(msg.bot, "sendMessage", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendMessage")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendablePhoto struct {
	bot                 Bot
	ChatId              int
	file
	Caption             string
	ParseMode           string
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendablePhoto) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.sendFile(msg.file, "photo","sendPhoto", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendPhoto")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableAudio struct {
	bot                 Bot
	ChatId              int
	file
	Caption             string
	ParseMode           string
	Duration            int
	Performer           string
	Title               string
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendableAudio) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("performer", msg.Performer)
	v.Add("title", msg.Title)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.sendFile(msg.file, "audio","sendAudio", v)

	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendAudio")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableDocument struct {
	bot                 Bot
	ChatId              int
	DocName             string // file name
	file
	Caption             string
	ParseMode           string
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendableDocument) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.sendFile(msg.file, "document","sendDocument", v)

	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendDocument")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableVideo struct {
	bot                 Bot
	ChatId              int
	file
	Duration            int
	Width               int
	Height              int
	Caption             string
	ParseMode           string
	SupportsStreaming   bool
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendableVideo) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("width", strconv.Itoa(msg.Width))
	v.Add("height", strconv.Itoa(msg.Height))
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("supports_streaming", strconv.FormatBool(msg.SupportsStreaming))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.sendFile(msg.file, "video","sendVideo", v)

	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendVideo")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableVoice struct {
	bot                 Bot
	ChatId              int
	file
	Caption             string
	ParseMode           string
	Duration            int
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendableVoice) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.sendFile(msg.file, "voice", "sendVoice", v)

	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendVoice")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableVideoNote struct {
	bot                 Bot
	ChatId              int
	file
	Duration            int
	Length              int
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendableVideoNote) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("length", strconv.Itoa(msg.Length))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.sendFile(msg.file, "video", "sendVideoNote", v)

	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendVideoNote")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

// TODO
type sendableMediaGroup struct {
	bot    Bot
	ChatId int
	//media
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendableMediaGroup) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	log.Println("TODO: media groups") // TODO
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	//v.Add("media")
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := Get(msg.bot, "sendMediaGroup", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendMediaGroup")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
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
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendableLocation) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("latitude", strconv.FormatFloat(msg.Latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(msg.Longitude, 'f', -1, 64))
	v.Add("live_period", strconv.Itoa(msg.LivePeriod))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := Get(msg.bot, "sendLocation", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendLocation")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
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
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendableVenue) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("latitude", strconv.FormatFloat(msg.Latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(msg.Longitude, 'f', -1, 64))
	v.Add("title", msg.Title)
	v.Add("address", msg.Address)
	v.Add("foursquare_id", msg.FoursquareId)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := Get(msg.bot, "sendVenue", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendVenue")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
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
	ReplyMarkup         *ReplyKeyboardMarkup
}

func (msg *sendableContact) Send() (*Message, error) {
	replyMarkup, err := marshalRepyMarkup(msg.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("phone_number", msg.PhoneNumber)
	v.Add("first_name", msg.FirstName)
	v.Add("last_name", msg.LastName)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := Get(msg.bot, "sendContact", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendContact")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.Bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

type sendableChatAction struct {
	bot    Bot
	ChatId int
	Action string
}

func (msg *sendableChatAction) Send() (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("Action", msg.Action)

	r, err := Get(msg.bot, "sendChatAction", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to sendChatAction")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}
	var newMsg bool
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}

func (b Bot) sendFile(msg file, fileType string, endpoint string, params url.Values) (*Response, error) {
	if msg.FileId != "" {
		params.Add(fileType, msg.FileId)
		return Get(b, endpoint, params)
	} else if msg.Path != "" {
		file, err := os.Open(msg.Path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		return Post(b, fileType, endpoint, params, file, msg.Name)
	} else if msg.Reader != nil {
		return Post(b, fileType, endpoint, params, msg.Reader, msg.Name)
	} else {
		return nil, errors.New("the message had no files that could be sent")
	}
}
