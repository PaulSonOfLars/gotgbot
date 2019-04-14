package ext

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf16"
)

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Url    string `json:"url"`
	User   *User  `json:"user"`
}

type ParsedMessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Url    string `json:"url"`
	User   *User  `json:"user"`
	Text   string `json:"text"`
}

type Audio struct {
	FileId    string `json:"file_id"`
	Duration  int    `json:"duration"`
	Performer string `json:"performer"`
	Title     string `json:"title"`
	MimeType  string `json:"mime_type"`
	FileSize  int    `json:"file_size"`
}

type Document struct {
	FileId   string    `json:"file_id"`
	Thumb    PhotoSize `json:"thumb"`
	FileName string    `json:"file_name"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

type PhotoSize struct {
	FileId   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

type Video struct {
	FileId   string    `json:"file_id"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Duration int       `json:"duration"`
	Thumb    PhotoSize `json:"thumb"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

type Voice struct {
	FileId   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

type VideoNote struct {
	FileId   string    `json:"file_id"`
	Length   int       `json:"length"`
	Duration int       `json:"duration"`
	Thumb    PhotoSize `json:"thumb"`
	FileSize int       `json:"file_size"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserId      int    `json:"user_id"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Venue struct {
	Location     Location `json:"location"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	FoursquareId string   `json:"foursquare_id"`
}

type PreCheckoutQuery struct {
	Id               string    `json:"id"`
	From             *User     `json:"from"`
	Currency         string    `json:"currency"`
	TotalAmount      int       `json:"total_amount"`
	InvoicePayload   string    `json:"invoice_payload"`
	ShippingOptionId string    `json:"shipping_option_id"`
	OrderInfo        OrderInfo `json:"order_info"`
}

type PollOption struct {
	Text       string
	VoterCount int
}

type Poll struct {
	Id       string
	Question string
	Options  []PollOption
	IsClosed bool
}

type Message struct {
	Bot                   Bot
	MessageId             int                `json:"message_id"`
	From                  *User              `json:"from"`
	Date                  int                `json:"date"`
	Chat                  *Chat              `json:"chat"`
	ForwardFrom           *User              `json:"forward_from"`
	ForwardFromChat       *Chat              `json:"forward_from_chat"`
	ForwardFromMessageId  int                `json:"forward_from_message_id"`
	ForwardSignature      string             `json:"forward_signature"`
	ForwardDate           int                `json:"forward_date"`
	ReplyToMessage        *Message           `json:"reply_to_message"`
	EditDate              int                `json:"edit_date"`
	MediaGroupId          string             `json:"media_group_id"`
	AuthorSignature       string             `json:"author_signature"`
	Text                  string             `json:"text"`
	Entities              []MessageEntity    `json:"entities"`
	CaptionEntities       []MessageEntity    `json:"caption_entities"`
	Audio                 *Audio             `json:"audio"`
	Document              *Document          `json:"document"`
	Animation             *Animation         `json:"animation"`
	Game                  *Game              `json:"game"`
	Photo                 []PhotoSize        `json:"photo"`
	Sticker               *Sticker           `json:"sticker"`
	Video                 *Video             `json:"video"`
	Voice                 *Voice             `json:"voice"`
	VideoNote             *VideoNote         `json:"video_note"`
	Caption               string             `json:"caption"`
	Contact               *Contact           `json:"contact"`
	Location              *Location          `json:"location"`
	Venue                 *Venue             `json:"venue"`
	NewChatMembers        []User             `json:"new_chat_members"`
	LeftChatMember        *User              `json:"left_chat_member"`
	NewChatTitle          string             `json:"new_chat_title"`
	NewChatPhoto          []PhotoSize        `json:"new_chat_photo"`
	DeleteChatPhoto       bool               `json:"delete_chat_photo"`
	GroupChatCreated      bool               `json:"group_chat_created"`
	SupergroupChatCreated bool               `json:"supergroup_chat_created"`
	ChannelChatCreated    bool               `json:"channel_chat_created"`
	MigrateToChatId       int                `json:"migrate_to_chat_id"`
	MigrateFromChatId     int                `json:"migrate_from_chat_id"`
	PinnedMessage         *Message           `json:"pinned_message"`
	Invoice               *Invoice           `json:"invoice"`
	SuccessfulPayment     *SuccessfulPayment `json:"successful_payment"`
	ConnectedWebsite      string             `json:"connected_website"`
	PassportData          PassportData       `json:"passport_data"`
	Poll                  Poll               `json:"poll"`

	// internals
	utf16Text       []uint16
	utf16Caption    []uint16
	originalText    string
	originalCaption string
}

