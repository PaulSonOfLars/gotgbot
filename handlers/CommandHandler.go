package handlers

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type baseCommand struct {
	command string
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
	return Command{
		baseCommand: baseCommand{
			command: strings.ToLower(command),
		},
		response: response,
	}
}

func NewArgsCommand(command string, response func(b ext.Bot, u gotgbot.Update, args []string)) ArgsCommand {
	return ArgsCommand{
		baseCommand: baseCommand{
			command: strings.ToLower(command),
		},
		response: response,
	}
}

func (h Command) HandleUpdate(u gotgbot.Update, d gotgbot.Dispatcher) {
	h.response(d.Bot, u)
}

func (h ArgsCommand) HandleUpdate(u gotgbot.Update, d gotgbot.Dispatcher) {
	h.response(d.Bot, u, strings.Fields(u.EffectiveMessage.Text)[1:])
}

func (h baseCommand) CheckUpdate(u gotgbot.Update) (bool, error) {
	if u.Message != nil && u.Message.Text != "" &&
		len(u.Message.Entities) > 0 && u.Message.Entities[0].Type == "bot_command" {
		cmd := strings.Split(strings.Fields(u.Message.Text)[0], "@")
		return cmd[0] == "/"+h.command && (len(cmd) <= 1 || cmd[1] == u.Message.Bot.UserName), nil
	}
	return false, nil
}
