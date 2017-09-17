package Ext

import (
	"log"
	"encoding/json"
	"strconv"
	"bot/library/Types"
)

var url = "https://api.telegram.org/bot"

type Bot struct {
	Token string

}

func (b Bot) GetMe() Types.User {
	m := make(map[string]string)

	r := Get(b, "getChat", m)

	var u Types.User
	json.Unmarshal(r.Result, &u)

	if !r.Ok {
		log.Fatal("You done goofed, API Res for getMe was not OK")
	}

	return u

}

func (b Bot) GetUserProfilePhotos(user_id int) Types.UserProfilePhotos {
	m := make(map[string]string)
	m["user_id"] = strconv.Itoa(user_id)


	r := Get(b, "getUserProfilePhotos", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var userProfilePhotos Types.UserProfilePhotos
	json.Unmarshal(r.Result, &userProfilePhotos)

	return userProfilePhotos
}


func (b Bot) GetFile(file_id string) Types.File {
	m := make(map[string]string)
	m["file_id"] = file_id

	r := Get(b, "getFile", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for getFile was not OK")
	}

	var f Types.File
	json.Unmarshal(r.Result, &f)

	return f
}

// TODO: options here
// TODO: r.OK or unmarshal??
func (b Bot) AnswerCallbackQuery(callback_query_id string) bool {
	m := make(map[string]string)
	m["callback_query_id"] = callback_query_id

	r := Get(b, "answerCallbackQuery", m)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for answerCallbackQuery was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}