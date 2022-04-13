package gotgbot

import (
	"net/http"
	"time"
)

//go:generate go run ./scripts/generate

// Bot is the core Bot object used to send and receive messages.
type Bot struct {
	User
	Token          string
	APIURL         string
	Client         http.Client
	RequestTimeout time.Duration
}

// BotOpts declares all optional parameters for the NewBot function.
type BotOpts struct {
	APIURL         string
	Client         http.Client
	RequestTimeout time.Duration
}

// NewBot returns a new Bot struct populated with the necessary defaults.
func NewBot(token string, opts *BotOpts) (*Bot, error) {
	b := Bot{
		Token:          token,
		RequestTimeout: time.Second * 10, // 10 seconds timeout for initial GetMe request, which can be slow.
		APIURL:         DefaultAPIURL,
	}

	timeout := DefaultTimeout
	if opts != nil {
		b.Client = opts.Client
		b.APIURL = opts.APIURL

		timeout = opts.RequestTimeout
	}

	u, err := b.GetMe()
	if err != nil {
		return nil, err
	}

	b.User = *u
	b.RequestTimeout = timeout
	return &b, nil
}
