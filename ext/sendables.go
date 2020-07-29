package ext

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type Sendable interface {
	Send() (*Message, error)
}

// NewSendableMessage create a new message struct to send
func (b Bot) NewSendableMessage(chatId int, text string) *sendableTextMessage {
	return &sendableTextMessage{
		bot:               b,
		ChatId:            chatId,
		Text:              text,
		DisableWebPreview: b.DisableWebPreview,
	}
}

// NewSendableEditMessageText create a new message editing struct to send
func (b Bot) NewSendableEditMessageText(chatId int, messageId int, text string) *sendableEditMessageText {
	return &sendableEditMessageText{
		bot:               b,
		ChatId:            chatId,
		MessageId:         messageId,
		Text:              text,
		DisableWebPreview: b.DisableWebPreview,
	}
}

// NewSendableEditMessageCaption create a new caption editing struct to send
func (b Bot) NewSendableEditMessageCaption(chatId int, messageId int, caption string) *sendableEditMessageCaption {
	return &sendableEditMessageCaption{
		bot:       b,
		ChatId:    chatId,
		MessageId: messageId,
		Caption:   caption,
	}
}

// NewSendableEditMessageReplyMarkup creates a new markup editing struct to send
func (b Bot) NewSendableEditMessageReplyMarkup(chatId int, messageId int, markup ReplyMarkup) *sendableEditMessageReplyMarkup {
	return &sendableEditMessageReplyMarkup{
		bot:         b,
		ChatId:      chatId,
		MessageId:   messageId,
		ReplyMarkup: markup,
	}
}

// NewSendablePhoto creates a new photo struct to send
func (b Bot) NewSendablePhoto(chatId int, caption string) *sendablePhoto {
	return &sendablePhoto{bot: b, ChatId: chatId, Caption: caption}
}

// NewSendableAudio creates a new audio struct to send
func (b Bot) NewSendableAudio(chatId int, caption string) *sendableAudio {
	return &sendableAudio{bot: b, ChatId: chatId, Caption: caption}
}

// NewSendableDocument creates a new document struct to send
func (b Bot) NewSendableDocument(chatId int, caption string) *sendableDocument {
	return &sendableDocument{bot: b, ChatId: chatId, Caption: caption}
}

// NewSendableVideo creates a new video struct to send
func (b Bot) NewSendableVideo(chatId int, caption string) *sendableVideo {
	return &sendableVideo{bot: b, ChatId: chatId, Caption: caption}
}

// NewSendableVoice creates a new voice struct to send
func (b Bot) NewSendableVoice(chatId int, caption string) *sendableVoice {
	return &sendableVoice{bot: b, ChatId: chatId, Caption: caption}
}

// NewSendableVideoNote creates a new videonote struct to send
func (b Bot) NewSendableVideoNote(chatId int) *sendableVideoNote {
	return &sendableVideoNote{bot: b, ChatId: chatId}
}

// NewSendableMediaGroup creates a new mediagroup struct to send
func (b Bot) NewSendableMediaGroup(chatId int) *sendableMediaGroup {
	return &sendableMediaGroup{bot: b, ChatId: chatId}
}

// NewSendableEditMessageMedia creates a new editmessage media struct to send
func (b Bot) NewSendableEditMessageMedia(chatId int, messageId int) *sendableEditMessageMedia {
	return &sendableEditMessageMedia{
		bot:       b,
		ChatId:    chatId,
		MessageId: messageId,
	}
}

func (b Bot) NewSendableLocation(chatId int) *sendableLocation {
	return &sendableLocation{bot: b, ChatId: chatId}
}

func (b Bot) NewSendableEditMessageLiveLocation(chatId int, latitude float64, longitude float64) *sendableEditMessageLiveLocation {
	return &sendableEditMessageLiveLocation{bot: b, ChatId: chatId, Latitude: latitude, Longitude: longitude}
}

func (b Bot) NewSendableStopMessageLiveLocation(chatId int) *sendableStopMessageLiveLocation {
	return &sendableStopMessageLiveLocation{bot: b, ChatId: chatId}
}

// NewSendableVenue creates a new venue struct to send
func (b Bot) NewSendableVenue(chatId int) *sendableVenue {
	return &sendableVenue{bot: b, ChatId: chatId}
}

