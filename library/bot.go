package library

import (
	"log"
	"encoding/json"
	"strconv"
)

var url = "https://api.telegram.org/bot"

type Bot struct {
	Token string

}

func (b Bot) SendMessage(msg string, chat_id int) Response {
	//req, _ := http.NewRequest("GET", url + b.Token+ "/sendMessage", nil)
	//q := req.URL.Query()
	//q.Add("chat_id", strconv.Itoa(chat_id))
	//q.Add("text", msg)
	//req.URL.RawQuery = q.Encode()
	//
	//resp , err := client.Do(req)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()
	//
	//var r Response
	//json.NewDecoder(resp.Body).Decode(&r)
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["text"] = msg

	r := Get(b, "sendMessage", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Message
	json.Unmarshal(r.Result, &mess)

	return r
}

func (b Bot) GetChat(chat_id int) Chat {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)

	r := Get(b, "getChat", m)

	var c Chat
	json.Unmarshal(r.Result, &c)

	if !r.Ok {
		log.Fatal("You done goofed, API Res for GetChat was not OK")
	}

	return c

}

func (b Bot) GetMe() User {
	m := make(map[string]string)

	r := Get(b, "getChat", m)

	var u User
	json.Unmarshal(r.Result, &u)

	if !r.Ok {
		log.Fatal("You done goofed, API Res for getMe was not OK")
	}

	return u

}
