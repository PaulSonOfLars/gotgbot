package types

type InlineQuery struct {
	Id       string   `json:"id"`
	From     User     `json:"from"`
	Location Location `json:"location"`
	Query    string   `json:"query"`
	Offset   string   `json:"offset"`
}

type InlineQueryResult struct{}

type InlineQueryResultArticle struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	Title               string               `json:"title"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	Url                 string               `json:"url"`
	HideUrl             bool                 `json:"hide_url"`
	Description         string               `json:"description"`
	ThumbUrl            string               `json:"thumb_url"`
	ThumbWidth          int                  `json:"thumb_width"`
	ThumbHeight         int                  `json:"thumb_height"`
}

type InlineQueryResultPhoto struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	PhotoUrl            string               `json:"photo_url"`
	ThumbUrl            string               `json:"thumb_url"`
	PhotoWidth          int                  `json:"photo_width"`
	PhotoHeight         int                  `json:"photo_height"`
	Title               string               `json:"title"`
	Description         string               `json:"description"`
	Caption             string               `json:"caption"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultGif struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	GifUrl              string               `json:"gif_url"`
	GifWidth            int                  `json:"gif_width"`
	GifHeight           int                  `json:"gif_height"`
	GifDuration         int                  `json:"gif_duration"`
	ThumbUrl            string               `json:"thumb_url"`
	Title               string               `json:"title"`
	Caption             string               `json:"caption"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultMpeg4Gif struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	Mpeg4Url            string               `json:"mpeg4_url"`
	Mpeg4Width          int                  `json:"mpeg4_width"`
	Mpeg4Height         int                  `json:"mpeg4_height"`
	Mpeg4Duration       int                  `json:"mpeg4_duration"`
	ThumbUrl            string               `json:"thumb_url"`
	Title               string               `json:"title"`
	Caption             string               `json:"caption"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultVideo struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	VideoUrl            string               `json:"video_url"`
	MimeType            string               `json:"mime_type"`
	ThumbUrl            string               `json:"thumb_url"`
	Title               string               `json:"title"`
	Caption             string               `json:"caption"`
	VideoWidth          int                  `json:"video_width"`
	VideoHeight         int                  `json:"video_height"`
	VideoDuration       int                  `json:"video_duration"`
	Description         string               `json:"description"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultAudio struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	AudioUrl            string               `json:"audio_url"`
	Title               string               `json:"title"`
	Caption             string               `json:"caption"`
	Performer           string               `json:"performer"`
	AudioDuration       int                  `json:"audio_duration"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultVoice struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	VoiceUrl            string               `json:"voice_url"`
	Title               string               `json:"title"`
	Caption             string               `json:"caption"`
	VoiceDuration       int                  `json:"voice_duration"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultDocument struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	Title               string               `json:"title"`
	Caption             string               `json:"caption"`
	DocumentUrl         string               `json:"document_url"`
	MimeType            string               `json:"mime_type"`
	Description         string               `json:"description"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
	ThumbUrl            string               `json:"thumb_url"`
	ThumbWidth          int                  `json:"thumb_width"`
	ThumbHeight         int                  `json:"thumb_height"`
}

type InlineQueryResultLocation struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	Latitude            float64              `json:"latitude"`
	Longitude           float64              `json:"longitude"`
	Title               string               `json:"title"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
	ThumbUrl            string               `json:"thumb_url"`
	ThumbWidth          int                  `json:"thumb_width"`
	ThumbHeight         int                  `json:"thumb_height"`
}

type InlineQueryResultVenue struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	Latitude            float64              `json:"latitude"`
	Longitude           float64              `json:"longitude"`
	Title               string               `json:"title"`
	Address             string               `json:"address"`
	FoursquareId        string               `json:"foursquare_id"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
	ThumbUrl            string               `json:"thumb_url"`
	ThumbWidth          int                  `json:"thumb_width"`
	ThumbHeight         int                  `json:"thumb_height"`
}

type InlineQueryResultContact struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	PhoneNumber         string               `json:"phone_number"`
	FirstName           string               `json:"first_name"`
	LastName            string               `json:"last_name"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
	ThumbUrl            string               `json:"thumb_url"`
	ThumbWidth          int                  `json:"thumb_width"`
	ThumbHeight         int                  `json:"thumb_height"`
}

type InlineQueryResultGame struct {
	InlineQueryResult
	Type          string               `json:"type"`
	Id            string               `json:"id"`
	GameShortName string               `json:"game_short_name"`
	ReplyMarkup   InlineKeyboardMarkup `json:"reply_markup"`
}

type InlineQueryResultCachedPhoto struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	PhotoFileId         string               `json:"photo_file_id"`
	Title               string               `json:"title"`
	Description         string               `json:"description"`
	Caption             string               `json:"caption"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultCachedGif struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	GifFileId           string               `json:"gif_file_id"`
	Title               string               `json:"title"`
	Caption             string               `json:"caption"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultCachedMpeg4Gif struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	Mpeg4FileId         string               `json:"mpeg4_file_id"`
	Title               string               `json:"title"`
	Caption             string               `json:"caption"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultCachedSticker struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	StickerFileId       string               `json:"sticker_file_id"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultCachedDocument struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	Title               string               `json:"title"`
	DocumentFileId      string               `json:"document_file_id"`
	Description         string               `json:"description"`
	Caption             string               `json:"caption"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultCachedVideo struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	VideoFileId         string               `json:"video_file_id"`
	Title               string               `json:"title"`
	Description         string               `json:"description"`
	Caption             string               `json:"caption"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultCachedVoice struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	VoiceFileId         string               `json:"voice_file_id"`
	Title               string               `json:"title"`
	Caption             string               `json:"caption"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InlineQueryResultCachedAudio struct {
	InlineQueryResult
	Type                string               `json:"type"`
	Id                  string               `json:"id"`
	AudioFileId         string               `json:"audio_file_id"`
	Caption             string               `json:"caption"`
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`
	InputMessageContent InputMessageContent  `json:"input_message_content"`
}

type InputMessageContent struct{}

type InputTextMessageContent struct {
	InputMessageContent
	MessageText           string `json:"message_text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

type InputLocationMessageContent struct {
	InputMessageContent
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type InputVenueMessageContent struct {
	InputMessageContent
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Title        string  `json:"title"`
	Address      string  `json:"address"`
	FoursquareId string  `json:"foursquare_id"`
}

type InputContactMessageContent struct {
	InputMessageContent
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

type ChosenInlineResult struct {
	ResultId        string   `json:"result_id"`
	From            *User     `json:"from"`
	Location        Location `json:"location"`
	InlineMessageId string   `json:"inline_message_id"`
	Query           string   `json:"query"`
}
