package ext

type Game struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Photo        []PhotoSize     `json:"photo"`
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities"`
	Animation    Animation       `json:"animation"`
}

type Animation struct {
	FileId   string    `json:"file_id"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Duration int       `json:"duration"`
	Thumb    PhotoSize `json:"thumb"`
	FileName string    `json:"file_name"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

type GameHighScore struct {
	Position int  `json:"position"`
	User     User `json:"user"`
	Score    int  `json:"score"`
}

func (b Bot) SendGame(chatId int, gameShortName string) (*Message, error) {
	return b.NewSendableGame(chatId, gameShortName).Send()
}

func (b Bot) SetGameScore(userId int, score int, chatId int, messageId int) (bool, error) {
	return b.NewSendableSetGameScore(userId, score, chatId, messageId).Send()
}

func (b Bot) SetGameScoreInline(userId int, score int, inlineMessageId string) (bool, error) {
	return b.NewSendableSetGameScoreInline(userId, score, inlineMessageId).Send()
}

func (b Bot) GetGameHighScores(userId int, chatId int, messageId int) ([]GameHighScore, error) {
	return b.NewSendableGetGameHighScore(userId, chatId, messageId).Send()
}

func (b Bot) GetGameHighScoresInline(userId int, inlineMessageId string) ([]GameHighScore, error) {
	return b.NewSendableGetGameHighScoreInline(userId, inlineMessageId).Send()
}
