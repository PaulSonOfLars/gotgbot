package gotgbot

// The consts listed below represent all the parse_mode options that can be sent to telegram.
const (
	ParseModeHTML       = "HTML"
	ParseModeMarkdownV2 = "MarkdownV2"
	ParseModeMarkdown   = "Markdown"
	ParseModeNone       = ""
)

// ChatMemberStatusOwner is an alias to handle the commonly confused type names.
const ChatMemberStatusOwner = ChatMemberStatusCreator
