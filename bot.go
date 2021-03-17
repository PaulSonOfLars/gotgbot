package gotgbot

import (
	"net/http"
	"time"
)

//go:generate go run ./scripts/generate

type Bot struct {
	Token       string
	User        User
	APIURL      string
	Client      http.Client
	GetTimeout  time.Duration
	PostTimeout time.Duration
}

// BotOpts Declares all optional parameters for the NewBot function.
type BotOpts struct {
	APIURL      string
	Client      http.Client
	GetTimeout  time.Duration
	PostTimeout time.Duration
}

func NewBot(token string, opts *BotOpts) (*Bot, error) {
	b := Bot{
		Token:      token,
		GetTimeout: time.Second * 10,// 10 seconds timeout for initial GetMe request, which can be slow.
	}

	getTimeout := DefaultGetTimeout
	postTimeout := DefaultPostTimeout
	if opts != nil {
		b.Client = opts.Client
		b.APIURL = opts.APIURL

		getTimeout = opts.GetTimeout
		postTimeout = opts.PostTimeout
	}

	u, err := b.GetMe()
	if err != nil {
		return nil, err
	}

	b.User = *u
	b.GetTimeout = getTimeout
	b.PostTimeout = postTimeout
	return &b, nil
}
