package Ext

import (
	"gotgbot/Types"
	"strconv"
	"log"
	"encoding/json"
	"net/url"
)
// TODO: all the optionals here. Best option is probably to use a builder.
func (b Bot) SendInvoice(chatId int, title string, description string, payload string,
						providerToken string, startParameter string, currency string,
						prices []Types.LabeledPrice) Message {
	pricesStr, err := json.Marshal(prices)
	if err != nil {
		log.Println("Err in send_invoice")
		log.Fatal(err)
	}
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("title", title)
	v.Add("description", description)
	v.Add("payload", payload)
	v.Add("provider_token", providerToken)
	v.Add("start_parameter", startParameter)
	v.Add("currency", currency)
	v.Add("prices", string(pricesStr))

	r := Get(b, "sendInvoice", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}

// TODO: shipping options
// TODO: err_msg
func (b Bot) AnswerShippingQuery(shippingQueryId string, ok bool) bool {
	v := url.Values{}
	v.Add("shipping_query_id", shippingQueryId)
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

func (b Bot) AnswerPreCheckoutQuery(preCheckoutQueryId string, ok bool) bool {
	v := url.Values{}
	v.Add("pre_checkout_query_id", preCheckoutQueryId)
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