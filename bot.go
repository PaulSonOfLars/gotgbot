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

func NewBot(token string) (*Bot, error) {
	b := Bot{
		Token:  token,
		User:   User{},
		Client: http.Client{},
	}

	u, err := b.GetMe()
	if err != nil {
		return nil, err
	}

	b.User = *u
	return &b, nil
}
