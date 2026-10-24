package gotgbot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:generate go run ./scripts/generate

var (
	ErrNilBotClient       = errors.New("nil BotClient")
	ErrInvalidTokenFormat = errors.New("invalid token format")
)

// Bot is the default Bot struct used to send and receive messages to the telegram API.
type Bot struct {
	// Token stores the bot's secret token obtained from t.me/BotFather, and used to interact with telegram's API.
	Token string

	// The bot's User info, as returned by Bot.GetMe. Populated when created through the NewBot method.
	User
	// The bot client to use to make requests
	BotClient

	// Any extra usernames that might be useful in command handler resolution.
	additionalUsernames map[string]struct{}
	usernamesMu         sync.RWMutex
}

// BotOpts declares all optional parameters for the NewBot function.
type BotOpts struct {
	// BotClient allows for passing in custom configurations of BotClient, such as handling extra errors or providing
	// metrics.
	BotClient BotClient
	// DisableTokenCheck can be used to disable the token validity check.
	// Useful when running in time-constrained environments where the startup time should be minimised, and where the
	// token can be assumed to be valid (eg lambdas).
	// Warning: Disabling the token check will mean that the Bot.User struct will no longer be populated.
	DisableTokenCheck bool
	// Request opts to use for checking token validity with Bot.GetMe. Can be slow - a high timeout (eg 10s) is
	// recommended.
	RequestOpts *RequestOpts
}

// NewBot returns a new Bot struct populated with the necessary defaults.
func NewBot(token string, opts *BotOpts) (*Bot, error) {
	botClient := BotClient(&BaseBotClient{
		Client:             http.Client{},
		UseTestEnvironment: false,
		DefaultRequestOpts: nil,
	})

	// Large timeout on the initial GetMe request as this can sometimes be slow.
	getMeReqOpts := &RequestOpts{
		Timeout: 10 * time.Second,
	}

	checkTokenValidity := true
	if opts != nil {
		if opts.BotClient != nil {
			botClient = opts.BotClient
		}

		if opts.RequestOpts != nil {
			getMeReqOpts = opts.RequestOpts
		}
		checkTokenValidity = !opts.DisableTokenCheck
	}

	b := Bot{
		Token:     token,
		BotClient: botClient,
	}

	if checkTokenValidity {
		// Get bot info. This serves two purposes:
		// 1. Check token is valid.
		// 2. Populate the bot struct "User" field.
		botUser, err := b.GetMe(&GetMeOpts{RequestOpts: getMeReqOpts})
		if err != nil {
			return nil, fmt.Errorf("failed to check bot token: %w", err)
		}
		b.User = *botUser
	} else {
		// If token checks are disabled, we populate the bot's ID from the token.
		split := strings.Split(token, ":")
		if len(split) != 2 {
			return nil, fmt.Errorf("%w: expected '123:abcd', got %s", ErrInvalidTokenFormat, token)
		}

		id, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bot ID from token: %w", err)
		}
		b.User = User{
			Id:    id,
			IsBot: true,
			// We mark these fields as missing so we can know why they're not available
			FirstName: "<missing>",
			Username:  "<missing>",
		}
	}

	return &b, nil
}

// UseMiddleware allows you to wrap the existing bot client to enhance functionality
//
// Deprecated: Instead of using middlewares, consider implementing the BotClient interface.
func (bot *Bot) UseMiddleware(mw func(client BotClient) BotClient) *Bot {
	bot.BotClient = mw(bot.BotClient)
	return bot
}

func (bot *Bot) Request(method string, params map[string]any, opts *RequestOpts) (json.RawMessage, error) {
	return bot.RequestWithContext(context.Background(), method, params, opts)
}

func (bot *Bot) RequestWithContext(ctx context.Context, method string, params map[string]any, opts *RequestOpts) (json.RawMessage, error) {
	if bot.BotClient == nil {
		return nil, ErrNilBotClient
	}

	return bot.BotClient.RequestWithContext(ctx, bot.Token, method, params, opts)
}

// UpdateAdditionalUsernames replaces the usernames currently set for this bot.
// This is relevant for any features which rely on Bot.MatchesUsername, such as command handlers.
// Safe to call concurrently.
func (bot *Bot) UpdateAdditionalUsernames(us []string) {
	newUsernames := map[string]struct{}{}
	for _, u := range us {
		clean := strings.ToLower(strings.TrimPrefix(u, "@"))
		if clean == "" {
			continue
		}

		newUsernames[clean] = struct{}{}
	}

	bot.usernamesMu.Lock()
	defer bot.usernamesMu.Unlock()
	bot.additionalUsernames = newUsernames
}

// MatchesUsername assumes u does not include the @ prefix.
// Safe to call concurrently.
func (bot *Bot) MatchesUsername(u string) bool {
	if u == "" {
		return false
	}

	if strings.EqualFold(u, bot.User.Username) {
		return true
	}

	bot.usernamesMu.RLock()
	defer bot.usernamesMu.RUnlock()

	_, ok := bot.additionalUsernames[strings.ToLower(u)]
	return ok
}