func (b Bot) Message(chatId int, text string) Message {
	return Message{Bot: b}
}

func (b Bot) ParseMessage(message json.RawMessage) (mess *Message, err error) {
	mess = &Message{Bot: b}
	return mess, json.Unmarshal(message, mess)
}

func (m Message) ReplyText(text string) (*Message, error) {
	return m.Bot.ReplyText(m.Chat.Id, text, m.MessageId)
}

func (m Message) ReplyTextf(format string, a ...interface{}) (*Message, error) {
	return m.Bot.ReplyText(m.Chat.Id, fmt.Sprintf(format, a...), m.MessageId)
}

func (m Message) ReplyHTML(text string) (*Message, error) {
	return m.Bot.ReplyHTML(m.Chat.Id, text, m.MessageId)
}

func (m Message) ReplyHTMLf(format string, a ...interface{}) (*Message, error) {
	return m.Bot.ReplyHTML(m.Chat.Id, fmt.Sprintf(format, a...), m.MessageId)
}

func (m Message) ReplyMarkdown(text string) (*Message, error) {
	return m.Bot.ReplyMarkdown(m.Chat.Id, text, m.MessageId)
}

func (m Message) ReplyMarkdownf(format string, a ...interface{}) (*Message, error) {
	return m.Bot.ReplyMarkdown(m.Chat.Id, fmt.Sprintf(format, a...), m.MessageId)
}

func (m Message) EditText(text string) (*Message, error) {
	return m.Bot.EditMessageText(m.Chat.Id, m.MessageId, text)
}

func (m Message) EditTextf(format string, a ...interface{}) (*Message, error) {
	return m.Bot.EditMessageText(m.Chat.Id, m.MessageId, fmt.Sprintf(format, a...))
}

func (m Message) EditHTML(text string) (*Message, error) {
	return m.Bot.EditMessageHTML(m.Chat.Id, m.MessageId, text)
}

func (m Message) EditHTMLf(format string, a ...interface{}) (*Message, error) {
	return m.Bot.EditMessageHTML(m.Chat.Id, m.MessageId, fmt.Sprintf(format, a...))
}

func (m Message) EditMarkdown(text string) (*Message, error) {
	return m.Bot.EditMessageMarkdown(m.Chat.Id, m.MessageId, text)
}

func (m Message) EditMarkdownf(format string, a ...interface{}) (*Message, error) {
	return m.Bot.EditMessageMarkdown(m.Chat.Id, m.MessageId, fmt.Sprintf(format, a...))
}

func (m Message) ReplyAudioStr(audio string) (*Message, error) {
	return m.Bot.ReplyAudioStr(m.Chat.Id, audio, m.MessageId)
}

func (m Message) ReplyDocumentStr(document string) (*Message, error) {
	return m.Bot.ReplyDocumentStr(m.Chat.Id, document, m.MessageId)
}

func (m Message) ReplyLocation(latitude float64, longitude float64) (*Message, error) {
	return m.Bot.ReplyLocation(m.Chat.Id, latitude, longitude, m.MessageId)
}

func (m Message) ReplyPhotoStr(photo string) (*Message, error) {
	return m.Bot.ReplyPhotoStr(m.Chat.Id, photo, m.MessageId)
}

func (m Message) ReplyStickerStr(sticker string) (*Message, error) {
	return m.Bot.ReplyStickerStr(m.Chat.Id, sticker, m.MessageId)
}

func (m Message) ReplyVenue(latitude float64, longitude float64, title string, address string) (*Message, error) {
	return m.Bot.ReplyVenue(m.Chat.Id, latitude, longitude, title, address, m.MessageId)
}

func (m Message) ReplyVideoStr(video string) (*Message, error) {
	return m.Bot.ReplyVideoStr(m.Chat.Id, video, m.MessageId)
}

func (m Message) ReplyVideoNoteStr(videoNote string) (*Message, error) {
	return m.Bot.ReplyVideoNoteStr(m.Chat.Id, videoNote, m.MessageId)
}

func (m Message) ReplyVoiceStr(voice string) (*Message, error) {
	return m.Bot.ReplyVoiceStr(m.Chat.Id, voice, m.MessageId)
}

func (m Message) Delete() (bool, error) {
	return m.Bot.DeleteMessage(m.Chat.Id, m.MessageId)
}

