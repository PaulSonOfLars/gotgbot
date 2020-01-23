package ext

type User struct {
	Bot                     Bot    `json:"-"`
	Id                      int    `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Username                string `json:"username"`
	LanguageCode            string `json:"language_code"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
}

type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"`
	Photos     [][]PhotoSize `json:"photos"`
}

// GetProfilePhotos Retrieves a user's profile pictures
func (user User) GetProfilePhotos(offset int, limit int) (*UserProfilePhotos, error) {
	return user.Bot.GetUserProfilePhotos(user.Id)
}
