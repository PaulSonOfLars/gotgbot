package ext

import (
	"net/url"
	"strconv"
	"encoding/json"
	"gotgbot/types"
	"github.com/pkg/errors"
)

type newMessage struct {
	bot              Bot
	chatId           int
	replyToMessageId int
	text             string
	parseMode        string
}

const (
	Markdown = "Markdown"
	Html = "HTML"
)


func (b Bot) NewSendableMessage(chatId int, text string) *newMessage {
	return &newMessage{bot: b, chatId: chatId, text: text}
}

func (msg *newMessage) SetParseMode(parseMode string) {
	msg.parseMode = parseMode
}

func (msg *newMessage) ReplyToMsg(replyTo *types.Message) {
	msg.replyToMessageId = replyTo.Message_id
}

func (msg *newMessage) ReplyToMsgId(replyTo int) {
	msg.replyToMessageId = replyTo
}

func (msg *newMessage) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(msg.chatId))
	v.Add("text", msg.text)
	v.Add("reply_to_message_id", strconv.Itoa(msg.replyToMessageId))
	v.Add("parse_mode", msg.parseMode)

	r := Get(msg.bot, "sendMessage", v)

	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newMsg := &Message{}
	newMsg.bot = msg.bot
	json.Unmarshal(r.Result, newMsg)
	return newMsg, nil
}