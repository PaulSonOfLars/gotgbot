package handlers

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type baseCommand struct {
	Triggers    []rune
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
			Triggers: []rune("/"),
			command: strings.ToLower(command),
		},
		response: response,
	}
}

func NewArgsCommand(command string, response func(b ext.Bot, u gotgbot.Update, args []string) error) ArgsCommand {
	return ArgsCommand{
		baseCommand: baseCommand{
			Triggers: []rune("/"),
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
	if u.EffectiveMessage == nil || u.EffectiveMessage.Text == "" {
		return false, nil
	}
	if !h.AllowEdited && u.EditedMessage != nil {
		return false, nil
	}

	var cmd string
	for _, x := range h.Triggers {
		if []rune(u.EffectiveMessage.Text)[0] == x {
			stuff := strings.Split(strings.ToLower(strings.Fields(u.EffectiveMessage.Text)[0]), "@")
			// todo: remove repeated ToLower of username
			if len(stuff) > 1 && stuff[1] != strings.ToLower(u.EffectiveMessage.Bot.UserName) {
				return false, nil
			}
			cmd = strings.ToLower(stuff[0])[1:]
			break
		}
	}
	if cmd == "" {
		return false, nil
	}

	if len(u.EffectiveMessage.Entities) != 0 && u.EffectiveMessage.Entities[0].Offset == 0 && u.EffectiveMessage.Entities[0].Type != "bot_command" {
		return false, nil
	}

	return cmd == h.command, nil
}