// NewSendableContact creates a new contact struct to send
func (b Bot) NewSendableContact(chatId int) *sendableContact {
	return &sendableContact{bot: b, ChatId: chatId}
}

// NewSendableChatAction creates a new chat action struct to send
func (b Bot) NewSendableChatAction(chatId int) *sendableChatAction {
	return &sendableChatAction{bot: b, ChatId: chatId}
}

// NewSendableAnimation creates a new animation struct to send
func (b Bot) NewSendableAnimation(chatId int) *sendableAnimation {
	return &sendableAnimation{bot: b, ChatId: chatId}
}

// NewSendablePoll creates a new poll struct to send.
func (b Bot) NewSendablePoll(chatId int, question string, options []string) *sendablePoll {
	return &sendablePoll{bot: b, ChatId: chatId, Question: question, Options: options}
}

// NewSendableDice creates a new poll struct to send.
func (b Bot) NewSendableDice(chatId int) *sendableDice {
	return &sendableDice{bot: b, ChatId: chatId}
}

// NewSendableAnswerCallbackQuery creates a new callbackQuery struct to send.
func (b Bot) NewSendableAnswerCallbackQuery(queryId string) *sendableCallbackQuery {
	return &sendableCallbackQuery{bot: b, CallbackQueryId: queryId}
}

type InputFile struct {
	b      Bot
	Name   string
	FileId string
	Reader io.Reader
	URL    string
}

func (f InputFile) send(endpoint string, params url.Values, fileType string) (json.RawMessage, error) {
	params.Add("name", f.Name)
	if f.FileId != "" {
		params.Add(fileType, f.FileId)
		return f.b.Get(endpoint, params)
	} else if f.URL != "" {
		params.Add(fileType, f.URL)
		return f.b.Get(endpoint, params)
	} else if f.Reader != nil {
		return f.b.Post(endpoint, params, map[string]PostFile{
			fileType: {
				File:     f.Reader,
				FileName: f.Name,
			}})
	} else {
		return nil, errors.New("the message had no files that could be sent")
	}
}

func (f InputFile) GetMediaType(attachAs string, data map[string]PostFile) string {
	if f.FileId != "" {
		return f.FileId
	} else if f.URL != "" {
		return f.URL
	} else if f.Reader != nil {
		data[attachAs] = PostFile{
			File:     f.Reader,
			FileName: f.Name,
		}
		return "attach://" + attachAs
	}
	return ""
}

func (b Bot) NewFileId(fileId string) InputFile {
	return InputFile{
		b:      b,
		FileId: fileId,
	}
}

func (b Bot) NewFileURL(url string) InputFile {
	return InputFile{
		b:   b,
		URL: url,
	}
}

func (b Bot) NewFileReader(name string, r io.Reader) InputFile {
	if name == "" {
		name = "file"
	}
	return InputFile{
		b:      b,
		Name:   name,
		Reader: r,
	}
}

type InputMedia interface {
	toJson(idx int) (map[string]string, map[string]PostFile)
}

type InputMediaAnimation struct {
	Media     InputFile
	Caption   string
	ParseMode string

	Thumb    *InputFile
	Width    int
	Height   int
	Duration int
}

func (ima InputMediaAnimation) toJson(idx int) (map[string]string, map[string]PostFile) {
	data := make(map[string]PostFile)
	media := ima.Media.GetMediaType("media"+strconv.Itoa(idx), data)
	m := map[string]string{
		"type":       "animation",
		"media":      media,
		"caption":    ima.Caption,
		"parse_mode": ima.ParseMode,
		"width":      strconv.Itoa(ima.Width),
		"height":     strconv.Itoa(ima.Height),
		"duration":   strconv.Itoa(ima.Duration),
	}
	if ima.Thumb != nil {
		m["thumb"] = ima.Thumb.GetMediaType("thumb"+strconv.Itoa(idx), data)
	}
	return m, data
}

type InputMediaDocument struct {
	Media     InputFile
	Caption   string
	ParseMode string

	Thumb *InputFile
}

func (imd InputMediaDocument) toJson(idx int) (map[string]string, map[string]PostFile) {
	data := make(map[string]PostFile)
	media := imd.Media.GetMediaType("media"+strconv.Itoa(idx), data)
	m := map[string]string{
		"type":       "document",
		"media":      media,
		"caption":    imd.Caption,
		"parse_mode": imd.ParseMode,
	}
	if imd.Thumb != nil {
		m["thumb"] = imd.Thumb.GetMediaType("thumb"+strconv.Itoa(idx), data)
	}
	return m, data
}

