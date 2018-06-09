package ext


type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
}

type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

type OrderInfo struct {
	Name            string          `json:"name"`
	PhoneNumber     string          `json:"phone_number"`
	Email           string          `json:"email"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

type ShippingOption struct {
	Id     string         `json:"id"`
	Title  string         `json:"title"`
	Prices []LabeledPrice `json:"prices"`
}

type SuccessfulPayment struct {
	Currency                string    `json:"currency"`
	TotalAmount             int       `json:"total_amount"`
	InvoicePayload          string    `json:"invoice_payload"`
	ShippingOptionId        string    `json:"shipping_option_id"`
	OrderInfo               OrderInfo `json:"order_info"`
	TelegramPaymentChargeId string    `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeId string    `json:"provider_payment_charge_id"`
}

type ShippingQuery struct {
	Id              string          `json:"id"`
	From            *User            `json:"from"`
	InvoicePayload  string          `json:"invoice_payload"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}


// TODO: all the optionals here. Best option is probably to use a builder.
func (b Bot) SendInvoice(chatId int, title string, description string, payload string,
	providerToken string, startParameter string, currency string,
	prices []LabeledPrice) (*Message, error) {
	return b.NewSendableInvoice(chatId, title, description, payload, providerToken, startParameter, currency, prices).Send()
}

func (b Bot) AnswerShippingQuery(shippingQueryId string, ok bool) (bool, error) {
	return b.NewSendableAnswerShippingQuery(shippingQueryId, ok).Send()
}

func (b Bot) AnswerPreCheckoutQuery(preCheckoutQueryId string, ok bool) (bool, error) {
	return b.NewSendableAnswerPreCheckoutQuery(preCheckoutQueryId, ok).Send()
}
