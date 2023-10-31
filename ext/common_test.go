package ext

import (
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

func (u *Updater) InjectUpdate(token string, upd gotgbot.Update) error {
	bData, ok := u.botMapping.getBot(token)
	if !ok {
		return ErrNotFound
	}

	return u.Dispatcher.(*Dispatcher).ProcessUpdate(bData.bot, &upd, nil)
}
