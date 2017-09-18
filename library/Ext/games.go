package Ext

import (
	"bot/library/Types"
	"log"
	"encoding/json"
	"strconv"
	"net/url"
)

// TODO: reply version
// TODO: no notif version
func (b Bot) SendGame(chat_id int, game_short_name string) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("game_short_name", game_short_name)

	r := Get(b, "sendGame", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for sendGame was not OK")
	}

	return b.ParseMessage(r.Result)
}


// TODO Check return - message or bool?
// TODO: Force version
// TODO: possible error
func (b Bot) SetGameScore(user_id int, score int, chat_id int, message_id int) bool {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(user_id))
	v.Add("score", strconv.Itoa(score))
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("message_id", strconv.Itoa(message_id))

	r := Get(b, "setGameScore", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for setGameScore was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO Check return - message or bool?
// TODO: Force version
// TODO: Possible error
func (b Bot) SetGameScoreInline(user_id int, score int, inline_message_id string) bool {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(user_id))
	v.Add("score", strconv.Itoa(score))
	v.Add("inline_message_id", inline_message_id)

	r := Get(b, "setGameScore", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for setGameScore was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) GetGameHighScores(user_id int, chat_id int, message_id int) []Types.GameHighScore {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(user_id))
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("message_id", strconv.Itoa(message_id))

	r := Get(b, "getGameHighScores", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getGameHighScores was not OK")
	}

	var ghs []Types.GameHighScore
	json.Unmarshal(r.Result, &ghs)

	return ghs
}

func (b Bot) GetGameHighScoresInline(user_id int, inline_message_id string) []Types.GameHighScore {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(user_id))
	v.Add("inline_message_id", inline_message_id)

	r := Get(b, "getGameHighScores", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getGameHighScores was not OK")
	}

	var ghs []Types.GameHighScore
	json.Unmarshal(r.Result, &ghs)

	return ghs
}