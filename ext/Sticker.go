package ext

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
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
	bot           Bot        `json:"-"`
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

func (b Bot) SendStickerStr(chatId int, stickerId string) (*Message, error) {
	sticker := b.NewSendableSticker(chatId)
	sticker.FileId = stickerId
	return sticker.Send()
}

func (b Bot) SendStickerPath(chatId int, path string) (*Message, error) {
	sticker := b.NewSendableSticker(chatId)
	sticker.Path = path
	return sticker.Send()
}

func (b Bot) SendStickerReader(chatId int, reader io.Reader) (*Message, error) {
	sticker := b.NewSendableSticker(chatId)
	sticker.Reader = reader
	return sticker.Send()
}

func (b Bot) ReplyStickerStr(chatId int, stickerId string, replyToMessageId int) (*Message, error) {
	sticker := b.NewSendableSticker(chatId)
	sticker.FileId = stickerId
	sticker.ReplyToMessageId = replyToMessageId
	return sticker.Send()
}

func (b Bot) ReplyStickerPath(chatId int, path string, replyToMessageId int) (*Message, error) {
	sticker := b.NewSendableSticker(chatId)
	sticker.Path = path
	sticker.ReplyToMessageId = replyToMessageId
	return sticker.Send()
}

func (b Bot) ReplyStickerReader(chatId int, reader io.Reader, replyToMessageId int) (*Message, error) {
	sticker := b.NewSendableSticker(chatId)
	sticker.Reader = reader
	sticker.ReplyToMessageId = replyToMessageId
	return sticker.Send()
}

func (b Bot) GetStickerSet(name string) (*StickerSet, error) {
	v := url.Values{}
	v.Add("name", name)

	r, err := b.Get("getStickerSet", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to getStickerSet")
	}
	var ss StickerSet
	return &ss, json.Unmarshal(r, &ss)
}

func (b Bot) UploadStickerFileStr(userId int, pngStickerId string) (*File, error) {
	uploadSticker := b.NewSendableUploadStickerFile(userId)
	uploadSticker.FileId = pngStickerId
	return uploadSticker.Send()
}

func (b Bot) UploadStickerFilePath(userId int, path string) (*File, error) {
	uploadSticker := b.NewSendableUploadStickerFile(userId)
	uploadSticker.Path = path
	return uploadSticker.Send()
}

func (b Bot) UploadStickerFileReader(userId int, reader io.Reader) (*File, error) {
	uploadSticker := b.NewSendableUploadStickerFile(userId)
	uploadSticker.Reader = reader
	return uploadSticker.Send()
}

func (b Bot) CreateNewStickerSetStr(userId int, name string, title string, pngStickerid string, emojis string) (bool, error) {
	createNew := b.NewSendableCreateNewStickerSet(userId, name, title, emojis)
	createNew.FileId = pngStickerid
	return createNew.Send()
}

func (b Bot) CreateNewStickerSetPath(userId int, name string, title string, path string, emojis string) (bool, error) {
	createNew := b.NewSendableCreateNewStickerSet(userId, name, title, emojis)
	createNew.Path = path
	return createNew.Send()
}

func (b Bot) CreateNewStickerSetReader(userId int, name string, title string, reader io.Reader, emojis string) (bool, error) {
	createNew := b.NewSendableCreateNewStickerSet(userId, name, title, emojis)
	createNew.Reader = reader
	return createNew.Send()
}

func (b Bot) AddStickerToSetStr(userId int, name string, pngStickerId string, emojis string) (bool, error) {
	addSticker := b.NewSendableAddStickerToSet(userId, name, emojis)
	addSticker.FileId = pngStickerId
	return addSticker.Send()
}

func (b Bot) AddStickerToSetPath(userId int, name string, path string, emojis string) (bool, error) {
	addSticker := b.NewSendableAddStickerToSet(userId, name, emojis)
	addSticker.Path = path
	return addSticker.Send()
}

func (b Bot) AddStickerToSetReader(userId int, name string, reader io.Reader, emojis string) (bool, error) {
	addSticker := b.NewSendableAddStickerToSet(userId, name, emojis)
	addSticker.Reader = reader
	return addSticker.Send()
}

func (b Bot) SetStickerPositionInSet(sticker string, position int) (bool, error) {
	v := url.Values{}
	v.Add("sticker", sticker)
	v.Add("position", strconv.Itoa(position))

	r, err := b.Get("setStickerPositionInSet", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to setStickerPositionInSet")
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

func (b Bot) DeleteStickerFromSet(sticker string) (bool, error) {
	v := url.Values{}
	v.Add("sticker", sticker)

	r, err := b.Get("deleteStickerFromSet", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to deleteStickerFromSet")
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}
