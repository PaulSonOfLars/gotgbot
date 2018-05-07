package ext

import (
	"gotgbot/types"
	"strconv"
	"encoding/json"
	"net/url"
	"github.com/pkg/errors"
)

// TODO: inputfile version
// TODO: reply_markup version
func (b Bot) SendStickerStr(chatId int, sticker string) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("sticker", sticker)

	r, err := Get(b, "sendSticker", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendSticker")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return b.ParseMessage(r.Result), nil

}

func (b Bot) ReplyStickerStr(chatId int, sticker string, replyToMessageId int) (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("sticker", sticker)
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r, err := Get(b, "sendSticker", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendSticker")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return b.ParseMessage(r.Result), nil

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

// TODO: input file stuff
//func (b Bot) UploadStickerFile(user_id int, png_sticker types.InputFile) types.File {
//
//}

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
