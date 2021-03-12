package ext

import "github.com/PaulSonOfLars/gotgbot/v2"

type Handler interface {
	// CheckUpdate checks whether the update should handled by this handler.
	CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool
	// HandleUpdate processes the update.
	HandleUpdate(b *gotgbot.Bot, ctx *Context) error
	// Name gets the handler name; used to differentiate handlers programmatically. Names should be unique.
	Name() string
}
