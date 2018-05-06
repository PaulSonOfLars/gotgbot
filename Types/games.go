package Types

type Game struct {
	Title         string
	Description   string
	Photo         []PhotoSize
	Text          string
	Text_entities []MessageEntity
	Animation     Animation
}

type Animation struct {
	File_id   string
	Thumb     PhotoSize
	File_name string
	Mime_type string
	File_size int
}

type GameHighScore struct {
	Position int
	User     User
	Score    int
}
