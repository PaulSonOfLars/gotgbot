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

}

type PhotoSize struct {
	File_id   string
	Width     int
	Height    int
	File_size int

}

type Sticker struct {

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

}

type SuccessfulPayment struct {

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