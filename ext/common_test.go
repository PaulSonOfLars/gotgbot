package ext

import (
	"errors"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type DummyHandler struct {
	F func(b *gotgbot.Bot, ctx *Context) error
	N string
}

func (d DummyHandler) CheckUpdate(b *gotgbot.Bot, ctx *Context) bool {
	return true
}

func (d DummyHandler) HandleUpdate(b *gotgbot.Bot, ctx *Context) error {
	return d.F(b, ctx)
}

func (d DummyHandler) Name() string {
	return "dummy" + d.N
}

var ErrBadDispatcher = errors.New("can only inject updates if the dispatcher is of type *Dispatcher")

func (u *Updater) InjectUpdate(token string, upd gotgbot.Update) error {
	bData, ok := u.botMapping.getBot(token)
	if !ok {
		return ErrNotFound
	}

	d, ok := u.Dispatcher.(*Dispatcher)
	if !ok {
		return ErrBadDispatcher
	}
	return d.ProcessUpdate(bData.bot, &upd, nil)
}
