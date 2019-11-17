package ext

// AnswerCallbackQuery answer a callback query
func (b Bot) AnswerCallbackQuery(callbackQueryId string) (bool, error) {
	return b.NewSendableAnswerCallbackQuery(callbackQueryId).Send()
}

// AnswerCallbackQueryText answer a callback query with text
func (b Bot) AnswerCallbackQueryText(callbackQueryId string, text string, alert bool) (bool, error) {
	cbq := b.NewSendableAnswerCallbackQuery(callbackQueryId)
	cbq.Text = text
	cbq.ShowAlert = alert
	return cbq.Send()
}

// AnswerCallbackQueryURL answer a callback query with a URL
func (b Bot) AnswerCallbackQueryURL(callbackQueryId string, url string) (bool, error) {
	cbq := b.NewSendableAnswerCallbackQuery(callbackQueryId)
	cbq.Url = url
	return cbq.Send()
}
