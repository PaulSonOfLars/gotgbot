package ext

import (
	"github.com/PaulSonOfLars/gotgbot/types"
	"strconv"
	"encoding/json"
	"net/url"
	"github.com/pkg/errors"
	"io"
)

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

func (b Bot) GetStickerSet(name string) (*types.StickerSet, error) {
	v := url.Values{}
	v.Add("name", name)

	r, err := Get(b, "getStickerSet", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to getStickerSet")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	var ss types.StickerSet
	json.Unmarshal(r.Result, &ss)

	return &ss, nil
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

// TODO: contains mask + mask position version
// TODO: InputFile version
// TODO: check return
func (b Bot) CreateNewStickerSetStr(userId int, name string, title string, pngSticker string, emojis string) (bool, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(userId))
	v.Add("name", name)
	v.Add("title", title)
	v.Add("png_sticker", pngSticker)
	v.Add("emojis", emojis)

	r, err := Get(b, "createNewStickerSet", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to createNewStickerSet")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

// TODO: InputFile version
// TODO: mask position version
func (b Bot) AddStickerToSetStr(userId int, name string, pngSticker string, emojis string) (bool, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(userId))
	v.Add("name", name)
	v.Add("png_sticker", pngSticker)
	v.Add("emojis", emojis)

	r, err := Get(b, "addStickerToSet", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to addStickerToSet")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

func (b Bot) SetStickerPositionInSet(sticker string, position int) (bool, error) {
	v := url.Values{}
	v.Add("sticker", sticker)
	v.Add("position", strconv.Itoa(position))

	r, err := Get(b, "setStickerPositionInSet", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to setStickerPositionInSet")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}

func (b Bot) DeleteStickerFromSet(sticker string) (bool, error) {
	v := url.Values{}
	v.Add("sticker", sticker)

	r, err := Get(b, "deleteStickerFromSet", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to deleteStickerFromSet")
	}
	if !r.Ok {
		return false, errors.New(r.Description)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb, nil
}
