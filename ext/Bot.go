package ext

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type Bot struct {
	Token     string
	Id        int
	FirstName string
	UserName  string
}

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

func (b Bot) AnswerCallbackQuery(callbackQueryId string) (bool, error) {
	v := url.Values{}
	v.Add("callback_query_id", callbackQueryId)

	return b.boolSender("answerCallbackQuery", v)
}

func (b Bot) AnswerCallbackQueryText(callbackQueryId string, text string, alert bool) (bool, error) {
	v := url.Values{}
	v.Add("callback_query_id", callbackQueryId)
	v.Add("text", text)
	v.Add("alert", strconv.FormatBool(alert))

	return b.boolSender("answerCallbackQuery", v)
}

func (b Bot) AnswerCallbackQueryURL(callbackQueryId string, URL string) (bool, error) {
	v := url.Values{}
	v.Add("callback_query_id", callbackQueryId)
	v.Add("url", URL)

	return b.boolSender("answerCallbackQuery", v)
}
