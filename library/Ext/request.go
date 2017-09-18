package Ext

import (
	"net/http"
	"encoding/json"
	"bot/library/Types"
	"net/url"
	"log"
)

var api_url = "https://api.telegram.org/bot"

var client = &http.Client{}

type Response struct {
	Ok     bool
	Result json.RawMessage
}

type Result struct {
	Message_id int
	From       *Types.User
	Chat       *Types.Chat
	Date       int
	Text       string
}

func Get(bot Bot, method string, params url.Values) Response {
	req, err := http.NewRequest("GET", api_url + bot.Token + "/" + method, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.URL.RawQuery = params.Encode()

	resp , err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var r Response
	json.NewDecoder(resp.Body).Decode(&r)
	return r
}