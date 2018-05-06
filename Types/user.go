package Types

type User struct {
	Id            int
	Is_bot        bool
	First_name    string
	Last_name     string
	Username      string
	Language_code string
}

type UserProfilePhotos struct {
	Total_count int
	Photos      [][]PhotoSize
}