type InputMediaAudio struct {
	Media     InputFile
	Caption   string
	ParseMode string
	Thumb     *InputFile
	Duration  int
	Performer string
	Title     string
}

func (ima InputMediaAudio) toJson(idx int) (map[string]string, map[string]PostFile) {
	data := make(map[string]PostFile)
	media := ima.Media.GetMediaType("media"+strconv.Itoa(idx), data)
	m := map[string]string{
		"type":       "audio",
		"media":      media,
		"caption":    ima.Caption,
		"parse_mode": ima.ParseMode,
		"duration":   strconv.Itoa(ima.Duration),
		"performer":  ima.Performer,
		"title":      ima.Title,
	}
	if ima.Thumb != nil {
		m["thumb"] = ima.Thumb.GetMediaType("thumb"+strconv.Itoa(idx), data)
	}
	return m, data
}

type InputMediaPhoto struct {
	Media     InputFile
	Caption   string
	ParseMode string
}

func (imp InputMediaPhoto) toJson(idx int) (map[string]string, map[string]PostFile) {
	data := make(map[string]PostFile)
	media := imp.Media.GetMediaType("media"+strconv.Itoa(idx), data)
	return map[string]string{
		"type":       "photo",
		"media":      media,
		"caption":    imp.Caption,
		"parse_mode": imp.ParseMode,
	}, data
}

type InputMediaVideo struct {
	Media            InputFile
	Caption          string
	ParseMode        string
	Thumb            *InputFile
	Width            int
	Height           int
	Duration         int
	SupportStreaming bool
}

func (imv InputMediaVideo) toJson(idx int) (map[string]string, map[string]PostFile) {
	data := make(map[string]PostFile)
	media := imv.Media.GetMediaType("media"+strconv.Itoa(idx), data)
	m := map[string]string{
		"type":              "video",
		"media":             media,
		"caption":           imv.Caption,
		"parse_mode":        imv.ParseMode,
		"width":             strconv.Itoa(imv.Width),
		"height":            strconv.Itoa(imv.Height),
		"duration":          strconv.Itoa(imv.Duration),
		"support_streaming": strconv.FormatBool(imv.SupportStreaming),
	}
	if imv.Thumb != nil {
		m["thumb"] = imv.Thumb.GetMediaType("thumb"+strconv.Itoa(idx), data)
	}
	return m, data
}

type sendableTextMessage struct {
	bot                 Bot
	ChatId              int
	Text                string
	ParseMode           string
	DisableWebPreview   bool
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableTextMessage) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("text", msg.Text)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_web_page_preview", strconv.FormatBool(msg.DisableWebPreview))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.Get("sendMessage", v)
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableEditMessageText struct {
	bot               Bot
	ChatId            int
	MessageId         int
	InlineMessageId   string
	Text              string
	ParseMode         string
	DisableWebPreview bool
	ReplyMarkup       ReplyMarkup
}

func (msg *sendableEditMessageText) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("message_id", strconv.Itoa(msg.MessageId))
	v.Add("inline_message_id", msg.InlineMessageId)
	v.Add("text", msg.Text)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_web_page_preview", strconv.FormatBool(msg.DisableWebPreview))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.Get("editMessageText", v)
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableEditMessageCaption struct {
	bot             Bot
	ChatId          int
	MessageId       int
	InlineMessageId string
	Caption         string
	ParseMode       string
	ReplyMarkup     ReplyMarkup
}

func (msg *sendableEditMessageCaption) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("message_id", strconv.Itoa(msg.MessageId))
	v.Add("inline_message_id", msg.InlineMessageId)
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.Get("editMessageCaption", v)
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableEditMessageReplyMarkup struct {
	bot             Bot
	ChatId          int
	MessageId       int
	InlineMessageId string
	ReplyMarkup     ReplyMarkup
}

