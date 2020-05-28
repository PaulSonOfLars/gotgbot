package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type File struct {
	bot          Bot    `json:"-"`
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FilePath     string `json:"file_path"`
}

type sendableSticker struct {
	bot    Bot `json:"-"`
	ChatId int
	file
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         ReplyMarkup
}

func (b Bot) NewSendableSticker(chatId int) *sendableSticker {
	return &sendableSticker{bot: b, ChatId: chatId}
}

func (s *sendableSticker) Send() (*Message, error) {
	replyMarkup := []byte("{}")
	if s.ReplyMarkup != nil {
		var err error
		replyMarkup, err = s.ReplyMarkup.Marshal()
		if err != nil {
			return nil, err
		}
	}

	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(s.ChatId))
	// v.Add("disable_notification", strconv.FormatBool(s.DisableNotification))
	if s.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(s.ReplyToMessageId))
	}
	if s.ReplyMarkup != nil {
		v.Add("reply_markup", string(replyMarkup))
	}

	r, err := s.bot.sendFile(s.file, "sticker", "sendSticker", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sendSticker")
	}

	return s.bot.ParseMessage(r.Result)
}

type sendableUploadStickerFile struct {
	bot    Bot `json:"-"`
	UserId int
	file
}

func (b Bot) NewSendableUploadStickerFile(userId int) *sendableUploadStickerFile {
	return &sendableUploadStickerFile{bot: b, UserId: userId}
}

func (usf *sendableUploadStickerFile) Send() (*File, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(usf.UserId))

	r, err := usf.bot.sendFile(usf.file, "sticker", "uploadStickerFile", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to uploadStickerFile")
	}

	newFile := &File{}
	newFile.bot = usf.bot
	return newFile, json.Unmarshal(r.Result, newFile)
}

// TODO: check whether uploading tgs_stickers works
type sendableCreateNewStickerSet struct {
	bot         Bot    `json:"-"`
	StickerType string `json:"-"` // "png_sticker" or "tgs_sticker"
	UserId      int
	Name        string
	Title       string
	file
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

	if cns.StickerType == "" {
		cns.StickerType = "png_sticker"
	}

	v := url.Values{}
	v.Add("user_id", strconv.Itoa(cns.UserId))
	v.Add("name", cns.Name)
	v.Add("title", cns.Title)
	v.Add("emojis", cns.Emojis)
	v.Add("contains_mask", strconv.FormatBool(cns.ContainsMasks))
	v.Add("mask_position", string(maskPos))

	r, err := cns.bot.sendFile(cns.file, cns.StickerType, "createNewStickerSet", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to createNewStickerSet")
	}

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

type sendableAddStickerToSet struct {
	bot         Bot    `json:"-"`
	StickerType string `json:"-"` // "png_sticker" or "tgs_sticker"
	UserId      int
	Name        string
	file
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

	if asts.StickerType == "" {
		asts.StickerType = "png_sticker"
	}

	v := url.Values{}
	v.Add("user_id", strconv.Itoa(asts.UserId))
	v.Add("name", asts.Name)
	v.Add("emojis", asts.Emojis)
	v.Add("mask_position", string(maskPos))

	r, err := asts.bot.sendFile(asts.file, asts.StickerType, "addStickerToSet", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to addStickerToSet")
	}

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}

type sendableSetStickerSetThumb struct {
	bot    Bot `json:"-"`
	UserId int
	file
}

func (b Bot) NewSendableSetStickerSetThumb(userId int) *sendableSetStickerSetThumb {
	return &sendableSetStickerSetThumb{bot: b, UserId: userId}
}

func (ssst *sendableSetStickerSetThumb) Send() (bool, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(ssst.UserId))

	r, err := ssst.bot.sendFile(ssst.file, "sticker", "setStickerSetThumb", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to setStickerSetThumb")
	}

	var bb bool
	return bb, json.Unmarshal(r.Result, &bb)
}
