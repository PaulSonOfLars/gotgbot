package gotgbot

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

//go:generate go run ./scripts/generate

// Bot is the default Bot struct used to send and receive messages to the telegram API.
type Bot struct {
	// The bot's User info, as returned by Bot.GetMe. Populated when created through the NewBot method.
	User
	// The bot client to use to make requests
	BotClient
}

// BotOpts declares all optional parameters for the NewBot function.
type BotOpts struct {
	// HTTP client with any custom settings (eg proxy information) that might be necessary.
	Client http.Client
	// Disables the token validity check.
	// Useful when running in time-constrained environments where the startup time should be minimised, and where the
	// token can be assumed to be valid (eg lambdas).
	// Warning: Disabling the token check will mean that the Bot.User struct will no longer be populated.
	DisableTokenCheck bool
	// UseTestEnvironment defines whether this bot was created to run on telegram's test environment.
	// Enabling this uses a slightly different API path.
	// See https://core.telegram.org/bots/webapps#using-bots-in-the-test-environment for more details.
	UseTestEnvironment bool
	// Request opts to use for checking token validity with Bot.GetMe. Can be slow - a high timeout (eg 10s) is
	// recommended.
	RequestOpts *RequestOpts
	// Default opts to use for all requests, when no other request opts are specified.
	DefaultRequestOpts *RequestOpts
}

// NewBot returns a new Bot struct populated with the necessary defaults.
func NewBot(token string, opts *BotOpts) (*Bot, error) {
	botClient := &BaseBotClient{
		Token:              token,
		Client:             http.Client{},
		DefaultRequestOpts: nil,
	}

	// Large timeout on the initial GetMe request as this can sometimes be slow.
	getMeReqOpts := &RequestOpts{
		Timeout: 10 * time.Second,
		APIURL:  DefaultAPIURL,
	}

	checkTokenValidity := true
	if opts != nil {
		botClient.Client = opts.Client
		botClient.UseTestEnvironment = opts.UseTestEnvironment
		if opts.DefaultRequestOpts != nil {
			botClient.DefaultRequestOpts = opts.DefaultRequestOpts
		}
		if opts.RequestOpts != nil {
			getMeReqOpts = opts.RequestOpts
		}
		checkTokenValidity = !opts.DisableTokenCheck
	}

	b := Bot{
		BotClient: botClient,
	}

	if checkTokenValidity {
		// Get bot info. This serves two purposes:
		// 1. Check token is valid.
		// 2. Populate the bot struct "User" field.
		botUser, err := b.GetMe(&GetMeOpts{RequestOpts: getMeReqOpts})
		if err != nil {
			return nil, err
		}
		b.User = *botUser
	}

	return &b, nil
}

func (bot *Bot) UseMiddleware(mw func(client BotClient) BotClient) *Bot {
	bot.BotClient = mw(bot.BotClient)
	return bot
}

var ErrNilBotClient = errors.New("nil BotClient")

func (bot *Bot) Request(method string, params map[string]string, data map[string]NamedReader, opts *RequestOpts) (json.RawMessage, error) {
	if bot.BotClient == nil {
		return nil, ErrNilBotClient
	}

	ctx, cancel := bot.BotClient.TimeoutContext(opts)
	defer cancel()

	return bot.BotClient.RequestWithContext(ctx, method, params, data, opts)
}