func (msg *sendableEditMessageReplyMarkup) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("message_id", strconv.Itoa(msg.MessageId))
	v.Add("inline_message_id", msg.InlineMessageId)
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.Get("editMessageCaption", v)
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendablePhoto struct {
	bot                 Bot
	ChatId              int
	Photo               InputFile
	Caption             string
	ParseMode           string
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendablePhoto) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.Photo.send("sendPhoto", v, "photo")
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableAudio struct {
	bot                 Bot
	ChatId              int
	Audio               InputFile
	Caption             string
	ParseMode           string
	Duration            int
	Performer           string
	Title               string
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableAudio) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
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

	r, err := msg.Audio.send("sendAudio", v, "audio")
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableDocument struct {
	bot                 Bot
	ChatId              int
	Document            InputFile
	Caption             string
	ParseMode           string
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableDocument) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.Document.send("sendDocument", v, "document")
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableVideo struct {
	bot                 Bot
	ChatId              int
	Video               InputFile
	Duration            int
	Width               int
	Height              int
	Caption             string
	ParseMode           string
	SupportsStreaming   bool
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableVideo) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
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

	r, err := msg.Video.send("sendVideo", v, "video")
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableVoice struct {
	bot                 Bot
	ChatId              int
	Voice               InputFile
	Caption             string
	ParseMode           string
	Duration            int
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableVoice) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.Voice.send("sendVoice", v, "voice")
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableVideoNote struct {
	bot                 Bot
	ChatId              int
	VideoNote           InputFile
	Duration            int
	Length              int
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableVideoNote) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("length", strconv.Itoa(msg.Length))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.VideoNote.send("sendVideoNote", v, "video_note")
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableEditMessageMedia struct {
	bot             Bot
	ChatId          int
	MessageId       int
	InlineMessageId string
	InputMedia      InputMedia
	ReplyMarkup     ReplyMarkup
}

