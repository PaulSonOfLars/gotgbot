package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// Define middleware BotClient.
type sendWithoutReplyBotClient struct {
	// Inline existing client to call, allowing us to chain middlewares.
	// Inlining also avoids us having to redefine helper methods part of the interface.
	gotgbot.BotClient
}

// Define wrapper around existing RequestWithContext method.
// Note: this is the only method that needs redefining.
func (b sendWithoutReplyBotClient) RequestWithContext(ctx context.Context, token string, method string, params map[string]string, data map[string]gotgbot.NamedReader, opts *gotgbot.RequestOpts) (json.RawMessage, error) {
	// For all sendable methods, we want to allow sending if the message has been deleted.
	// So, we edit the params to allow for that.
	// We also log this, for the sake of the example. :)
	if strings.HasPrefix(method, "send") || method == "copyMessage" {
		log.Println("Applying middleware to", method)
		params["allow_sending_without_reply"] = "true"
	}

	// Call the next bot client instance in the middleware chain.
	val, err := b.BotClient.RequestWithContext(ctx, token, method, params, data, opts)
	if err != nil {
		// Middlewares can also be used to increase error visibility, in case they aren't logged elsewhere.
		log.Println("warning, got an error:", err)
	}
	return val, err
}

func newSendWithoutReplyClient() sendWithoutReplyBotClient {
	return sendWithoutReplyBotClient{
		BotClient: &gotgbot.BaseBotClient{
			Client:             http.Client{},
			UseTestEnvironment: false,
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  gotgbot.DefaultAPIURL,
			},
		},
	}
}
