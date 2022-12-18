package handlers_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"unicode/utf16"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func NewTestBot() *gotgbot.Bot {
	server := httptest.NewServer(nil)
	return &gotgbot.Bot{
		User: gotgbot.User{
			Id:        0,
			IsBot:     false,
			FirstName: "gobot",
			LastName:  "",
			Username:  "gotgbot",
		},
		BotClient: &gotgbot.BaseBotClient{
			Token:  "use-me",
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: 0,
				APIURL:  server.URL,
			},
		},
	}
}

func NewMessage(userId int64, chatId int64, message string) *ext.Context {
	return newMessage(userId, chatId, message, nil)
}

func NewCommandMessage(userId int64, chatId int64, command string, args []string) *ext.Context {
	msg, ents := buildCommand(command, args)
	return newMessage(userId, chatId, msg, ents)
}

func buildCommand(cmd string, args []string) (string, []gotgbot.MessageEntity) {
	return fmt.Sprintf("/%s %s", cmd, strings.Join(args, " ")),
		[]gotgbot.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: int64(len(utf16.Encode([]rune("/" + cmd)))),
			},
		}
}

func newMessage(userId int64, chatId int64, message string, entities []gotgbot.MessageEntity) *ext.Context {
	chatType := "supergroup"
	if userId == chatId {
		chatType = "private"
	}

	return ext.NewContext(&gotgbot.Update{
		UpdateId: rand.Int63(), // should this be consistent?
		Message: &gotgbot.Message{
			MessageId: rand.Int63(), // should this be consistent?
			Text:      message,
			Entities:  entities,
			From: &gotgbot.User{
				Id:        userId,
				FirstName: "bob",
			},
			Chat: gotgbot.Chat{
				Id:    chatId,
				Title: "bob's chat",
				Type:  chatType,
			},
		},
	}, nil)
}