func (msg *sendableEditMessageMedia) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("message_id", strconv.Itoa(msg.MessageId))
	v.Add("inline_message_id", msg.InlineMessageId)
	v.Add("reply_markup", string(replyMarkup))

	mediaMap, data := msg.InputMedia.toJson(0)

	bs, err := json.Marshal(mediaMap)
	v.Add("media", string(bs))

	var r json.RawMessage
	if len(data) != 0 {
		r, err = msg.bot.Post("editMessageMedia", v, data)
	} else {
		r, err = msg.bot.Get("editMessageMedia", v)
	}
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableMediaGroup struct {
	bot                 Bot
	ChatId              int
	ArrInputMedia       []InputMedia
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableMediaGroup) Send() (*Message, error) {
	var replyMarkup []byte
	var err error
	if msg.ReplyMarkup != nil {
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	var media []byte
	data := make(map[string]PostFile)
	if msg.ArrInputMedia != nil {
		array := make([][]byte, len(msg.ArrInputMedia))
		for idx, mediaItem := range msg.ArrInputMedia {
			mediaMap, mediaData := mediaItem.toJson(idx)
			array[idx], err = json.Marshal(mediaMap)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to marshal item in array input media")
			}
			// assign all data elems
			for k, v := range mediaData {
				data[k] = v
			}
		}
		vals, err := json.Marshal(array)
		if err != nil {
			return nil, err
		}
		media = vals
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("media", string(media))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	var r json.RawMessage
	if len(data) != 0 {
		r, err = msg.bot.Post("sendMediaGroup", v, data)
	} else {
		r, err = msg.bot.Get("sendMediaGroup", v)
	}
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableLocation struct {
	bot                 Bot
	ChatId              int
	Latitude            float64
	Longitude           float64
	LivePeriod          int
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableLocation) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("latitude", strconv.FormatFloat(msg.Latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(msg.Longitude, 'f', -1, 64))
	v.Add("live_period", strconv.Itoa(msg.LivePeriod))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.Get("sendLocation", v)
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableEditMessageLiveLocation struct {
	bot             Bot
	ChatId          int
	MessageId       int
	InlineMessageId string
	Latitude        float64
	Longitude       float64
	ReplyMarkup     ReplyMarkup
}

func (msg *sendableEditMessageLiveLocation) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("message_id", strconv.Itoa(msg.MessageId))
	v.Add("inline_message_id", msg.InlineMessageId)
	v.Add("latitude", strconv.FormatFloat(msg.Latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(msg.Longitude, 'f', -1, 64))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.Get("editMessageLiveLocation", v)
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableStopMessageLiveLocation struct {
	bot             Bot
	ChatId          int
	MessageId       int
	InlineMessageId string
	ReplyMarkup     ReplyMarkup
}

func (msg *sendableStopMessageLiveLocation) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("message_id", strconv.Itoa(msg.MessageId))
	v.Add("inline_message_id", msg.InlineMessageId)
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.Get("stopMessageLiveLocation", v)
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

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
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableVenue) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
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

	r, err := msg.bot.Get("sendVenue", v)
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableContact struct {
	bot                 Bot
	ChatId              int
	PhoneNumber         string
	FirstName           string
	LastName            string
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableContact) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("phone_number", msg.PhoneNumber)
	v.Add("first_name", msg.FirstName)
	v.Add("last_name", msg.LastName)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.Get("sendContact", v)
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
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

	r, err := msg.bot.Get("sendChatAction", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

type sendableAnimation struct {
	bot       Bot
	ChatId    int
	Animation InputFile
	Duration  int
	Width     int
	Height    int
	// Thumb // TODO: support this
	Caption             string
	ParseMode           string
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (msg *sendableAnimation) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("duration", strconv.Itoa(msg.Duration))
	v.Add("width", strconv.Itoa(msg.Width))
	v.Add("height", strconv.Itoa(msg.Height))
	// v.Add("thumb", msg.Thumb)
	v.Add("caption", msg.Caption)
	v.Add("parse_mode", msg.ParseMode)
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.Animation.send("sendAnimation", v, "animation")
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendablePoll struct {
	bot                   Bot
	ChatId                int
	Question              string
	Options               []string
	IsAnonymous           bool
	Type                  string
	AllowsMultipleAnswers bool
	CorrectOptionId       int
	Explanation           string
	ExplanationParseMode  string
	OpenPeriod            int
	CloseDate             int
	IsClosed              bool
	DisableNotification   bool
	ReplyToMessageId      int
	ReplyMarkup           ReplyMarkup
}

func (msg *sendablePoll) Send() (*Message, error) {
	var replyMarkup []byte
	if msg.ReplyMarkup != nil {
		var err error
		replyMarkup, err = msg.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}
	if msg.Options == nil {
		msg.Options = []string{}
	}

	optionsBytes, err := json.Marshal(msg.Options)
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.ChatId))
	v.Add("question", msg.Question)
	v.Add("options", string(optionsBytes))
	v.Add("is_anonymous", strconv.FormatBool(msg.IsAnonymous))
	v.Add("type", msg.Type)
	v.Add("allows_multiple_answers", strconv.FormatBool(msg.AllowsMultipleAnswers))
	v.Add("correct_option_id", strconv.Itoa(msg.CorrectOptionId))
	v.Add("explanation", msg.Explanation)
	v.Add("explanation_parse_mode", msg.ExplanationParseMode)
	v.Add("open_period", strconv.Itoa(msg.OpenPeriod))
	v.Add("close_date", strconv.Itoa(msg.CloseDate))
	v.Add("is_closed", strconv.FormatBool(msg.IsClosed))
	v.Add("disable_notification", strconv.FormatBool(msg.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(msg.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := msg.bot.Get("sendPoll", v)
	if err != nil {
		return nil, err
	}

	return msg.bot.ParseMessage(r)
}

type sendableDice struct {
	bot                 Bot
	ChatId              int
	Emoji               string
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (d *sendableDice) Send() (*Message, error) {
	var replyMarkup []byte
	if d.ReplyMarkup != nil {
		var err error
		replyMarkup, err = d.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(d.ChatId))
	v.Add("emoji", d.Emoji)
	v.Add("disable_notification", strconv.FormatBool(d.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(d.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := d.bot.Get("sendDice", v)
	if err != nil {
		return nil, err
	}

	return d.bot.ParseMessage(r)
}

type sendableCallbackQuery struct {
	bot             Bot
	CallbackQueryId string
	Text            string
	ShowAlert       bool
	Url             string
	CacheTime       int
}

func (cbq *sendableCallbackQuery) Send() (bool, error) {
	v := url.Values{}
	v.Add("callback_query_id", cbq.CallbackQueryId)
	v.Add("text", cbq.Text)
	v.Add("show_alert", strconv.FormatBool(cbq.ShowAlert))
	v.Add("url", cbq.Url)
	v.Add("cache_time", strconv.Itoa(cbq.CacheTime))

	return cbq.bot.boolSender("answerCallbackQuery", v)
}
