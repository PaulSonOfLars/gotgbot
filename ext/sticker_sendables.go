package ext

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

type File struct {
	bot Bot
	FileId   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	FilePath string `json:"file_path"`
}

type sendableSticker struct {
	bot                 Bot
	ChatId              int
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
	//v.Add("disable_notification", strconv.FormatBool(s.DisableNotification))
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
	if !r.Ok {
		return nil, errors.New(r.Description)
	}

	return s.bot.ParseMessage(r.Result), nil
}

type sendableUploadStickerFile struct {
	bot    Bot
	UserId int
	file
}

func (b Bot) NewSendableUploadStickerFile(userId int) *sendableUploadStickerFile {
	return &sendableUploadStickerFile{bot: b, UserId: userId}
}

func (usf *sendableUploadStickerFile) Send() (*File, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(usf.UserId))

	r, err := usf.bot.sendFile(usf.file, "sticker","uploadStickerFile", v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to uploadStickerFile")
	}
	if !r.Ok {
		return nil, errors.New(r.Description)
	}
	newFile := &File{}
	newFile.bot = usf.bot
	json.Unmarshal(r.Result, newFile)
	return newFile, nil
}

type sendableCreateNewSticker struct {
	bot           Bot
	UserId        int
	Name          string
	Title         string
	file
	Emojis        string
	ContainsMasks bool
	MaskPosition  *MaskPosition
}

func (b Bot) NewSendableCreateNewSticker(userId int, name string, title string, emojis string) *sendableCreateNewSticker {
	return &sendableCreateNewSticker{bot: b, UserId: userId, Name: name, Title: title, Emojis: emojis}
}

func (cns *sendableCreateNewSticker) Send() (bool, error) {
	maskPos, err := json.Marshal(cns.MaskPosition)
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse mask position")
	}

	v := url.Values{}
	v.Add("user_id", strconv.Itoa(cns.UserId))
	v.Add("name", cns.Name)
	v.Add("title", cns.Title)
	v.Add("emojis", cns.Emojis)
	v.Add("contains_mask", strconv.FormatBool(cns.ContainsMasks))
	v.Add("mask_position", string(maskPos))

	r, err := cns.bot.sendFile(cns.file, "sticker","createNewStickerSet", v)
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

type sendableAddStickerToSet struct {
	bot          Bot
	UserId       int
	Name         string
	file
	Emojis       string
	MaskPosition *MaskPosition
}


func (b Bot) NewSendableAddStickerToSet(userId int, name string, emojis string) *sendableAddStickerToSet {
	return &sendableAddStickerToSet{bot: b, UserId: userId, Name: name, Emojis: emojis}
}

func (asts *sendableAddStickerToSet) Send() (bool, error) {
	maskPos, err := json.Marshal(asts.MaskPosition)
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse mask position")
	}

	v := url.Values{}
	v.Add("user_id", strconv.Itoa(asts.UserId))
	v.Add("name", asts.Name)
	v.Add("emojis", asts.Emojis)
	v.Add("mask_position", string(maskPos))

	r, err := asts.bot.sendFile(asts.file, "sticker","addStickerToSet", v)
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
