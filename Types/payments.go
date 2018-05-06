package Types


type Invoice struct {
	Title           string
	Description     string
	Start_parameter string
	Currency        string
	Total_amount    int
}

type LabeledPrice struct {
	Label  string
	Amount int
}

type ShippingAddress struct {
	Country_code string
	State        string
	City         string
	Street_line1 string
	Street_line2 string
	Post_code    string
}

type OrderInfo struct {
	Name             string
	Phone_number     string
	Email            string
	Shipping_address ShippingAddress
}

type ShippingOption struct {
	Id     string
	Title  string
	Prices []LabeledPrice
}

type SuccessfulPayment struct {
	Currency                   string
	Total_amount               int
	Invoice_payload            string
	Shipping_option_id         string
	Order_info                 OrderInfo
	Telegram_payment_charge_id string
	Provider_payment_charge_id string
}

type ShippingQuery struct {
	Id               string
	From             User
	Invoice_payload  string
	Shipping_address ShippingAddress
}
