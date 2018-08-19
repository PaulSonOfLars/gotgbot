package handlers

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type baseCommand struct {
	AllowEdited bool
	command     string
}

type Command struct {
	baseCommand
	response func(b ext.Bot, u gotgbot.Update) error
}

type ArgsCommand struct {
	baseCommand
	response func(b ext.Bot, u gotgbot.Update, args []string) error
}

func NewCommand(command string, response func(b ext.Bot, u gotgbot.Update) error) Command {
	return Command{
		baseCommand: baseCommand{
			command: strings.ToLower(command),
		},
		response: response,
	}
}

func NewArgsCommand(command string, response func(b ext.Bot, u gotgbot.Update, args []string) error) ArgsCommand {
	return ArgsCommand{
		baseCommand: baseCommand{
			command: strings.ToLower(command),
		},
		response: response,
	}
}

func (h Command) HandleUpdate(u gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.response(d.Bot, u)
}

func (h ArgsCommand) HandleUpdate(u gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.response(d.Bot, u, strings.Fields(u.EffectiveMessage.Text)[1:])
}

func (h baseCommand) CheckUpdate(u gotgbot.Update) (bool, error) {
	if ((u.Message != nil && u.Message.Text != "") || (h.AllowEdited && u.EditedMessage != nil && u.EditedMessage.Text != "")) &&
		len(u.EffectiveMessage.Entities) > 0 && u.EffectiveMessage.Entities[0].Type == "bot_command" {
		cmd := strings.Split(strings.Fields(strings.ToLower(u.EffectiveMessage.Text))[0], "@")
		// TODO: remove repeated tolower of bot username
		return cmd[0] == "/"+h.command && (len(cmd) <= 1 || cmd[1] == strings.ToLower(u.EffectiveMessage.Bot.UserName)), nil
	}
	return false, nil
}
