package handlers

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type baseCommand struct {
	baseHandler
	Triggers     []rune
	AllowEdited  bool
	AllowChannel bool
	Command      string
}

type Command struct {
	baseCommand
	Response func(b ext.Bot, u *gotgbot.Update) error
}

type ArgsCommand struct {
	baseCommand
	Response func(b ext.Bot, u *gotgbot.Update, args []string) error
}

func NewCommand(command string, response func(b ext.Bot, u *gotgbot.Update) error) Command {
	cmd := strings.ToLower(command)
	return Command{
		baseCommand: baseCommand{
			baseHandler:  baseHandler{
				Name: cmd,
			},
			Triggers:     []rune("/"),
			AllowEdited:  false,
			AllowChannel: false,
			Command:      cmd,
		},
		Response: response,
	}
}

func NewArgsCommand(command string, response func(b ext.Bot, u *gotgbot.Update, args []string) error) ArgsCommand {
	cmd := strings.ToLower(command)
	return ArgsCommand{
		baseCommand: baseCommand{
			baseHandler:  baseHandler{
				Name: cmd,
			},
			Triggers:     []rune("/"),
			AllowEdited:  false,
			AllowChannel: false,
			Command:      cmd,
		},
		Response: response,
	}
}

func (h Command) HandleUpdate(u *gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.Response(d.Bot, u)
}

func (h ArgsCommand) HandleUpdate(u *gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.Response(d.Bot, u, strings.Fields(u.EffectiveMessage.Text)[1:])
}

// todo optimise if statements?
func (h baseCommand) CheckUpdate(u *gotgbot.Update) (bool, error) {
	if u.EffectiveMessage == nil || u.EffectiveMessage.Text == "" {
		return false, nil
	}
	// if no edits and message is edited
	if !h.AllowEdited && u.EditedMessage != nil {
		return false, nil
	}
	// if no channel and message is channel message
	if !h.AllowChannel && u.ChannelPost != nil {
		return false, nil
	}
	// if no channel, no edits, and post is edited
	if !h.AllowChannel && !h.AllowEdited && u.EditedChannelPost != nil {
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

	return cmd == h.Command, nil
}
