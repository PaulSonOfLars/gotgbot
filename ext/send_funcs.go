package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/parsemode"
)

func (b Bot) SendMessage(chatId int, text string) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	return newMsg.Send()
}

func (b Bot) SendMessageHTML(chatId int, text string) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	newMsg.ParseMode = parsemode.Html
	return newMsg.Send()
}

func (b Bot) SendMessageMarkdown(chatId int, text string) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	newMsg.ParseMode = parsemode.Markdown
	return newMsg.Send()
}

func (b Bot) SendMessageMarkdownV2(chatId int, text string) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	newMsg.ParseMode = parsemode.MarkdownV2
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

func (b Bot) ReplyMarkdownV2(chatId int, text string, replyToMessageId int) (*Message, error) {
	newMsg := b.NewSendableMessage(chatId, text)
	newMsg.ReplyToMessageId = replyToMessageId
	newMsg.ParseMode = parsemode.MarkdownV2
	return newMsg.Send()
}

func (b Bot) ForwardMessage(chatId int, fromChatId int, messageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("from_chat_id", strconv.Itoa(fromChatId))
	v.Add("message_id", strconv.Itoa(messageId))

	r, err := b.Get("forwardMessage", v)
	if err != nil {
		return nil, err
	}

	return b.ParseMessage(r)
}

func (b Bot) SendPhoto(chatId int, photo InputFile) (*Message, error) {
	return b.replyPhoto(chatId, photo, "", 0)
}

func (b Bot) SendPhotoCaption(chatId int, photo InputFile, caption string) (*Message, error) {
	return b.replyPhoto(chatId, photo, caption, 0)
}

func (b Bot) ReplyPhoto(chatId int, photo InputFile, replyToMessageId int) (*Message, error) {
	return b.replyPhoto(chatId, photo, "", replyToMessageId)
}

func (b Bot) ReplyPhotoCaption(chatId int, photo InputFile, caption string, replyToMessageId int) (*Message, error) {
	return b.replyPhoto(chatId, photo, caption, replyToMessageId)
}

func (b Bot) replyPhoto(chatId int, photo InputFile, caption string, replyToMessageId int) (*Message, error) {
	photoMsg := b.NewSendablePhoto(chatId, caption)
	photoMsg.InputFile = photo
	photoMsg.ReplyToMessageId = replyToMessageId
	return photoMsg.Send()
}

func (b Bot) SendAudio(chatId int, audio InputFile) (*Message, error) {
	return b.replyAudio(chatId, audio, "", 0)
}

func (b Bot) ReplyAudio(chatId int, audio InputFile, replyToMessageId int) (*Message, error) {
	return b.replyAudio(chatId, audio, "", replyToMessageId)
}

func (b Bot) replyAudio(chatId int, audio InputFile, caption string, replyToMessageId int) (*Message, error) {
	audioMsg := b.NewSendableAudio(chatId, caption)
	audioMsg.InputFile = audio
	audioMsg.ReplyToMessageId = replyToMessageId
	return audioMsg.Send()
}

func (b Bot) SendDocument(chatId int, document InputFile) (*Message, error) {
	return b.replyDocument(chatId, document, "", 0)
}

func (b Bot) SendDocumentCaption(chatId int, document InputFile, caption string) (*Message, error) {
	return b.replyDocument(chatId, document, caption, 0)
}

func (b Bot) ReplyDocument(chatId int, document InputFile, replyToMessageId int) (*Message, error) {
	return b.replyDocument(chatId, document, "", replyToMessageId)
}

func (b Bot) ReplyDocumentCaption(chatId int, document InputFile, caption string, replyToMessageId int) (*Message, error) {
	return b.replyDocument(chatId, document, caption, replyToMessageId)
}

func (b Bot) replyDocument(chatId int, document InputFile, caption string, replyToMessageId int) (*Message, error) {
	docMsg := b.NewSendableDocument(chatId, caption)
	docMsg.InputFile = document
	docMsg.ReplyToMessageId = replyToMessageId
	return docMsg.Send()
}

func (b Bot) SendVideo(chatId int, video InputFile) (*Message, error) {
	return b.replyVideo(chatId, video, "", 0)
}

func (b Bot) ReplyVideo(chatId int, video InputFile, replyToMessageId int) (*Message, error) {
	return b.replyVideo(chatId, video, "", replyToMessageId)
}

func (b Bot) replyVideo(chatId int, video InputFile, caption string, replyToMessageId int) (*Message, error) {
	videoMsg := b.NewSendableVideo(chatId, caption)
	videoMsg.InputFile = video
	videoMsg.ReplyToMessageId = replyToMessageId
	return videoMsg.Send()
}

func (b Bot) SendVoice(chatId int, voice InputFile) (*Message, error) {
	return b.replyVoice(chatId, voice, "", 0)
}

func (b Bot) ReplyVoice(chatId int, voice InputFile, replyToMessageId int) (*Message, error) {
	return b.replyVoice(chatId, voice, "", replyToMessageId)
}

func (b Bot) replyVoice(chatId int, voice InputFile, caption string, replyToMessageId int) (*Message, error) {
	voiceMsg := b.NewSendableVoice(chatId, caption)
	voiceMsg.InputFile = voice
	voiceMsg.ReplyToMessageId = replyToMessageId
	return voiceMsg.Send()
}

func (b Bot) SendVideoNote(chatId int, videoNote InputFile) (*Message, error) {
	return b.replyVideoNote(chatId, videoNote, 0)
}

func (b Bot) ReplyVideoNote(chatId int, videoNote InputFile, replyToMessageId int) (*Message, error) {
	return b.replyVideoNote(chatId, videoNote, replyToMessageId)
}

func (b Bot) replyVideoNote(chatId int, videoNote InputFile, replyToMessageId int) (*Message, error) {
	videoMsg := b.NewSendableVideoNote(chatId)
	videoMsg.InputFile = videoNote
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

func (b Bot) sendPoll(chatId int, question string, option []string, disableNotification bool, replyToMessageId int) (*Message, error) {
	pollMsg := b.NewSendablePoll(chatId, question, option)
	pollMsg.DisableNotification = disableNotification
	pollMsg.ReplyToMessageId = replyToMessageId
	return pollMsg.Send()
}

func (b Bot) SendPoll(chatId int, question string, options []string) (*Message, error) {
	return b.sendPoll(chatId, question, options, false, 0)
}

func (b Bot) ReplyPoll(chatId int, question string, options []string, replyToMessageId int) (*Message, error) {
	return b.sendPoll(chatId, question, options, false, replyToMessageId)
}

func (b Bot) stopPoll(chatId int, messageId int, replyMarkup *InlineKeyboardMarkup) (*Poll, error) {
	var replyMarkupBytes []byte
	if replyMarkup != nil {
		var err error
		replyMarkupBytes, err = replyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))
	v.Add("reply_markup", string(replyMarkupBytes))

	r, err := b.Get("forwardMessage", v)
	if err != nil {
		return nil, err
	}

	poll := &Poll{Bot: b}
	return poll, json.Unmarshal(r, poll)
}

func (b Bot) StopPoll(chatId int, messageId int) (*Poll, error) {
	return b.stopPoll(chatId, messageId, nil)
}

func (b Bot) StopPollMarkup(chatId int, messageId int, replyMarkup InlineKeyboardMarkup) (*Poll, error) {
	return b.stopPoll(chatId, messageId, &replyMarkup)
}
