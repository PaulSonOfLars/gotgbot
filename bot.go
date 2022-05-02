package gotgbot

import (
	"net/http"
	"time"
)

//go:generate go run ./scripts/generate

// Bot is the core Bot object used to send and receive messages.
type Bot struct {
	// The bot's User info, as returned by Bot.GetMe. Populated when created through the NewBot method.
	User
	// Token stores the bot's secret token obtained from t.me/BotFather, and used to interact with telegram's API.
	Token string
	// Client is the HTTP Client used for all HTTP requests made for this bot.
	Client http.Client
	// The default request opts for this bot instance.
	DefaultRequestOpts *RequestOpts
}

// BotOpts declares all optional parameters for the NewBot function.
type BotOpts struct {
	// HTTP client with any custom settings (eg proxy information) that might be necessary.
	Client http.Client
	// Default request opts to use when no other request opts are specified.
	DefaultRequestOpts *RequestOpts
	// Request opts to use for checking token validity with Bot.GetMe. Can be slow - a high timeout (eg 10s) is
	// recommended.
	RequestOpts *RequestOpts
}

// NewBot returns a new Bot struct populated with the necessary defaults.
func NewBot(token string, opts *BotOpts) (*Bot, error) {
	// Barebones bot - token not verified yet, no settings set
	b := Bot{Token: token}

	// Large timeout on the initial GetMe request as this can sometimes be slow.
	getMeReqOpts := &RequestOpts{
		Timeout: 10 * time.Second,
		APIURL:  DefaultAPIURL,
	}

	if opts != nil {
		b.Client = opts.Client
		if opts.DefaultRequestOpts != nil {
			b.DefaultRequestOpts = opts.DefaultRequestOpts
		}
		if opts.RequestOpts != nil {
			getMeReqOpts = opts.RequestOpts
		}
	}

	// Get bot info. This serves two purposes:
	// 1. Check token is valid.
	// 2. Populate the bot struct "User" field.
	botUser, err := b.GetMe(&GetMeOpts{RequestOpts: getMeReqOpts})
	if err != nil {
		return nil, err
	}

	b.User = *botUser
	return &b, nil
}
