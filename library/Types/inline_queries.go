package Types

type InlineQuery struct {
	Id       string
	From     User
	Location Location
	Query    string
	Offset   string
}

type InlineQueryResult struct {}

type InlineQueryResultArticle struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Title                 string
	Input_message_content InputMessageContent
	Reply_markup          InlineKeyboardMarkup
	Url                   string
	Hide_url              bool
	Description           string
	Thumb_url             string
	Thumb_width           int
	Thumb_height          int
}

type InlineQueryResultPhoto struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Photo_url             string
	Thumb_url             string
	Photo_width           int
	Photo_height          int
	Title                 string
	Description           string
	Caption               string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent
}

type InlineQueryResultGif struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Gif_url               string
	Gif_width             int
	Gif_height            int
	Gif_duration          int
	Thumb_url             string
	Title                 string
	Caption               string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent
}

type InlineQueryResultMpeg4Gif struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Mpeg4_url             string
	Mpeg4_width           int
	Mpeg4_height          int
	Mpeg4_duration        int
	Thumb_url             string
	Title                 string
	Caption               string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent
}

type InlineQueryResultVideo struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Video_url             string
	Mime_type             string
	Thumb_url             string
	Title                 string
	Caption               string
	Video_width           int
	Video_height          int
	Video_duration        int
	Description           string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InlineQueryResultAudio struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Audio_url             string
	Title                 string
	Caption               string
	Performer             string
	Audio_duration        int
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InlineQueryResultVoice struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Voice_url             string
	Title                 string
	Caption               string
	Voice_duration        int
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InlineQueryResultDocument struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Title                 string
	Caption               string
	Document_url          string
	Mime_type             string
	Description           string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent
	Thumb_url             string
	Thumb_width           int
	Thumb_height          int

}

type InlineQueryResultLocation struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Latitude              float64
	Longitude             float64
	Title                 string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent
	Thumb_url             string
	Thumb_width           int
	Thumb_height          int

}

type InlineQueryResultVenue struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Latitude              float64
	Longitude             float64
	Title                 string
	Address               string
	Foursquare_id         string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent
	Thumb_url             string
	Thumb_width           int
	Thumb_height          int

}

type InlineQueryResultContact struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Phone_number          string
	First_name            string
	Last_name             string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent
	Thumb_url             string
	Thumb_width           int
	Thumb_height          int

}

type InlineQueryResultGame struct {
	InlineQueryResult
	Type            string
	Id              string
	Game_short_name string
	Reply_markup    InlineKeyboardMarkup

}

type InlineQueryResultCachedPhoto struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Photo_file_id         string
	Title                 string
	Description           string
	Caption               string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InlineQueryResultCachedGif struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Gif_file_id           string
	Title                 string
	Caption               string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InlineQueryResultCachedMpeg4Gif struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Mpeg4_file_id         string
	Title                 string
	Caption               string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InlineQueryResultCachedSticker struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Sticker_file_id       string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InlineQueryResultCachedDocument struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Title                 string
	Document_file_id      string
	Description           string
	Caption               string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InlineQueryResultCachedVideo struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Video_file_id         string
	Title                 string
	Description           string
	Caption               string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InlineQueryResultCachedVoice struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Voice_file_id         string
	Title                 string
	Caption               string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InlineQueryResultCachedAudio struct {
	InlineQueryResult
	Type                  string
	Id                    string
	Audio_file_id         string
	Caption               string
	Reply_markup          InlineKeyboardMarkup
	Input_message_content InputMessageContent

}

type InputMessageContent struct {}

type InputTextMessageContent struct {
	InputMessageContent
	Message_text             string
	Parse_mode               string
	Disable_web_page_preview bool
}

type InputLocationMessageContent struct {
	InputMessageContent
	Latitude  float64
	Longitude float64

}

type InputVenueMessageContent struct {
	InputMessageContent
	Latitude      float64
	Longitude     float64
	Title         string
	Address       string
	Foursquare_id string
}

type InputContactMessageContent struct {
	InputMessageContent
	Phone_number string
	First_name  string
	Last_name   string

}

type ChosenInlineResult struct {
	Result_id         string
	From              User
	Location          Location
	Inline_message_id string
	Query             string

}
