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

func NewPrefixCommand(command string, prefixes []rune, response func(b ext.Bot, u *gotgbot.Update) error) Command {
	cmd := strings.ToLower(command)
	return Command{
		baseCommand: baseCommand{
			baseHandler: baseHandler{
				Name: cmd,
			},
			Triggers:     prefixes,
			AllowEdited:  false,
			AllowChannel: false,
			Command:      cmd,
		},
		Response: response,
	}
}

func NewCommand(command string, response func(b ext.Bot, u *gotgbot.Update) error) Command {
	return NewPrefixCommand(command, []rune("/"), response)
}

func NewPrefixArgsCommand(command string, prefixes []rune, response func(b ext.Bot, u *gotgbot.Update, args []string) error) ArgsCommand {
	cmd := strings.ToLower(command)
	return ArgsCommand{
		baseCommand: baseCommand{
			baseHandler: baseHandler{
				Name: cmd,
			},
			Triggers:     prefixes,
			AllowEdited:  false,
			AllowChannel: false,
			Command:      cmd,
		},
		Response: response,
	}
}

func NewArgsCommand(command string, response func(b ext.Bot, u *gotgbot.Update, args []string) error) ArgsCommand {
	return NewPrefixArgsCommand(command, []rune("/"), response)
}

func (h Command) HandleUpdate(u *gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.Response(*d.Bot, u)
}

func (h ArgsCommand) HandleUpdate(u *gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.Response(*d.Bot, u, strings.Fields(u.EffectiveMessage.Text)[1:])
}

func (h baseCommand) CheckUpdate(u *gotgbot.Update) (bool, error) {
	if u.EffectiveMessage == nil {
		return false, nil
	}

	var text string
	if u.EffectiveMessage.Text != "" {
		text = u.EffectiveMessage.Text
	} else if u.EffectiveMessage.Caption != "" {
		text = u.EffectiveMessage.Caption
	} else {
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

	rText := []rune(text)
	botName := strings.ToLower(u.EffectiveMessage.Bot.UserName)
	var cmd string
	for _, x := range h.Triggers {
		if rText[0] != x {
			continue
		}

		split := strings.Split(strings.ToLower(strings.Fields(text)[0]), "@")
		if len(split) > 1 && split[1] != botName {
			return false, nil
		}
		cmd = split[0][1:]
		break
	}
	if cmd == "" {
		return false, nil
	}

	if len(u.EffectiveMessage.Entities) != 0 && u.EffectiveMessage.Entities[0].Offset == 0 && u.EffectiveMessage.Entities[0].Type != "bot_command" {
		return false, nil
	}

	return cmd == h.Command, nil
}
