package handlers

import (
	"strings"
	"unicode/utf8"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Command struct {
	Triggers     []rune
	AllowEdited  bool
	AllowChannel bool
	Command      string // should be lowercase for case-insensitivity
	Response     Response
}

func NewCommand(c string, r Response) Command {
	return Command{
		Triggers:     []rune{'/'},
		AllowEdited:  false,
		AllowChannel: false,
		Command:      strings.ToLower(c),
		Response:     r,
	}
}

func (c Command) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.Message != nil {
		if ctx.Message.Text == "" && ctx.Message.Caption == "" {
			return false
		}
		return c.checkMessage(b, ctx.Message)
	}

	// if no edits and message is edited
	if c.AllowEdited && ctx.EditedMessage != nil {
		if ctx.EditedMessage.Text == "" && ctx.EditedMessage.Caption == "" {
			return false
		}
		return c.checkMessage(b, ctx.EditedMessage)
	}
	// if no channel and message is channel message
	if c.AllowChannel && ctx.ChannelPost != nil {
		if ctx.ChannelPost.Text == "" && ctx.ChannelPost.Caption == "" {
			return false
		}
		return c.checkMessage(b, ctx.ChannelPost)
	}
	// if no channel, no edits, and post is edited
	if c.AllowChannel && c.AllowEdited && ctx.EditedChannelPost != nil {
		if ctx.EditedChannelPost.Text == "" && ctx.EditedChannelPost.Caption == "" {
			return false
		}
		return c.checkMessage(b, ctx.EditedChannelPost)
	}

	return false
}

func (c Command) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.Response(b, ctx)
}

func (c Command) Name() string {
	return "command_" + c.Command
}

func (c Command) checkMessage(b *gotgbot.Bot, msg *gotgbot.Message) bool {
	text := msg.Text
	if msg.Caption != "" {
		text = msg.Caption
	}

	var cmd string
	for _, t := range c.Triggers {
		if r, _ := utf8.DecodeRuneInString(text); r != t {
			continue
		}

		split := strings.Split(strings.ToLower(strings.Fields(text)[0]), "@")
		if len(split) > 1 && split[1] != strings.ToLower(b.User.Username) {
			return false
		}
		cmd = split[0][1:]
		break
	}
	if cmd == "" {
		return false
	}

	if len(msg.Entities) != 0 && msg.Entities[0].Offset == 0 && msg.Entities[0].Type != "bot_command" {
		return false
	}

	return cmd == c.Command
}
