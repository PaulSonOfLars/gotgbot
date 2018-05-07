package ext

import "gotgbot/types"

type User struct {
	types.User
	bot Bot
}

func (user User) GetProfilePhotos(offset int, limit int) (*types.UserProfilePhotos, error) {
	return user.bot.GetUserProfilePhotos(user.Id)
}