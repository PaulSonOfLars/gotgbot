package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type sendableInvoice struct {
	bot                       Bot
	ChatId                    int
	Title                     string
	Description               string
	Payload                   string
	ProviderToken             string
	StartParameter            string
	Currency                  string
	Prices                    []LabeledPrice
	ProviderData              string
	PhotoUrl                  string
	PhotoSize                 int
	PhotoWidth                int
	PhotoHeight               int
	NeedName                  bool
	NeedPhoneNumber           bool
	NeedEmail                 bool
	NeedShippingAddress       bool
	SendPhoneNumberToProvider bool
	SendEmailToProvider       bool
	IsFlexible                bool
	DisableNotification       bool
	ReplyToMessageId          int
	ReplyMarkup               ReplyMarkup
}

func (b Bot) NewSendableInvoice(chatId int, title string, description string, payload string, providerToken string, startParameter string, currency string, prices []LabeledPrice) *sendableInvoice {
	return &sendableInvoice{
		bot:            b,
		ChatId:         chatId,
		Title:          title,
		Description:    description,
		Payload:        payload,
		ProviderToken:  providerToken,
		StartParameter: startParameter,
		Currency:       currency,
		Prices:         prices,
	}
}

func (i *sendableInvoice) Send() (*Message, error) {
	pricesStr, err := json.Marshal(i.Prices)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal invoice prices")
	}
	var replyMarkup []byte
	if i.ReplyMarkup != nil {
		replyMarkup, err = i.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(i.ChatId))
	v.Add("title", i.Title)
	v.Add("description", i.Description)
	v.Add("payload", i.Payload)
	v.Add("provider_token", i.ProviderToken)
	v.Add("start_parameter", i.StartParameter)
	v.Add("currency", i.Currency)
	v.Add("prices", string(pricesStr))
	v.Add("provider_data", i.ProviderData)
	v.Add("photo_url", i.PhotoUrl)
	v.Add("photo_size", strconv.Itoa(i.PhotoSize))
	v.Add("photo_width", strconv.Itoa(i.PhotoWidth))
	v.Add("photo_height", strconv.Itoa(i.PhotoHeight))
	v.Add("need_name", strconv.FormatBool(i.NeedName))
	v.Add("need_phone_number", strconv.FormatBool(i.NeedPhoneNumber))
	v.Add("need_email", strconv.FormatBool(i.NeedEmail))
	v.Add("need_shipping_address", strconv.FormatBool(i.NeedShippingAddress))
	v.Add("send_phone_number_to_provider", strconv.FormatBool(i.SendPhoneNumberToProvider))
	v.Add("send_email_to_provider", strconv.FormatBool(i.SendEmailToProvider))
	v.Add("is_flexible", strconv.FormatBool(i.IsFlexible))
	v.Add("disable_notification", strconv.FormatBool(i.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(i.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := Get(i.bot, "sendInvoice", v)
	if err != nil {
		return nil, errors.New(r.Description)
	}

	return i.bot.ParseMessage(r.Result)
}

type sendableAnswerShippingQuery struct {
	bot             Bot
	ShippingQueryId string
	Ok              bool
	ShippingOptions []ShippingOption
	ErrorMessage    string
}

func (b Bot) NewSendableAnswerShippingQuery(shippingQueryId string, ok bool) *sendableAnswerShippingQuery {
	return &sendableAnswerShippingQuery{ShippingQueryId: shippingQueryId, Ok: ok}
}

func (asq *sendableAnswerShippingQuery) Send() (bool, error) {
	shippingOptions, err := json.Marshal(asq.ShippingOptions)
	if err != nil {
		return false, errors.Wrapf(err, "could not marshal shipping query shipping options")
	}

	v := url.Values{}
	v.Add("shipping_query_id", asq.ShippingQueryId)
	v.Add("ok", strconv.FormatBool(asq.Ok))
	v.Add("shipping_options", string(shippingOptions))
	v.Add("error_message", asq.ErrorMessage)

	r, err := Get(asq.bot, "answerShippingQuery", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to answerShippingQuery")
	}

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

type sendableAnswerPreCheckoutQuery struct {
	bot             Bot
	ShippingQueryId string
	Ok              bool
	ShippingOptions []ShippingOption
	ErrorMessage    string
}

func (b Bot) NewSendableAnswerPreCheckoutQuery(shippingQueryId string, ok bool) *sendableAnswerPreCheckoutQuery {
	return &sendableAnswerPreCheckoutQuery{ShippingQueryId: shippingQueryId, Ok: ok}
}

func (apcq *sendableAnswerPreCheckoutQuery) Send() (bool, error) {
	shippingOptions, err := json.Marshal(apcq.ShippingOptions)
	if err != nil {
		return false, errors.Wrapf(err, "could not marshal pre checkout query shipping options")
	}

	v := url.Values{}
	v.Add("shipping_query_id", apcq.ShippingQueryId)
	v.Add("ok", strconv.FormatBool(apcq.Ok))
	v.Add("shipping_options", string(shippingOptions))
	v.Add("error_message", apcq.ErrorMessage)

	r, err := Get(apcq.bot, "answerPreCheckoutQuery", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to answerPreCheckoutQuery")
	}
	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}
