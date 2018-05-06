package ext

import (
	"net/http"
	"encoding/json"
	"net/url"
	"log"
)

var apiUrl = "https://api.telegram.org/bot"

var client = &http.Client{}

type Response struct {
	Ok          bool
	Result      json.RawMessage
	ErrorCode  int `json:"error_code"`
	Description string
	Parameters  json.RawMessage
}

func Get(bot Bot, method string, params url.Values) Response {
	req, err := http.NewRequest("GET", apiUrl+bot.Token+"/"+method, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var r Response
	json.NewDecoder(resp.Body).Decode(&r)
	return r
}
