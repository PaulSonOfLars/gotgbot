package ext

import (
	"github.com/PaulSonOfLars/gotgbot/types"
	"encoding/json"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

type File struct {
	types.File
	bot Bot
}

type sendableSticker struct {
	bot                 Bot
	ChatId              int
	file
	DisableNotification bool
	ReplyToMessageId    int
	ReplyMarkup         types.ReplyKeyboardMarkup
}

func (b Bot) NewSendableSticker(chatId int) *sendableSticker {
	return &sendableSticker{bot: b, ChatId: chatId}
}

func (s *sendableSticker) Send() (*Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(s.ChatId))

	r, err := sendFile(s.file, "sendSticker", v)
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

	r, err := sendFile(usf.file, "uploadStickerFile", v)
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
