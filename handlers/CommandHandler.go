package handlers

import (
	"strings"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type baseCommand struct {
	command  string
}

type Command struct {
	baseCommand
	response func(b ext.Bot, u gotgbot.Update)
}

type ArgsCommand struct {
	baseCommand
	response func(b ext.Bot, u gotgbot.Update, args []string)
}

func NewCommand(command string, response func(b ext.Bot, u gotgbot.Update)) Command {
	h := Command{}
	h.command = command
	h.response = response
	return h
}

func NewArgsCommand(command string, response func(b ext.Bot, u gotgbot.Update, args []string)) ArgsCommand {
	h := ArgsCommand{}
	h.command = command
	h.response = response
	return h
}

func (h Command) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher) {
	h.response(d.Bot, update)
}

func (h ArgsCommand) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher){
	h.response(d.Bot, update, strings.Fields(update.EffectiveMessage.Text)[1:])
}

func (h baseCommand) CheckUpdate(update gotgbot.Update) (bool, error) {
	return update.Message != nil && update.Message.Text != "" &&
		strings.Split(strings.Fields(update.Message.Text)[0], "@")[0] == "/"+h.command, nil
}
