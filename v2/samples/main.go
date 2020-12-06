package main

import (
	"net/http"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func main() {
	b := gotgbot.Bot{
		Token:  os.Getenv("TOKEN"),
		User:   gotgbot.User{},
		Client: http.Client{},
	}

	//sendTestMsg(b)
	//
	//sendDocMsg(b)
	//
	//sendMediaGroupMsg(b)
	//
	sendMixedMediaGroupMsg(b)
}

func sendMediaGroupMsg(b gotgbot.Bot) {
	f1, err := os.Open("api.json")
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	f2, err := os.Open("api.json")
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	var inputMedia []gotgbot.InputMedia
	inputMedia = append(inputMedia, gotgbot.InputMediaDocument{
		Media: gotgbot.NamedFile{
			File:     f1,
			FileName: "first",
		},
		Caption: "boo1",
	}, gotgbot.InputMediaDocument{
		Media:   f2,
		Caption: "boo2",
	})

	_, err = b.SendMediaGroup(-1001328491754, inputMedia, gotgbot.SendMediaGroupOpts{})
	if err != nil {
		panic(err)
	}
}

func sendMixedMediaGroupMsg(b gotgbot.Bot) {
	m4a, err := os.Open("test.MP4")
	if err != nil {
		panic(err)
	}
	defer m4a.Close()

	jpg, err := os.Open("test.jpg")
	if err != nil {
		panic(err)
	}
	defer jpg.Close()

	var inputMedia []gotgbot.InputMedia
	inputMedia = append(inputMedia, gotgbot.InputMediaVideo{
		Media: gotgbot.NamedFile{
			File:     m4a,
			FileName: "first",
		},
		Caption: "boo1",
	}, gotgbot.InputMediaPhoto{
		Media:   jpg,
		Caption: "boo2",
	})

	_, err = b.SendMediaGroup(-1001328491754, inputMedia, gotgbot.SendMediaGroupOpts{})
	if err != nil {
		panic(err)
	}
}

func sendDocMsg(b gotgbot.Bot) {
	f, err := os.Open("api.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = b.SendDocument(-1001328491754, f, gotgbot.SendDocumentOpts{
		Thumb:                       nil,
		Caption:                     "api with a caption",
		ParseMode:                   "",
		CaptionEntities:             nil,
		DisableContentTypeDetection: false,
		DisableNotification:         false,
		ReplyToMessageId:            0,
		AllowSendingWithoutReply:    false,
		ReplyMarkup:                 gotgbot.InlineKeyboardMarkup{},
	})
	if err != nil {
		panic(err)
	}
}

func sendTestMsg(b gotgbot.Bot) {
	_, err := b.SendMessage(-1001328491754, "hey _italics_", gotgbot.SendMessageOpts{
		ParseMode:                "Markdown",
		Entities:                 nil,
		DisableWebPagePreview:    false,
		DisableNotification:      false,
		ReplyToMessageId:         0,
		AllowSendingWithoutReply: false,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{
				Text:                         "heyo",
				Url:                          "google.com",
				LoginUrl:                     nil,
				CallbackData:                 "",
				SwitchInlineQuery:            "",
				SwitchInlineQueryCurrentChat: "",
				CallbackGame:                 nil,
				Pay:                          false,
			}}},
		},
	})
	if err != nil {
		panic(err)
	}
}
