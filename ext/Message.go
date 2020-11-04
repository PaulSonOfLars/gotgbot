package ext

import (
	"encoding/json"
	"fmt"
	"html"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"

	"github.com/PaulSonOfLars/gotgbot/parsemode"
)

type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	Url      string `json:"url"`
	User     *User  `json:"user"`
	Language string `json:"language"`
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
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Duration     int        `json:"duration"`
	Performer    string     `json:"performer"`
	Title        string     `json:"title"`
	FileName     string     `json:"file_name"`
	MimeType     string     `json:"mime_type"`
	FileSize     int        `json:"file_size"`
	Thumb        *PhotoSize `json:"thumb"`
}

type Document struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Thumb        *PhotoSize `json:"thumb"`
	FileName     string     `json:"file_name"`
	MimeType     string     `json:"mime_type"`
	FileSize     int        `json:"file_size"`
}

type PhotoSize struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size"`
}

type Video struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumb        *PhotoSize `json:"thumb"`
	FileName     string     `json:"file_name"`
	MimeType     string     `json:"mime_type"`
	FileSize     int        `json:"file_size"`
}

type Voice struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

type VideoNote struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Length       int        `json:"length"`
	Duration     int        `json:"duration"`
	Thumb        *PhotoSize `json:"thumb"`
	FileSize     int        `json:"file_size"`
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
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

type Poll struct {
	Bot                   Bot             `json:"-"`
	Id                    string          `json:"id"`
	Question              string          `json:"question"`
	Options               []PollOption    `json:"options"`
	TotalVoterCount       int             `json:"total_voter_count"`
	IsClosed              bool            `json:"is_closed"`
	IsAnonymous           bool            `json:"is_anonymous"`
	Type                  string          `json:"type"`
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers"`
	CorrectOptionId       int             `json:"correct_option_id"`
	Explanation           string          `json:"explanation"`
	ExplanationEntities   []MessageEntity `json:"explanation_entities"`
	OpenPeriod            int             `json:"open_period"`
	CloseDate             int             `json:"close_date"`
}

type PollAnswer struct {
	Bot       Bot    `json:"-"`
	PollId    string `json:"poll_id"`
	User      *User  `json:"user"`
	OptionIds []int  `json:"option_ids"`
}

