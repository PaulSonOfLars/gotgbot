package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"go.uber.org/zap"

	"github.com/pkg/errors"
)

type Bot struct {
	Token             string
	Id                int
	FirstName         string
	UserName          string
	Logger            *zap.SugaredLogger `json:"-"`
	DisableWebPreview bool
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// GetMe gets the bot info
func (b Bot) GetMe() (*User, error) {
	v := url.Values{}

	r, err := Get(b, "getMe", v)
	if err != nil {
		return nil, errors.Wrapf(err, "could not getMe")
	}

	var u User
	return &u, json.Unmarshal(r.Result, &u)
}

// GetMyCommands gets the list of bot commands assigned to the bot.
func (b Bot) GetMyCommands() ([]BotCommand, error) {
	v := url.Values{}

	r, err := Get(b, "getMyCommands", v)
	if err != nil {
		return nil, errors.Wrapf(err, "could not getMyCommands")
	}

	var bc []BotCommand
	return bc, json.Unmarshal(r.Result, &bc)
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

	r, err := Get(b, "getUserProfilePhotos", v)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get user profile photos")
	}

	var userProfilePhotos UserProfilePhotos
	return &userProfilePhotos, json.Unmarshal(r.Result, &userProfilePhotos)
}

// GetFile Retrieve a file from the bot api
func (b Bot) GetFile(fileId string) (*File, error) {
	v := url.Values{}
	v.Add("file_id", fileId)

	r, err := Get(b, "getFile", v)
	if err != nil {
		return nil, errors.Wrapf(err, "could not complete getFile request")
	}

	var f File
	return &f, json.Unmarshal(r.Result, &f)
}
