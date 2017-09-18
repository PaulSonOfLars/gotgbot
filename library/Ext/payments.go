package Ext

import (
	"bot/library/Types"
	"strconv"
	"log"
	"encoding/json"
	"net/url"
)
// TODO: all the optionals here. Best option is probably to use a builder.
func (b Bot) SendInvoice(chat_id int, title string, description string, payload string,
						provider_token string, start_parameter string, currency string,
						prices []Types.LabeledPrice) Message {
	prices_str, err := json.Marshal(prices)
	if err != nil {
		log.Println("Err in send_invoice")
		log.Fatal(err)
	}
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("title", title)
	v.Add("description", description)
	v.Add("payload", payload)
	v.Add("provider_token", provider_token)
	v.Add("start_parameter", start_parameter)
	v.Add("currency", currency)
	v.Add("prices", string(prices_str))

	r := Get(b, "sendInvoice", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}

// TODO: shipping options
// TODO: err_msg
func (b Bot) AnswerShippingQuery(shipping_query_id string, ok bool) bool {
	v := url.Values{}
	v.Add("shipping_query_id", shipping_query_id)
	v.Add("ok", strconv.FormatBool(ok))

	r := Get(b, "answerShippingQuery", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) AnswerPreCheckoutQuery(pre_checkout_query_id string, ok bool) bool {
	v := url.Values{}
	v.Add("pre_checkout_query_id", pre_checkout_query_id)
	v.Add("ok", strconv.FormatBool(ok))

	r := Get(b, "answerPreCheckoutQuery", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}