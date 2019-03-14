package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type sendableGame struct {
	bot                 Bot
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

	r, err := Get(g.bot, "sendGame", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to execute sendGame request")
	}
	if !r.Ok {
		return nil, errors.Wrapf(err, "invalid sendGame request")
	}

	return g.bot.ParseMessage(r.Result)
}

type sendableSetGameScore struct {
	bot                Bot
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

	r, err := Get(sgs.bot, "setGameScore", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to execute setGameScore request")
	}
	if !r.Ok {
		return false, errors.Wrapf(err, "invalid setGameScore request")
	}

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

type sendableGetGameHighScores struct {
	bot             Bot
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

	r, err := Get(gghs.bot, "getGameHighScores", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to execute getGameHighScores request")
	}
	if !r.Ok {
		return nil, errors.Wrapf(err, "invalid getGameHighScores request")
	}

	var ghs []GameHighScore
	return ghs, json.Unmarshal(r.Result, &ghs)
}
