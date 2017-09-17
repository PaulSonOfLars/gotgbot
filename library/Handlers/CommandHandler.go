package Handlers

import (
	"strings"
	"bot/library"
	"bot/library/Ext"
)

type Command struct {
	command  string
	response func(b Ext.Bot, u library.Update)
}

func NewCommand(command string, response func(b Ext.Bot, u library.Update)) Command {
	h := Command{}
	h.command = command
	h.response = response
	return h
}

func (h Command) HandleUpdate(update library.Update, d library.Dispatcher) {
	h.response(d.Bot, update)
}

func (h Command) CheckUpdate(update library.Update) bool {
	return update.Message != nil && update.Message.Text != "" &&
		strings.Split(strings.Fields(update.Message.Text)[0], "@")[0] == "/" + h.command
}