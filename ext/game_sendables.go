package ext

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type sendableGame struct {
	bot                 Bot `json:"-"`
	ChatId              int
	GameShortName       string
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         InlineKeyboardMarkup
}

func (b Bot) NewSendableGame(chatId int, gameShortName string) *sendableGame {
	return &sendableGame{bot: b, ChatId: chatId, GameShortName: gameShortName}
}

func (g *sendableGame) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(g.ChatId))
	v.Add("game_short_name", g.GameShortName)

	r, err := g.bot.Get("sendGame", v)
	if err != nil {
		return nil, err
	}

	return g.bot.ParseMessage(r)
}

type sendableSetGameScore struct {
	bot                Bot `json:"-"`
	UserId             int
	Score              int
	Force              bool
	DisableEditMessage bool
	ChatId             int
	MessageId          int
	InlineMessageId    string
}

func (b Bot) NewSendableSetGameScore(userId int, score int, chatId int, messageId int) *sendableSetGameScore {
	return &sendableSetGameScore{bot: b, UserId: userId, Score: score, ChatId: chatId, MessageId: messageId}
}

func (b Bot) NewSendableSetGameScoreInline(userId int, score int, inlineMessageId string) *sendableSetGameScore {
	return &sendableSetGameScore{bot: b, UserId: userId, Score: score, InlineMessageId: inlineMessageId}
}

func (sgs *sendableSetGameScore) Send() (bool, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(sgs.UserId))
	v.Add("score", strconv.Itoa(sgs.Score))
	v.Add("force", strconv.FormatBool(sgs.Force))
	v.Add("disable_edit_message", strconv.FormatBool(sgs.DisableEditMessage))
	v.Add("chat_id", strconv.Itoa(sgs.ChatId))
	v.Add("message_id", strconv.Itoa(sgs.MessageId))
	v.Add("inline_message_id", sgs.InlineMessageId)

	r, err := sgs.bot.Get("setGameScore", v)
	if err != nil {
		return false, err
	}
	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

type sendableGetGameHighScores struct {
	bot             Bot `json:"-"`
	UserId          int
	ChatId          int
	MessageId       int
	InlineMessageId string
}

func (b Bot) NewSendableGetGameHighScore(userId int, chatId int, messageId int) *sendableGetGameHighScores {
	return &sendableGetGameHighScores{bot: b, UserId: userId, ChatId: chatId, MessageId: messageId}
}

func (b Bot) NewSendableGetGameHighScoreInline(userId int, inlineMessageId string) *sendableGetGameHighScores {
	return &sendableGetGameHighScores{bot: b, UserId: userId, InlineMessageId: inlineMessageId}
}

func (gghs *sendableGetGameHighScores) Send() ([]GameHighScore, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(gghs.UserId))
	v.Add("chat_id", strconv.Itoa(gghs.ChatId))
	v.Add("message_id", strconv.Itoa(gghs.MessageId))
	v.Add("inline_message_id", gghs.InlineMessageId)

	r, err := gghs.bot.Get("getGameHighScores", v)
	if err != nil {
		return nil, err
	}
	var ghs []GameHighScore
	return ghs, json.Unmarshal(r, &ghs)
}
