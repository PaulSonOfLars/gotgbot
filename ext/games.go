package ext

import (
	"github.com/PaulSonOfLars/gotgbot/types"
)

func (b Bot) SendGame(chatId int, gameShortName string) (*Message, error) {
	return b.NewSendableGame(chatId, gameShortName).Send()
}

func (b Bot) SetGameScore(userId int, score int, chatId int, messageId int) (bool, error) {
	return b.NewSendableSetGameScore(userId, score, chatId, messageId).Send()
}

func (b Bot) SetGameScoreInline(userId int, score int, inlineMessageId string) (bool, error) {
	return b.NewSendableSetGameScoreInline(userId, score, inlineMessageId).Send()
}

func (b Bot) GetGameHighScores(userId int, chatId int, messageId int) ([]types.GameHighScore, error) {
	return b.NewSendableGetGameHighScore(userId, chatId, messageId).Send()
}

func (b Bot) GetGameHighScoresInline(userId int, inlineMessageId string) ([]types.GameHighScore, error) {
	return b.NewSendableGetGameHighScoreInline(userId, inlineMessageId).Send()
}
