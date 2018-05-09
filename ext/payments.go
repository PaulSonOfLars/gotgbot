package ext

import (
	"github.com/PaulSonOfLars/gotgbot/types"
	"strconv"
	"encoding/json"
	"net/url"
	"github.com/pkg/errors"
)

// TODO: all the optionals here. Best option is probably to use a builder.
func (b Bot) SendInvoice(chatId int, title string, description string, payload string,
	providerToken string, startParameter string, currency string,
	prices []types.LabeledPrice) (*Message, error) {
	pricesStr, err := json.Marshal(prices)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal invoice prices")
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

	r, err := Get(b, "sendInvoice", v)
	if err != nil {
		return nil, errors.New(r.Description)
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return b.ParseMessage(r.Result), nil

}

// TODO: shipping options
// TODO: err_msg
func (b Bot) AnswerShippingQuery(shippingQueryId string, ok bool) (bool, error) {
	v := url.Values{}
	v.Add("shipping_query_id", shippingQueryId)
	v.Add("ok", strconv.FormatBool(ok))

	r, err := Get(b, "answerShippingQuery", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to answerShippingQuery")
	}
	if !r.Ok {
		return false, errors.New("invalid answerShippingQuery")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)
	return bb, nil
}

func (b Bot) AnswerPreCheckoutQuery(preCheckoutQueryId string, ok bool) (bool, error) {
	v := url.Values{}
	v.Add("pre_checkout_query_id", preCheckoutQueryId)
	v.Add("ok", strconv.FormatBool(ok))

	r, err := Get(b, "answerPreCheckoutQuery", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to answerPreCheckoutQuery")
	}
	if !r.Ok {
		return false, errors.New("invalid answerPreCheckoutQuery")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)
	return bb, nil
}
