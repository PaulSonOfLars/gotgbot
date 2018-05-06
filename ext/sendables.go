package ext

import (
	"net/url"
	"strconv"
	"encoding/json"
	"gotgbot/types"
	"github.com/pkg/errors"
	"log"
)

type sendableTextMessage struct {
	bot              Bot
	chatId           int
	text             string
	parseMode        string
	webPreview       bool
	notification     bool
	replyToMessageId int
	// replyMarkup
}

const (
	Markdown = "Markdown"
	Html     = "HTML"
)

func (b Bot) NewSendableMessage(chatId int, text string) *sendableTextMessage {
	return &sendableTextMessage{bot: b, chatId: chatId, text: text}
}

func (b Bot) NewSendablePhoto(chatId int, caption string) *sendablePhoto {
	return &sendablePhoto{bot: b, chatId: chatId, caption: caption}
}

func (b Bot) NewSendableAudio(chatId int, caption string) *sendableAudio {
	return &sendableAudio{bot: b, chatId: chatId, caption: caption}
}

func (b Bot) NewSendableDocument(chatId int, caption string) *sendableDocument {
	return &sendableDocument{bot: b, chatId: chatId, caption: caption}
}

func (b Bot) NewSendableVideo(chatId int, caption string) *sendableVideo {
	return &sendableVideo{bot: b, chatId: chatId, caption: caption}
}

func (b Bot) NewSendableVoice(chatId int, caption string) *sendableVoice {
	return &sendableVoice{bot: b, chatId: chatId, caption: caption}
}

func (b Bot) NewSendableVideoNote(chatId int) *sendableVideoNote {
	return &sendableVideoNote{bot: b, chatId: chatId}
}

func (b Bot) NewSendableMediaGroup(chatId int) *sendableMediaGroup {
	return &sendableMediaGroup{bot: b, chatId: chatId}
}

func (b Bot) NewSendableLocation(chatId int) *sendableLocation {
	return &sendableLocation{bot: b, chatId: chatId}
}

func (b Bot) NewSendableVenue(chatId int) *sendableVenue {
	return &sendableVenue{bot: b, chatId: chatId}
}

func (b Bot) NewSendableContact(chatId int) *sendableContact {
	return &sendableContact{bot: b, chatId: chatId}
}

func (b Bot) NewSendableChatAction(chatId int) *sendableChatAction {
	return &sendableChatAction{bot: b, chatId: chatId}
}

func (msg *sendableTextMessage) SetParseMode(parseMode string) {
	msg.parseMode = parseMode
}

func (msg *sendableTextMessage) ReplyToMsg(replyTo *types.Message) {
	msg.replyToMessageId = replyTo.Message_id
}

func (msg *sendableTextMessage) ReplyToMsgId(replyTo int) {
	msg.replyToMessageId = replyTo
}

func (msg *sendableTextMessage) SetWebPagePreview(enable bool) {
	msg.webPreview = enable
}

func (msg *sendableTextMessage) SetNotification(enable bool) {
	msg.notification = enable
}

func (msg *sendableTextMessage) addKeyboard() {
	log.Println("ERROR: addkeyboard not implemented")
}

func (msg *sendableTextMessage) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	v.Add("text", msg.text)
	v.Add("parse_mode", msg.parseMode)
	v.Add("disable_web_page_preview", strconv.FormatBool(msg.webPreview))
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	chatId      int
	photoString string
	//photoFile
	caption          string
	parseMode        string
	notification     bool
	replyToMessageId int
	//replyMarkup
}

func (msg *sendablePhoto) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	if msg.photoString == "" {
		v.Add("photo", msg.photoString)

	} else {
		// TODO figure out type here
		log.Println("TODO: implement photofiles")
	}

	v.Add("caption", msg.caption)
	v.Add("parse_mode", msg.parseMode)
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	chatId      int
	audioString string
	//audioFile
	caption          string
	parseMode        string
	duration         int
	performer        string
	title            string
	notification     bool
	replyToMessageId int
	//replyMarkup
}

