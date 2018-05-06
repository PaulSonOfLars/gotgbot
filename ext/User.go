package ext

import "gotgbot/types"

type User struct {
	types.User
	bot Bot
}

func (user User) GetProfilePhotos(offset int, limit int) *types.UserProfilePhotos {
	return user.bot.GetUserProfilePhotos(user.Id)
}