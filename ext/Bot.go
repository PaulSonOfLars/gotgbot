package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type Bot struct {
	Token             string
	Id                int
	FirstName         string
	UserName          string
	Logger            *logrus.Logger `json:"-"`
	DisableWebPreview bool
}

// GetMe get the bot info
func (b Bot) GetMe() (*User, error) {
	v := url.Values{}

	r, err := Get(b, "getMe", v)
	if err != nil {
		return nil, errors.Wrapf(err, "could not getMe")
	}
	if !r.Ok {
		return nil, errors.New("invalid getMe request")
	}

	var u User
	return &u, json.Unmarshal(r.Result, &u)
}

// GetUserProfilePhotos Retrieves a user's profile pictures
func (b Bot) GetUserProfilePhotos(userId int) (*UserProfilePhotos, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(userId))

	r, err := Get(b, "getUserProfilePhotos", v)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get user profile photos")
	}
	if !r.Ok {
		return nil, errors.New("invalid getUserProfilePhotos request")
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
	if !r.Ok {
		return nil, errors.New("invalid getFile request")
	}

	var f File
	return &f, json.Unmarshal(r.Result, &f)
}
