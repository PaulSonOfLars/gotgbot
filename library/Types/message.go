package Types

type MessageEntity struct {
	Type   string
	Offset int
	Length int
	Url    string
	User   User
}

type Audio struct {
	File_id   string
	Duration  int
	Performer string
	Title     string
	Mime_type string
	File_size int
}

type Document struct {
	File_id   string
	Thumb     PhotoSize
	File_name string
	Mime_type string
	File_size int
}

type Game struct {
	Title         string
	Description   string
	Photo         []PhotoSize
	Text          string
	Text_entities []MessageEntity
	Animation     Animation
}

type Animation struct {
	File_id   string
	Thumb     PhotoSize
	File_name string
	Mime_type string
	File_size int
}

type GameHighScore struct {
	Position int
	User     User
	Score    int
}

type PhotoSize struct {
	File_id   string
	Width     int
	Height    int
	File_size int

}

type Sticker struct {
	File_id       string
	Width         int
	Height        int
	Thumb         PhotoSize
	Emoji         string
	Set_name      string
	Mask_position MaskPosition
	File_size     int
}

type StickerSet struct {
	Name           string
	Title          string
	Contains_masks bool
	Stickers       []Sticker
}

type MaskPosition struct {
	Point   string
	X_shift float64
	Y_shift float64
	Scale   float64
}

type Video struct {
	File_id   string
	Width     int
	Height    int
	Duration  int
	Thumb     PhotoSize
	Mime_type string
	File_size int
}

type Voice struct {
	File_id   string
	Duration  int
	Mime_type string
	File_size int
}

type VideoNote struct {
	File_id   string
	Length    int
	Duration  int
	Thumb     PhotoSize
	File_size int
}

type Contact struct {
	Phone_number string
	First_name   string
	Last_name    string
	User_id      int
}

type Location struct {
	Longitude float64
	Latitude  float64

}

type Venue struct {
	Location      Location
	Title         string
	Address       string
	Foursquare_id string

}

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

type PreCheckoutQuery struct {
	Id                 string
	From               User
	Currency           string
	Total_amount       int
	Invoice_payload    string
	Shipping_option_id string
	Order_info         OrderInfo
}

type Message struct {
	Message_id              int
	From                    *User
	Date                    int
	Chat                    *Chat
	Forward_from            *User
	Forward_from_chat       *Chat
	Forward_from_message_id int
	Forward_signature       string
	Forward_date            int
	Reply_to_message        *Message
	Edit_date               int
	Author_signature        string
	Text                    string
	Entities                []MessageEntity
	Audio                   *Audio
	Document                *Document
	Game                    *Game
	Photo                   *PhotoSize
	Sticker                 *Sticker
	Video                   *Video
	Voice                   *Voice
	Video_note              *VideoNote
	New_chat_members        []User
	Caption                 string
	Contact                 *Contact
	Location                *Location
	Venue                   *Venue
	New_chat_member         *User
	Left_chat_member        *User
	New_chat_title          string
	New_chat_phot           []PhotoSize
	Delete_chat_photo       bool
	Group_chat_created      bool
	Supergroup_chat_created bool
	Channel_chat_created    bool
	Migrate_to_chat_id      int
	Migrate_from_chat_id    int
	Pinned_message          *Message
	Invoice                 *Invoice
	Successful_payment      *SuccessfulPayment
}