func (msg *sendableAudio) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	if msg.audioString == "" {
		v.Add("audio", msg.audioString)
	} else {
		// TODO figure out type here
		log.Println("TODO: implement audiofiles")
	}

	v.Add("caption", msg.caption)
	v.Add("parse_mode", msg.parseMode)
	v.Add("duration", strconv.Itoa(msg.duration))
	v.Add("performer", msg.performer)
	v.Add("title", msg.title)
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	chatId    int
	docString string
	//docFile
	caption          string
	parseMode        string
	notification     bool
	replyToMessageId int
	//replyMarkup
}

func (msg *sendableDocument) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	if msg.docString == "" {
		v.Add("document", msg.docString)
	} else {
		// TODO figure out type here
		log.Println("TODO: implement docfiles")
	}
	v.Add("caption", msg.caption)
	v.Add("parse_mode", msg.parseMode)
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	chatId      int
	videoString string
	//videoFile
	duration          int
	width             int
	height            int
	caption           string
	parseMode         string
	supportsStreaming bool
	notification      bool
	replyToMessageId  int
	//replyMarkup
}

func (msg *sendableVideo) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	if msg.videoString == "" {
		v.Add("video", msg.videoString)
	} else {
		// TODO figure out type here
		log.Println("TODO: implement videofiles")
	}
	v.Add("duration", strconv.Itoa(msg.duration))
	v.Add("width", strconv.Itoa(msg.width))
	v.Add("height", strconv.Itoa(msg.height))
	v.Add("caption", msg.caption)
	v.Add("parse_mode", msg.parseMode)
	v.Add("supports_streaming", strconv.FormatBool(msg.supportsStreaming))
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	chatId      int
	voiceString string
	//voiceFile
	caption          string
	parseMode        string
	duration         int
	notification     bool
	replyToMessageId int
	//replyMarkup
}

func (msg *sendableVoice) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	if msg.voiceString == "" {
		v.Add("voice", msg.voiceString)
	} else {
		// TODO figure out type here
		log.Println("TODO: implement voicefiles")
	}
	v.Add("caption", msg.caption)
	v.Add("parse_mode", msg.parseMode)
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	chatId          int
	videoNoteString string
	//videoNoteFile
	duration         int
	length           int
	notification     bool
	replyToMessageId int
	//replyMarkup
}

func (msg *sendableVideoNote) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	if msg.videoNoteString== "" {
		v.Add("video_note", msg.videoNoteString)
	} else {
		// TODO figure out type here
		log.Println("TODO: implement videonotefiles")
	}
	v.Add("duration", strconv.Itoa(msg.duration))
	v.Add("length", strconv.Itoa(msg.length))
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	chatId int
	//media
	notification     bool
	replyToMessageId int
	//replyMarkup
}

func (msg *sendableMediaGroup) Send() (*Message, error) {
	log.Println("TODO: media groups") // TODO
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	//v.Add("media")
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	bot              Bot
	chatId           int
	latitude         float64
	longitude        float64
	livePeriod       int
	notification     bool
	replyToMessageId int
	//replyMarkup
}

func (msg *sendableLocation) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	v.Add("latitude", strconv.FormatFloat(msg.latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(msg.longitude, 'f', -1, 64))
	v.Add("live_period", strconv.Itoa(msg.livePeriod))
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	bot              Bot
	chatId           int
	latitude         float64
	longitude        float64
	title            string
	address          string
	foursquareId     string
	notification     bool
	replyToMessageId int
	//replyMarkup
}

func (msg *sendableVenue) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	v.Add("latitude", strconv.FormatFloat(msg.latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(msg.longitude, 'f', -1, 64))
	v.Add("title", msg.title)
	v.Add("address", msg.address)
	v.Add("foursquare_id", msg.foursquareId)
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	bot              Bot
	chatId           int
	phoneNumber      string
	firstName        string
	lastName         string
	notification     bool
	replyToMessageId int
	//replyMarkup
}

func (msg *sendableContact) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	v.Add("phone_number", msg.phoneNumber)
	v.Add("first_name", msg.firstName)
	v.Add("last_name", msg.lastName)
	v.Add("disable_notification", strconv.FormatBool(msg.notification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
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
	chatId int
	action string
}

func (msg *sendableChatAction) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	v.Add("action", msg.action)

	r := Get(msg.bot, "sendChatAction", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}
