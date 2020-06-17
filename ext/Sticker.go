package ext

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Sticker struct {
	FileId       string       `json:"file_id"`
	FileUniqueId string       `json:"file_unique_id"`
	Width        int          `json:"width"`
	Height       int          `json:"height"`
	IsAnimated   bool         `json:"is_animated"`
	Thumb        *PhotoSize   `json:"thumb"`
	Emoji        string       `json:"emoji"`
	SetName      string       `json:"set_name"`
	MaskPosition MaskPosition `json:"mask_position"`
	FileSize     int          `json:"file_size"`
}

type StickerSet struct {
	Name          string     `json:"name"`
	Title         string     `json:"title"`
	IsAnimated    bool       `json:"is_animated"`
	ContainsMasks bool       `json:"contains_masks"`
	Stickers      []Sticker  `json:"stickers"`
	Thumb         *PhotoSize `json:"thumb"`
}

type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float64 `json:"x_shift"`
	YShift float64 `json:"y_shift"`
	Scale  float64 `json:"scale"`
}

func (b Bot) SendSticker(chatId int, s InputFile) (*Message, error) {
	sticker := b.NewSendableSticker(chatId)
	sticker.InputFile = s
	return sticker.Send()
}

func (b Bot) ReplySticker(chatId int, s InputFile, replyToMessageId int) (*Message, error) {
	sticker := b.NewSendableSticker(chatId)
	sticker.InputFile = s
	sticker.ReplyToMessageId = replyToMessageId
	return sticker.Send()
}

func (b Bot) GetStickerSet(name string) (*StickerSet, error) {
	v := url.Values{}
	v.Add("name", name)

	r, err := b.Get("getStickerSet", v)
	if err != nil {
		return nil, err
	}
	var ss StickerSet
	return &ss, json.Unmarshal(r, &ss)
}

func (b Bot) UploadSticker(userId int, s InputFile) (*File, error) {
	uploadSticker := b.NewSendableUploadStickerFile(userId)
	uploadSticker.InputFile = s
	return uploadSticker.Send()
}

func (b Bot) CreateNewStickerSet(userId int, name string, title string, pngSticker InputFile, emojis string) (bool, error) {
	createNew := b.NewSendableCreateNewStickerSet(userId, name, title, emojis)
	createNew.InputFile = pngSticker
	return createNew.Send()
}

func (b Bot) AddStickerToSet(userId int, name string, pngSticker InputFile, emojis string) (bool, error) {
	addSticker := b.NewSendableAddStickerToSet(userId, name, emojis)
	addSticker.InputFile = pngSticker
	return addSticker.Send()
}

func (b Bot) SetStickerPositionInSet(sticker string, position int) (bool, error) {
	v := url.Values{}
	v.Add("sticker", sticker)
	v.Add("position", strconv.Itoa(position))

	r, err := b.Get("setStickerPositionInSet", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

func (b Bot) DeleteStickerFromSet(sticker string) (bool, error) {
	v := url.Values{}
	v.Add("sticker", sticker)

	r, err := b.Get("deleteStickerFromSet", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}
