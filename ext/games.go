package ext

import (
	"encoding/json"
	"strconv"
	"net/url"
	"github.com/pkg/errors"
	"github.com/PaulSonOfLars/gotgbot/types"
)

// TODO: reply version
// TODO: no notif version
func (b Bot) SendGame(chatId int, gameShortName string) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("game_short_name", gameShortName)

	r, err := Get(b, "sendGame", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to execute sendGame request")
	}
	if !r.Ok {
		return nil, errors.Wrapf(err, "invalid sendGame request")
	}

	return b.ParseMessage(r.Result), nil
}


// TODO Check return - message or bool?
// TODO: Force version
func (b Bot) SetGameScore(userId int, score int, chatId int, messageId int) (bool, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(userId))
	v.Add("score", strconv.Itoa(score))
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))

	r, err := Get(b, "setGameScore", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to execute setGameScore request")
	}
	if !r.Ok {
		return false, errors.Wrapf(err, "invalid setGameScore request")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)
	return bb, nil
}

// TODO Check return - message or bool?
// TODO: Force version
func (b Bot) SetGameScoreInline(userId int, score int, inlineMessageId string) (bool, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(userId))
	v.Add("score", strconv.Itoa(score))
	v.Add("inline_message_id", inlineMessageId)
	r, err := Get(b, "setGameScore", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to execute setGameScore request")
	}
	if !r.Ok {
		return false, errors.Wrapf(err, "invalid setGameScore request")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)
	return bb, nil
}

func (b Bot) GetGameHighScores(userId int, chatId int, messageId int) ([]types.GameHighScore, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(userId))
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("message_id", strconv.Itoa(messageId))

	r, err := Get(b, "getGameHighScores", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to execute getGameHighScores request")
	}
	if !r.Ok {
		return nil, errors.Wrapf(err, "invalid getGameHighScores request")
	}

	var ghs []types.GameHighScore
	json.Unmarshal(r.Result, &ghs)
	return ghs, nil
}

func (b Bot) GetGameHighScoresInline(userId int, inlineMessageId string) ([]types.GameHighScore, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(userId))
	v.Add("inline_message_id", inlineMessageId)

	r, err := Get(b, "getGameHighScores", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to execute getGameHighScores request")
	}
	if !r.Ok {
		return nil, errors.Wrapf(err, "invalid getGameHighScores request")
	}

	var ghs []types.GameHighScore
	json.Unmarshal(r.Result, &ghs)

	return ghs, nil
}