func (m Message) Forward(chatId int) (*Message, error) {
	return m.Bot.ForwardMessage(chatId, m.Chat.Id, m.MessageId)
}

func (m *Message) ParseEntities() (out []ParsedMessageEntity) {
	for _, ent := range m.Entities {
		out = append(out, m.ParseEntity(ent))
	}
	return out
}

func (m *Message) ParseCaptionEntities() (out []ParsedMessageEntity) {
	for _, ent := range m.CaptionEntities {
		out = append(out, m.ParseCaptionEntity(ent))
	}
	return out
}

func (m *Message) ParseEntityTypes(accepted map[string]struct{}) (out []ParsedMessageEntity) {
	for _, ent := range m.Entities {
		if _, ok := accepted[ent.Type]; ok {
			out = append(out, m.ParseEntity(ent))
		}
	}
	return out
}

func (m *Message) ParseCaptionEntityTypes(accepted map[string]struct{}) (out []ParsedMessageEntity) {
	for _, ent := range m.CaptionEntities {
		if _, ok := accepted[ent.Type]; ok {
			out = append(out, m.ParseCaptionEntity(ent))
		}
	}
	return out
}

func (m *Message) ParseEntity(entity MessageEntity) ParsedMessageEntity {
	if m.utf16Text == nil {
		m.utf16Text = utf16.Encode([]rune(m.Text))
	}
	text := string(utf16.Decode(m.utf16Text[entity.Offset : entity.Offset+entity.Length]))
	if entity.User != nil {
		entity.User.Bot = m.Bot
	}
	if entity.Type == "url" {
		entity.Url = text
	}
	return ParsedMessageEntity{
		Type:   entity.Type,
		Offset: len(string(utf16.Decode(m.utf16Text[:entity.Offset]))),
		Length: len(text),
		Url:    entity.Url,
		User:   entity.User,
		Text:   text,
	}
}

func (m *Message) ParseCaptionEntity(entity MessageEntity) ParsedMessageEntity {
	if m.utf16Caption == nil {
		m.utf16Caption = utf16.Encode([]rune(m.Caption))
	}
	text := string(utf16.Decode(m.utf16Caption[entity.Offset : entity.Offset+entity.Length]))
	if entity.User != nil {
		entity.User.Bot = m.Bot
	}
	if entity.Type == "url" {
		entity.Url = text
	}
	return ParsedMessageEntity{
		Type:   entity.Type,
		Offset: len(string(utf16.Decode(m.utf16Caption[:entity.Offset]))),
		Length: len(text),
		Url:    entity.Url,
		User:   entity.User,
		Text:   text,
	}
}

var mdMap = map[string]string{
	"bold":   "*",
	"italic": "_",
	"code":   "`",
}

func (m *Message) OriginalText() string {
	if m.originalText != "" {
		return m.originalText
	}
	if m.utf16Text == nil {
		m.utf16Text = utf16.Encode([]rune(m.Text))
	}

	m.originalText = getOrigMsgTxt(m.utf16Text, m.Entities)
	return m.originalText
}

func (m *Message) OriginalCaption() string {
	if m.originalCaption != "" {
		return m.originalCaption
	}
	if m.utf16Caption == nil {
		m.utf16Caption = utf16.Encode([]rune(m.Caption))
	}

	m.originalCaption = getOrigMsgTxt(m.utf16Caption, m.CaptionEntities)
	return m.originalCaption
}

func getOrigMsgTxt(utf16Data []uint16, ents []MessageEntity) (out string) {
	prev := 0
	for _, ent := range ents {
		newPrev := ent.Offset + ent.Length
		switch ent.Type {
		case "bold", "italic", "code":
			out += string(utf16.Decode(utf16Data[prev:ent.Offset])) + mdMap[ent.Type] + string(utf16.Decode(utf16Data[ent.Offset:newPrev])) + mdMap[ent.Type]
			prev = newPrev
		case "text_mention":
			out += string(utf16.Decode(utf16Data[prev:ent.Offset])) + "[" + string(utf16.Decode(utf16Data[ent.Offset:ent.Offset+ent.Length])) + "](tg://user?id=" + strconv.Itoa(ent.User.Id) + ")"
			prev = newPrev
		case "text_link":
			out += string(utf16.Decode(utf16Data[prev:ent.Offset])) + "[" + string(utf16.Decode(utf16Data[ent.Offset:ent.Offset+ent.Length])) + "](" + ent.Url + ")"
			prev = newPrev
		}
	}
	out += string(utf16.Decode(utf16Data[prev:]))
	return out
}
