package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type File struct {
	bot          Bot
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FilePath     string `json:"file_path"`
}

type sendableSticker struct {
	bot                 Bot
	ChatId              int
	Sticker             InputFile
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (b Bot) NewSendableSticker(chatId int) *sendableSticker {
	return &sendableSticker{bot: b, ChatId: chatId}
}

func (s *sendableSticker) Send() (*Message, error) {
	var replyMarkup []byte
	if s.ReplyMarkup != nil {
		var err error
		replyMarkup, err = s.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(s.ChatId))
	v.Add("disable_notification", strconv.FormatBool(s.DisableNotification))
	v.Add("reply_to_message_id", strconv.Itoa(s.ReplyToMessageId))
	v.Add("reply_markup", string(replyMarkup))

	r, err := s.Sticker.send("sendSticker", v, "sticker")
	if err != nil {
		return nil, err
	}

	return s.bot.ParseMessage(r)
}

type sendableUploadStickerFile struct {
	bot        Bot
	UserId     int
	PngSticker InputFile
}

func (b Bot) NewSendableUploadStickerFile(userId int) *sendableUploadStickerFile {
	return &sendableUploadStickerFile{bot: b, UserId: userId}
}

func (usf *sendableUploadStickerFile) Send() (*File, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(usf.UserId))

	r, err := usf.PngSticker.send("uploadStickerFile", v, "png_sticker")
	if err != nil {
		return nil, err
	}

	newFile := &File{}
	newFile.bot = usf.bot
	return newFile, json.Unmarshal(r, newFile)
}

type sendableCreateNewStickerSet struct {
	bot           Bot
	UserId        int
	Name          string
	Title         string
	PngSticker    *InputFile
	TgsSticker    *InputFile
	Emojis        string
	ContainsMasks bool
	MaskPosition  *MaskPosition
}

func (b Bot) NewSendableCreateNewStickerSet(userId int, name string, title string, emojis string) *sendableCreateNewStickerSet {
	return &sendableCreateNewStickerSet{bot: b, UserId: userId, Name: name, Title: title, Emojis: emojis}
}

func (cns *sendableCreateNewStickerSet) Send() (bool, error) {
	var maskPos []byte
	if cns.MaskPosition != nil {
		var err error
		maskPos, err = json.Marshal(cns.MaskPosition)
		if err != nil {
			return false, errors.Wrapf(err, "failed to parse mask position")
		}
	}

	v := url.Values{}
	v.Add("user_id", strconv.Itoa(cns.UserId))
	v.Add("name", cns.Name)
	v.Add("title", cns.Title)
	v.Add("emojis", cns.Emojis)
	v.Add("contains_mask", strconv.FormatBool(cns.ContainsMasks))
	v.Add("mask_position", string(maskPos))

	var r json.RawMessage
	var err error
	if cns.PngSticker != nil && cns.TgsSticker != nil {
		return false, errors.New("can only specify one stickertype; png or tgs")
	} else if cns.PngSticker != nil {
		r, err = cns.PngSticker.send("createNewStickerSet", v, "png_sticker")
	} else {
		r, err = cns.TgsSticker.send("createNewStickerSet", v, "tgs_sticker")
	}
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

type sendableAddStickerToSet struct {
	bot          Bot
	UserId       int
	Name         string
	PngSticker   *InputFile
	TgsSticker   *InputFile
	Emojis       string
	MaskPosition *MaskPosition
}

func (b Bot) NewSendableAddStickerToSet(userId int, name string, emojis string) *sendableAddStickerToSet {
	return &sendableAddStickerToSet{bot: b, UserId: userId, Name: name, Emojis: emojis}
}

func (asts *sendableAddStickerToSet) Send() (bool, error) {
	var maskPos []byte
	if asts.MaskPosition != nil {
		var err error
		maskPos, err = json.Marshal(asts.MaskPosition)
		if err != nil {
			return false, errors.Wrapf(err, "failed to parse mask position")
		}
	}

	v := url.Values{}
	v.Add("user_id", strconv.Itoa(asts.UserId))
	v.Add("name", asts.Name)
	v.Add("emojis", asts.Emojis)
	v.Add("mask_position", string(maskPos))

	var r json.RawMessage
	var err error
	if asts.PngSticker != nil && asts.TgsSticker != nil {
		return false, errors.New("can only specify one stickertype; png or tgs")
	} else if asts.PngSticker != nil {
		r, err = asts.PngSticker.send("addStickerToSet", v, "png_sticker")
	} else {
		r, err = asts.TgsSticker.send("addStickerToSet", v, "tgs_sticker")
	}
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

type sendableSetStickerSetThumb struct {
	bot    Bot
	UserId int
	Thumb  InputFile
}

func (b Bot) NewSendableSetStickerSetThumb(userId int) *sendableSetStickerSetThumb {
	return &sendableSetStickerSetThumb{bot: b, UserId: userId}
}

func (ssst *sendableSetStickerSetThumb) Send() (bool, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(ssst.UserId))

	r, err := ssst.Thumb.send("setStickerSetThumb", v, "thumb")
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}