type Dice struct {
	Bot   Bot    `json:"-"`
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

type Message struct {
	Bot                   Bot                   `json:"-"`
	MessageId             int                   `json:"message_id"`
	From                  *User                 `json:"from"`
	Date                  int                   `json:"date"`
	Chat                  *Chat                 `json:"chat"`
	ForwardFrom           *User                 `json:"forward_from"`
	ForwardFromChat       *Chat                 `json:"forward_from_chat"`
	ForwardFromMessageId  int                   `json:"forward_from_message_id"`
	ForwardSignature      string                `json:"forward_signature"`
	ForwardSenderName     string                `json:"forward_sender_name"`
	ForwardDate           int                   `json:"forward_date"`
	ReplyToMessage        *Message              `json:"reply_to_message"`
	ViaBot                *User                 `json:"via_bot"`
	EditDate              int                   `json:"edit_date"`
	MediaGroupId          string                `json:"media_group_id"`
	AuthorSignature       string                `json:"author_signature"`
	Text                  string                `json:"text"`
	Entities              []MessageEntity       `json:"entities"`
	CaptionEntities       []MessageEntity       `json:"caption_entities"`
	Audio                 *Audio                `json:"audio"`
	Document              *Document             `json:"document"`
	Animation             *Animation            `json:"animation"`
	Game                  *Game                 `json:"game"`
	Photo                 []PhotoSize           `json:"photo"`
	Sticker               *Sticker              `json:"sticker"`
	Video                 *Video                `json:"video"`
	Voice                 *Voice                `json:"voice"`
	VideoNote             *VideoNote            `json:"video_note"`
	Caption               string                `json:"caption"`
	Contact               *Contact              `json:"contact"`
	Location              *Location             `json:"location"`
	Venue                 *Venue                `json:"venue"`
	Poll                  *Poll                 `json:"poll"`
	Dice                  *Dice                 `json:"dice"`
	NewChatMembers        []User                `json:"new_chat_members"`
	LeftChatMember        *User                 `json:"left_chat_member"`
	NewChatTitle          string                `json:"new_chat_title"`
	NewChatPhoto          []PhotoSize           `json:"new_chat_photo"`
	DeleteChatPhoto       bool                  `json:"delete_chat_photo"`
	GroupChatCreated      bool                  `json:"group_chat_created"`
	SupergroupChatCreated bool                  `json:"supergroup_chat_created"`
	ChannelChatCreated    bool                  `json:"channel_chat_created"`
	MigrateToChatId       int                   `json:"migrate_to_chat_id"`
	MigrateFromChatId     int                   `json:"migrate_from_chat_id"`
	PinnedMessage         *Message              `json:"pinned_message"`
	Invoice               *Invoice              `json:"invoice"`
	SuccessfulPayment     *SuccessfulPayment    `json:"successful_payment"`
	ConnectedWebsite      string                `json:"connected_website"`
	PassportData          *PassportData         `json:"passport_data"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup"`

	// internals
	utf16Text           []uint16
	utf16Caption        []uint16
	originalTextMD      string
	originalTextMDV2    string
	originalTextHTML    string
	originalCaptionMD   string
	originalCaptionMDV2 string
	originalCaptionHTML string
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

func (m Message) EditCaption(text string) (*Message, error) {
	return m.Bot.EditMessageCaption(m.Chat.Id, m.MessageId, text)
}

func (m Message) EditCaptionf(format string, a ...interface{}) (*Message, error) {
	return m.Bot.EditMessageCaption(m.Chat.Id, m.MessageId, fmt.Sprintf(format, a...))
}

func (m Message) EditCaptionHTML(text string) (*Message, error) {
	return m.Bot.editMessageCaption(m.Chat.Id, m.MessageId, text, nil, parsemode.Html)
}

func (m Message) EditCaptionHTMLf(format string, a ...interface{}) (*Message, error) {
	return m.Bot.editMessageCaption(m.Chat.Id, m.MessageId, fmt.Sprintf(format, a...), nil, parsemode.Html)
}

func (m Message) EditCaptionMarkdown(text string) (*Message, error) {
	return m.Bot.editMessageCaption(m.Chat.Id, m.MessageId, text, nil, parsemode.Markdown)
}

func (m Message) EditCaptionMarkdownf(format string, a ...interface{}) (*Message, error) {
	return m.Bot.editMessageCaption(m.Chat.Id, m.MessageId, fmt.Sprintf(format, a...), nil, parsemode.Markdown)
}

func (m Message) ReplyAudio(audio InputFile) (*Message, error) {
	return m.Bot.ReplyAudio(m.Chat.Id, audio, m.MessageId)
}

func (m Message) ReplyDocument(document InputFile) (*Message, error) {
	return m.Bot.ReplyDocument(m.Chat.Id, document, m.MessageId)
}

func (m Message) ReplyLocation(latitude float64, longitude float64) (*Message, error) {
	return m.Bot.ReplyLocation(m.Chat.Id, latitude, longitude, m.MessageId)
}

func (m Message) ReplyPhoto(photo InputFile) (*Message, error) {
	return m.Bot.ReplyPhoto(m.Chat.Id, photo, m.MessageId)
}

func (m Message) ReplySticker(sticker InputFile) (*Message, error) {
	return m.Bot.ReplySticker(m.Chat.Id, sticker, m.MessageId)
}

func (m Message) ReplyVenue(latitude float64, longitude float64, title string, address string) (*Message, error) {
	return m.Bot.ReplyVenue(m.Chat.Id, latitude, longitude, title, address, m.MessageId)
}

func (m Message) ReplyVideo(video InputFile) (*Message, error) {
	return m.Bot.ReplyVideo(m.Chat.Id, video, m.MessageId)
}

func (m Message) ReplyVideoNote(videoNote InputFile) (*Message, error) {
	return m.Bot.ReplyVideoNote(m.Chat.Id, videoNote, m.MessageId)
}

func (m Message) ReplyVoice(voice InputFile) (*Message, error) {
	return m.Bot.ReplyVoice(m.Chat.Id, voice, m.MessageId)
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

var mdV2Map = map[string]string{
	"bold":          "*",
	"italic":        "_",
	"code":          "`",
	"pre":           "```",
	"underline":     "__",
	"strikethrough": "~",
}

var htmlMap = map[string]string{
	"bold":          "b",
	"italic":        "i",
	"code":          "code",
	"pre":           "pre",
	"underline":     "u",
	"strikethrough": "s",
}

func (m *Message) OriginalText() string {
	return m.originalMD()
}

func (m *Message) OriginalTextV2() string {
	return m.originalMDV2()
}

func (m *Message) OriginalHTML() string {
	return m.originalHTML()
}

func (m *Message) originalMD() string {
	if m.originalTextMD != "" {
		return m.originalTextMD
	}
	if m.utf16Text == nil {
		m.utf16Text = utf16.Encode([]rune(m.Text))
	}

	m.originalTextMD = getOrigMsgMD(m.utf16Text, m.Entities)
	return m.originalTextMD
}

func (m *Message) originalHTML() string {
	if m.originalTextHTML != "" {
		return m.originalTextHTML
	}
	if m.utf16Text == nil {
		m.utf16Text = utf16.Encode([]rune(m.Text))
	}

	m.originalTextHTML = getOrigMsgHTML(m.utf16Text, m.Entities)
	return m.originalTextHTML
}

func (m *Message) originalMDV2() string {
	if m.originalTextMDV2 != "" {
		return m.originalTextMDV2
	}
	if m.utf16Text == nil {
		m.utf16Text = utf16.Encode([]rune(m.Text))
	}

	m.originalTextMDV2 = getOrigMsgMDV2(m.utf16Text, m.Entities)
	return m.originalTextMDV2
}

func (m *Message) OriginalCaption() string {
	return m.originalCaptionTextMD()
}

func (m *Message) OriginalCaptionV2() string {
	return m.originalCaptionTextMDV2()
}

func (m *Message) OriginalCaptionHTML() string {
	return m.originalCaptionTextHTML()
}

func (m *Message) originalCaptionTextMD() string {
	if m.originalCaptionMD != "" {
		return m.originalCaptionMD
	}
	if m.utf16Caption == nil {
		m.utf16Caption = utf16.Encode([]rune(m.Caption))
	}

	m.originalCaptionMD = getOrigMsgMD(m.utf16Caption, m.CaptionEntities)
	return m.originalCaptionMD
}

func (m *Message) originalCaptionTextHTML() string {
	if m.originalCaptionHTML != "" {
		return m.originalCaptionHTML
	}
	if m.utf16Caption == nil {
		m.utf16Caption = utf16.Encode([]rune(m.Caption))
	}

	m.originalCaptionHTML = getOrigMsgHTML(m.utf16Caption, m.CaptionEntities)
	return m.originalCaptionHTML
}

func (m *Message) originalCaptionTextMDV2() string {
	if m.originalCaptionMDV2 != "" {
		return m.originalCaptionMDV2
	}
	if m.utf16Caption == nil {
		m.utf16Caption = utf16.Encode([]rune(m.Caption))
	}

	m.originalCaptionMDV2 = getOrigMsgMDV2(m.utf16Caption, m.CaptionEntities)
	return m.originalCaptionMDV2
}

// Does not support nesting. only look at upper entities.
func getOrigMsgMD(utf16Data []uint16, ents []MessageEntity) string {
	out := strings.Builder{}
	prev := 0
	for _, ent := range getUpperEntities(ents) {
		newPrev := ent.Offset + ent.Length
		prevText := string(utf16.Decode(utf16Data[prev:ent.Offset]))

		text := utf16.Decode(utf16Data[ent.Offset:newPrev])
		pre, cleanCntnt, post := splitEdgeWhitespace(string(text))
		cleanCntntRune := []rune(cleanCntnt)

		switch ent.Type {
		case "bold", "italic", "code":
			out.WriteString(prevText + pre + mdMap[ent.Type] + escapeContainedMDV1(cleanCntntRune, []rune(mdMap[ent.Type])) + mdMap[ent.Type] + post)
		case "text_mention":
			out.WriteString(prevText + pre + "[" + escapeContainedMDV1(cleanCntntRune, []rune("[]()")) + "](tg://user?id=" + strconv.Itoa(ent.User.Id) + ")" + post)
		case "text_link":
			out.WriteString(prevText + pre + "[" + escapeContainedMDV1(cleanCntntRune, []rune("[]()")) + "](" + ent.Url + ")" + post)
		default:
			continue
		}
		prev = newPrev
	}

	out.WriteString(string(utf16.Decode(utf16Data[prev:])))
	return out.String()
}

func getOrigMsgHTML(utf16Data []uint16, ents []MessageEntity) string {
	if len(ents) == 0 {
		return html.EscapeString(string(utf16.Decode(utf16Data)))
	}

	bd := strings.Builder{}
	prev := 0
	for _, e := range getUpperEntities(ents) {
		data, end := fillNestedHTML(utf16Data, e, prev, getChildEntities(e, ents))
		bd.WriteString(data)
		prev = end
	}

	bd.WriteString(html.EscapeString(string(utf16.Decode(utf16Data[prev:]))))
	return bd.String()
}

func getOrigMsgMDV2(utf16Data []uint16, ents []MessageEntity) string {
	if len(ents) == 0 {
		return string(utf16.Decode(utf16Data))
	}

	bd := strings.Builder{}
	prev := 0
	for _, e := range getUpperEntities(ents) {
		data, end := fillNestedMarkdownV2(utf16Data, e, prev, getChildEntities(e, ents))
		bd.WriteString(data)
		prev = end
	}

	bd.WriteString(string(utf16.Decode(utf16Data[prev:])))
	return bd.String()
}

func fillNestedHTML(data []uint16, ent MessageEntity, start int, entities []MessageEntity) (string, int) {
	entEnd := ent.Offset + ent.Length
	if len(entities) == 0 || entEnd < entities[0].Offset {
		// no nesting; just return straight away and move to next.
		return writeFinalHTML(data, ent, start, html.EscapeString(string(utf16.Decode(data[ent.Offset:entEnd])))), entEnd
	}
	subPrev := ent.Offset
	subEnd := ent.Offset
	bd := strings.Builder{}
	for _, e := range getUpperEntities(entities) {
		if e.Offset < subEnd || e == ent {
			continue
		}
		if e.Offset >= entEnd {
			break
		}

		out, end := fillNestedHTML(data, e, subPrev, getChildEntities(e, entities))
		bd.WriteString(out)
		subPrev = end
	}

	bd.WriteString(html.EscapeString(string(utf16.Decode(data[subPrev:entEnd]))))

	return writeFinalHTML(data, ent, start, bd.String()), entEnd
}

func fillNestedMarkdownV2(data []uint16, ent MessageEntity, start int, entities []MessageEntity) (string, int) {
	entEnd := ent.Offset + ent.Length
	if len(entities) == 0 || entEnd < entities[0].Offset {
		// no nesting; just return straight away and move to next.
		return writeFinalMarkdownV2(data, ent, start, string(utf16.Decode(data[ent.Offset:entEnd]))), entEnd
	}
	subPrev := ent.Offset
	subEnd := ent.Offset
	bd := strings.Builder{}
	for _, e := range getUpperEntities(entities) {
		if e.Offset < subEnd || e == ent {
			continue
		}
		if e.Offset >= entEnd {
			break
		}

		out, end := fillNestedMarkdownV2(data, e, subPrev, getChildEntities(e, entities))
		bd.WriteString(out)
		subPrev = end
	}

	bd.WriteString(string(utf16.Decode(data[subPrev:entEnd])))

	return writeFinalMarkdownV2(data, ent, start, bd.String()), entEnd
}

func splitEdgeWhitespace(text string) (pre string, cntnt string, post string) {
	bd := strings.Builder{}
	rText := []rune(text)
	for i := 0; i < len(rText) && unicode.IsSpace(rText[i]); i++ {
		bd.WriteRune(rText[i])
	}
	pre = bd.String()
	text = strings.TrimPrefix(text, pre)

	bd.Reset()
	for i := len(rText) - 1; i >= 0 && unicode.IsSpace(rText[i]); i-- {
		bd.WriteRune(rText[i])
	}
	post = bd.String()
	return pre, strings.TrimSuffix(text, post), post
}

func writeFinalHTML(data []uint16, ent MessageEntity, start int, cntnt string) string {
	prevText := html.EscapeString(string(utf16.Decode(data[start:ent.Offset])))
	switch ent.Type {
	case "bold", "italic", "code", "underline", "strikethrough":
		return prevText + "<" + htmlMap[ent.Type] + ">" + cntnt + "</" + htmlMap[ent.Type] + ">"
	case "pre":
		// <pre>text</pre>
		if ent.Language == "" {
			return prevText + "<pre>" + cntnt + "</pre>"
		}
		// <pre><code class="lang">text</code></pre>
		return prevText + `<pre><code class="` + ent.Language + `">` + cntnt + "</code></pre>"
	case "text_mention":
		return prevText + `<a href="tg://user?id=` + strconv.Itoa(ent.User.Id) + `">` + cntnt + "</a>"
	case "text_link":
		return prevText + `<a href="` + ent.Url + `">` + cntnt + "</a>"
	default:
		return prevText + cntnt
	}
}

func writeFinalMarkdownV2(data []uint16, ent MessageEntity, start int, cntnt string) string {
	prevText := string(utf16.Decode(data[start:ent.Offset]))
	pre, cleanCntnt, post := splitEdgeWhitespace(cntnt)
	switch ent.Type {
	case "bold", "italic", "code", "underline", "strikethrough", "pre":
		return prevText + pre + mdV2Map[ent.Type] + cleanCntnt + mdV2Map[ent.Type] + post
	case "text_mention":
		return prevText + pre + "[" + cleanCntnt + "](tg://user?id=" + strconv.Itoa(ent.User.Id) + ")" + post
	case "text_link":
		return prevText + pre + "[" + cleanCntnt + "](" + ent.Url + ")" + post
	default:
		return prevText + cntnt
	}
}

func getUpperEntities(ents []MessageEntity) []MessageEntity {
	prev := 0
	var uppers []MessageEntity
	for _, e := range ents {
		if e.Offset < prev {
			continue
		}
		uppers = append(uppers, e)
		prev = e.Offset + e.Length
	}
	return uppers
}

func getChildEntities(ent MessageEntity, ents []MessageEntity) []MessageEntity {
	end := ent.Offset + ent.Length
	var children []MessageEntity
	for _, e := range ents {
		if e.Offset < ent.Offset || e == ent {
			continue
		}
		if e.Offset >= end {
			break
		}
		children = append(children, e)
	}
	return children
}

func escapeContainedMDV1(data []rune, mdType []rune) string {
	out := strings.Builder{}
	for _, x := range data {
		if contains(x, mdType) {
			out.WriteRune('\\')
		}
		out.WriteRune(x)
	}
	return out.String()
}

func contains(r rune, rs []rune) bool {
	for _, rr := range rs {
		if r == rr {
			return true
		}
	}
	return false
}
