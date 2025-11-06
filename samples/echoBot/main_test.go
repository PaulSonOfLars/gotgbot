package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type RequestFunc func(ctx context.Context, token string, method string, params map[string]string, data map[string]gotgbot.FileReader, opts *gotgbot.RequestOpts) (json.RawMessage, error)

// We define a TestBotClient which allows us to overwrite a single method on a per-test basis; meaning each test can verify specific behaviour.
type TestBotClient struct {
	RequestFunc RequestFunc
}

func (b TestBotClient) GetAPIURL(opts *gotgbot.RequestOpts) string {
	// implement me if needed
	panic("not implemented")
}

func (b TestBotClient) FileURL(token string, tgFilePath string, opts *gotgbot.RequestOpts) string {
	// implement me if needed
	panic("not implemented")
}

// Define wrapper around existing RequestWithContext method.
func (b TestBotClient) RequestWithContext(ctx context.Context, token string, method string, params map[string]string, data map[string]gotgbot.FileReader, opts *gotgbot.RequestOpts) (json.RawMessage, error) {
	if b.RequestFunc == nil {
		panic("no requestfunc provided for test")
	}
	return b.RequestFunc(ctx, token, method, params, data, opts)
}

func NewBotClient(f RequestFunc) gotgbot.BotClient {
	return TestBotClient{
		RequestFunc: f,
	}
}

func Test_echo(t *testing.T) {
	const echoMessage = "Hello"
	sendMessageCount := 0

	b := testBot(func(ctx context.Context, token string, method string, params map[string]string, data map[string]gotgbot.FileReader, opts *gotgbot.RequestOpts) (json.RawMessage, error) {
		if method != "sendMessage" {
			t.Fatalf("Only expected API calls to sendMessage, got %s", method)
		}

		sentText := params["text"]
		if sentText != echoMessage {
			t.Errorf("expected text to be %s, got %s", echoMessage, sentText)
		}

		sendMessageCount++
		return json.Marshal(gotgbot.Message{
			MessageId: rand.Int63(),
			Text:      sentText,
		})
	})

	err := echo(b, ext.NewContext(b, testMessage(echoMessage), nil))
	if err != nil {
		t.Fatal(err)
	}
}

func testBot(f RequestFunc) *gotgbot.Bot {
	id := rand.Int63()

	return &gotgbot.Bot{
		Token: fmt.Sprintf("%d:TEST", id),
		User: gotgbot.User{
			Id:        id,
			IsBot:     true,
			FirstName: "Testing",
			Username:  "TestBot",
		},
		BotClient: NewBotClient(f),
	}
}

func testMessage(text string) *gotgbot.Update {
	return &gotgbot.Update{
		UpdateId: rand.Int63(),
		Message: &gotgbot.Message{
			MessageId: rand.Int63(),
			Text:      text,
			// populate other fields as needed
		},
	}

}
