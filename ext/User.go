package ext

type User struct {
	Bot          Bot
	Id           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"`
	Photos     [][]PhotoSize `json:"photos"`
}

func (user User) GetProfilePhotos(offset int, limit int) (*UserProfilePhotos, error) {
	return user.Bot.GetUserProfilePhotos(user.Id)
}
