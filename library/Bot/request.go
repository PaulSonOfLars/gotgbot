package Bot

import (
	"encoding/json"
	"log"
	"net/http"
	"bot/library/Types"
)

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

func Get(bot Bot, method string, params map[string]string) Response {
	req, err := http.NewRequest("GET", url + bot.Token + "/" + method, nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k,v)
	}
	req.URL.RawQuery = q.Encode()

	resp , err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var r Response
	json.NewDecoder(resp.Body).Decode(&r)
	return r
}
