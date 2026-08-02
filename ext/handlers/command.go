package handlers

import (
	"strings"
	"unicode/utf8"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Command is the go-to handler for setting up Commands in your bot. By default, it will use telegram-native commands
// that start with a forward-slash (/), but it can be customised to react to any message starting with a character.
//
// For example, a command handler on "help" with the triggers []rune("/!,") would trigger for "/help", "!help", or ",help".
type Command struct {
	Triggers     []rune
	AllowEdited  bool
	AllowChannel bool
	AllowBot     bool
	Command      string // set to a lowercase value for case-insensitivity
	Response     Response
}

// NewCommand creates a new case-insensitive command.
// By default, commands do not work on edited messages, or channel posts. These can be enabled by setting the
// AllowEdited and AllowChannel fields respectively.
func NewCommand(c string, r Response) Command {
	return Command{
		Triggers:     []rune{'/'},
		AllowEdited:  false,
		AllowChannel: false,
		AllowBot:     false,
		Command:      strings.ToLower(c),
		Response:     r,
	}
}

// SetAllowEdited Enables edited messages for this handler.
func (c Command) SetAllowEdited(allow bool) Command {
	c.AllowEdited = allow
	return c
}

// SetAllowChannel Enables channel messages for this handler.
func (c Command) SetAllowChannel(allow bool) Command {
	c.AllowChannel = allow
	return c
}

// SetAllowBot Enables bot messages for this handler.
func (c Command) SetAllowBot(allow bool) Command {
	c.AllowBot = allow
	return c
}

// SetTriggers sets the list of triggers to be used with this command.
func (c Command) SetTriggers(triggers []rune) Command {
	c.Triggers = triggers
	return c
}

func (c Command) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.EffectiveSender != nil && ctx.EffectiveSender.IsBot() && !c.AllowBot {
		return false
	}

	if ctx.Message != nil {
		return c.checkMessage(b, ctx.Message)
	}

	// if no edits and message is edited
	if c.AllowEdited && ctx.EditedMessage != nil {
		return c.checkMessage(b, ctx.EditedMessage)
	}
	// if no channel and message is channel message
	if c.AllowChannel && ctx.ChannelPost != nil {
		return c.checkMessage(b, ctx.ChannelPost)
	}
	// if no channel, no edits, and post is edited
	if c.AllowChannel && c.AllowEdited && ctx.EditedChannelPost != nil {
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
	if rich := msg.RichMessage; rich != nil {
		cmd, _ := c.extractRichCommand(rich)
		return cmd != "" && cmd == c.Command
	}

	text := msg.GetText()
	if text == "" {
		return false
	}

	ents := msg.GetEntities()
	if len(ents) != 0 && ents[0].Offset == 0 && ents[0].Type != gotgbot.MessageEntityTypeBotCommand {
		return false
	}

	// Find the end of the first field without scanning the whole string
	end := strings.IndexAny(text, " \t\n\v\f\r")
	var firstField string
	if end == -1 {
		firstField = text
	} else {
		firstField = text[:end]
	}

	split := strings.SplitN(strings.ToLower(firstField), "@", 2)

	// If the command targets a specific bot, ensure it's this one
	if len(split) > 1 && split[1] != strings.ToLower(b.User.Username) {
		return false
	}

	command := c.extractCommand(split[0])
	return command != "" && command == c.Command
}

func (c Command) extractCommand(command string) string {
	first, size := utf8.DecodeRuneInString(command)
	for _, t := range c.Triggers {
		if first == t {
			return command[size:]
		}
	}
	return ""
}

func (c Command) extractRichCommand(rm *gotgbot.RichMessage) (command string, ok bool) {
	if rm == nil || len(rm.Blocks) == 0 {
		return "", false
	}

	paragraph, ok := rm.Blocks[0].(gotgbot.RichBlockParagraph)
	if !ok {
		return "", false
	}

	wholeText := strings.TrimLeft(gotgbot.RichBlockContent(paragraph), " \n\t")
	return c.matchRichNode(paragraph.Text, wholeText)
}

func (c Command) matchRichNode(rt gotgbot.RichText, wholeText string) (string, bool) {
	switch v := rt.(type) {
	case gotgbot.RichTextArray:
		if len(v) == 0 {
			return "", false
		}
		return c.matchRichNode(v[0], wholeText)

	case gotgbot.RichTextString:
		cmd := c.extractCommand(string(v))
		return cmd, cmd != ""

	case gotgbot.RichTextBotCommand:
		cmdString := gotgbot.RichTextContent(v) // rendered token, e.g. "/start@mybot"

		cmd := c.extractCommand(cmdString)
		if cmd == "" {
			return "", false
		}

		return v.BotCommand, strings.HasPrefix(wholeText, cmdString)
	}
	return "", false
}
