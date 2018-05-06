package types

type Sticker struct {
	FileId       string       `json:"file_id"`
	Width        int          `json:"width"`
	Height       int          `json:"height"`
	Thumb        PhotoSize    `json:"thumb"`
	Emoji        string       `json:"emoji"`
	SetName      string       `json:"set_name"`
	MaskPosition MaskPosition `json:"mask_position"`
	FileSize     int          `json:"file_size"`
}

type StickerSet struct {
	Name          string    `json:"name"`
	Title         string    `json:"title"`
	ContainsMasks bool      `json:"contains_masks"`
	Stickers      []Sticker `json:"stickers"`
}

type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float64 `json:"x_shift"`
	YShift float64 `json:"y_shift"`
	Scale  float64 `json:"scale"`
}
