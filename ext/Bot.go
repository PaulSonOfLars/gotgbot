package ext

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/errors"
)

type Bot struct {
	Token     string
	Id        int
	FirstName string
	UserName  string

	Logger            *zap.SugaredLogger `json:"-"`
	DisableWebPreview bool
	Requester
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type Webhook struct {
	Serve     string // base url to where you listen
	ServePath string // path you listen to
	ServePort int    // port you listen on
	URL       string // where you set the webhook to send to
	// CertPath       string   // TODO
	IPAddress          string   // where you set the webhook to send to
	MaxConnections     int      // max connections; max 100, default 40
	AllowedUpdates     []string // which updates to allow
	DropPendingUpdates bool     // should start from scratch on the new updates
}

func (w Webhook) GetListenUrl() string {
	if w.Serve == "" {
		w.Serve = "0.0.0.0"
	}
	if w.ServePort == 0 {
		w.ServePort = 443
	}
	return fmt.Sprintf("%s:%d", w.Serve, w.ServePort)
}

func NewBot(l *zap.Logger, token string) (*Bot, error) {
	user, err := Bot{
		Token:  token,
		Logger: l.Sugar(),
		// getMe often times out, so add a large 5 second timeout for lots of leeway
		Requester: BaseRequester{Client: http.Client{Timeout: time.Second * 5}},
	}.GetMe()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create new bot")
	}
	return &Bot{
		Token:     token,
		Id:        user.Id,
		FirstName: user.FirstName,
		UserName:  user.Username,
		Logger:    l.Sugar(),
		Requester: DefaultTgBotRequester,
	}, nil
}

// Get -> convenience function to execute commands against the TG API.
func (b Bot) Get(method string, params url.Values) (json.RawMessage, error) {
	return b.Requester.Get(b.Logger, b.Token, method, params)
}

// Post -> convenience function to execute commands against the TG API.
// the "data" map is used to fill out the multipart form data to send to TG.
func (b Bot) Post(method string, params url.Values, data map[string]PostFile) (json.RawMessage, error) {
	return b.Requester.Post(b.Logger, b.Token, method, params, data)
}

// GetMe gets the bot info
func (b Bot) GetMe() (*User, error) {
	v := url.Values{}

	r, err := b.Get("getMe", v)
	if err != nil {
		return nil, err
	}

	var u User
	return &u, json.Unmarshal(r, &u)
}

func (b Bot) LogOut() (bool, error) {
	return b.boolSender("logOut", nil)
}

func (b Bot) Close() (bool, error) {
	return b.boolSender("close", nil)
}

// GetMyCommands gets the list of bot commands assigned to the bot.
func (b Bot) GetMyCommands() ([]BotCommand, error) {
	v := url.Values{}

	r, err := b.Get("getMyCommands", v)
	if err != nil {
		return nil, err
	}

	var bc []BotCommand
	return bc, json.Unmarshal(r, &bc)
}

// SetMyCommands gets the list of bot commands assigned to the bot.
func (b Bot) SetMyCommands(cmds []BotCommand) (bool, error) {
	if cmds == nil {
		cmds = []BotCommand{}
	}

	v := url.Values{}
	cmdBytes, err := json.Marshal(cmds)
	if err != nil {
		return false, errors.Wrapf(err, "failed to marshal commands")
	}
	v.Add("commands", string(cmdBytes))

	return b.boolSender("setMyCommands", v)
}

// GetUserProfilePhotos Retrieves a user's profile pictures
func (b Bot) GetUserProfilePhotos(userId int, offset int, limit int) (*UserProfilePhotos, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(userId))

	if offset != 0 {
		v.Add("offset", strconv.Itoa(offset))
	}

	if limit != 0 {
		v.Add("limit", strconv.Itoa(limit))
	}

	r, err := b.Get("getUserProfilePhotos", v)
	if err != nil {
		return nil, err
	}

	var userProfilePhotos UserProfilePhotos
	return &userProfilePhotos, json.Unmarshal(r, &userProfilePhotos)
}

// GetFile Retrieve a file from the bot api
func (b Bot) GetFile(fileId string) (*File, error) {
	v := url.Values{}
	v.Add("file_id", fileId)

	r, err := b.Get("getFile", v)
	if err != nil {
		return nil, err
	}

	var f File
	return &f, json.Unmarshal(r, &f)
}

type WebhookInfo struct {
	URL                  string   `json:"url"`
	HasCustomCertificate bool     `json:"has_custom_certificate"`
	PendingUpdateCount   int      `json:"pending_update_count"`
	IPAddress            string   `json:"ip_address"`
	LastErrorDate        int      `json:"last_error_date"`
	LastErrorMessage     int      `json:"last_error_message"`
	MaxConnections       int      `json:"max_connections"`
	AllowedUpdates       []string `json:"allowed_updates"`
}

// GetWebhookInfo Get webhook info from telegram servers
func (b Bot) GetWebhookInfo() (*WebhookInfo, error) {
	r, err := b.Get("getWebhookInfo", nil)
	if err != nil {
		return nil, err
	}

	var wh WebhookInfo
	return &wh, json.Unmarshal(r, &wh)
}

// SetWebhook Set the webhook url for telegram to contact with updates
func (b Bot) SetWebhook(path string, webhook Webhook) (bool, error) {
	allowedUpdates := webhook.AllowedUpdates
	if allowedUpdates == nil {
		allowedUpdates = []string{}
	}
	allowed, err := json.Marshal(allowedUpdates)
	if err != nil {
		return false, errors.Wrap(err, "cannot marshal allowedUpdates")
	}

	v := url.Values{}

	v.Add("url", strings.TrimSuffix(webhook.URL, "/")+"/"+strings.TrimPrefix(path, "/"))
	// v.Add("certificate", ) // todo: add certificate support
	v.Add("ip_address", webhook.IPAddress)
	v.Add("max_connections", strconv.Itoa(webhook.MaxConnections))
	v.Add("allowed_updates", string(allowed))
	v.Add("drop_pending_updates", strconv.FormatBool(webhook.DropPendingUpdates))

	r, err := b.Get("setWebhook", v)
	if err != nil {
		return false, err
	}

	var bb bool
	return bb, json.Unmarshal(r, &bb)
}

func (b Bot) DeleteWebhook() (bool, error) {
	return b.DeleteWebhookDrop(false)
}

func (b Bot) DeleteWebhookDrop(dropUpdates bool) (bool, error) {
	v := url.Values{}
	v.Add("drop_pending_updates", strconv.FormatBool(dropUpdates))

	r, err := b.Get("deleteWebhook", v)
	if err != nil {
		return false, err
	}
	var bb bool
	return bb, json.Unmarshal(r, &bb)
}
