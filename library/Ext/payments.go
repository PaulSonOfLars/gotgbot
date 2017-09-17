package Ext

import (
	"bot/library/Types"
	"strconv"
	"log"
	"encoding/json"
)
// TODO: all the optionals here. Best option is probably to use a builder.
func (b Bot) SendInvoice(chat_id int, title string, description string, payload string,
						provider_token string, start_parameter string, currency string,
						prices []Types.LabeledPrice) Types.Message {
	prices_str, err := json.Marshal(prices)
	if err != nil {
		log.Println("Err in send_invoice")
		log.Fatal(err)
	}
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["title"] = title
	m["description"] = description
	m["payload"] = payload
	m["provider_token"] = provider_token
	m["start_parameter"] = start_parameter
	m["currency"] = currency
	m["prices"] = string(prices_str)

	r := Get(b, "sendInvoice", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var mess Types.Message
	json.Unmarshal(r.Result, &mess)

	return mess
}

// TODO: shipping options
// TODO: err_msg
func (b Bot) AnswerShippingQuery(shipping_query_id string, ok bool) bool {
	m := make(map[string]string)
	m["shipping_query_id"] = shipping_query_id
	m["ok"] = strconv.FormatBool(ok)

	r := Get(b, "answerShippingQuery", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) AnswerPreCheckoutQuery(pre_checkout_query_id string, ok bool) bool {
	m := make(map[string]string)
	m["pre_checkout_query_id"] = pre_checkout_query_id
	m["ok"] = strconv.FormatBool(ok)

	r := Get(b, "answerPreCheckoutQuery", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}