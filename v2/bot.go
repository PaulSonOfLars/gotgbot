package gotgbot

import (
	"net/http"
	"time"
)

type Bot struct {
	Token       string
	User        User
	APIURL      string
	Client      http.Client
	GetTimeout  time.Duration
	PostTimeout time.Duration
}
