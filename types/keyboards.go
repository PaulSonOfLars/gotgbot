package types

type ReplyKeyboardMarkup struct {
	Keyboard          [][]KeyboardButton
	Resize_keyboard   bool
	One_time_keyboard bool
	Selective         bool
}

type KeyboardButton struct {
	Text             string
	Request_contact  bool
	Request_location bool
}

type ReplyKeyboardRemove struct {
	Remove_keyboard bool
	Selective       bool
}

type InlineKeyboardMarkup struct {
	Inline_keyboard [][]InlineKeyboardButton
}

type InlineKeyboardButton struct {
	Text                             string
	Url                              string
	Callback_data                    string
	Switch_inline_query              string
	Switch_inline_query_current_chat string
	//Callback_game                    *CallbackGame
	Pay                              bool
}

type CallbackQuery struct {
	Id                string
	From              *User
	Message           *Message
	Inline_message_id string
	Chat_instance     string
	Data              string
	Game_short_name   string
}

type ForceReply struct {
	ForceReply bool
	Selective  bool
}