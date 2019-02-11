package ext

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const ApiUrl = "https://api.telegram.org/bot"

var DefaultTgBotGetter = TgBotGetter{
	Client: &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 5,
	},
	ApiUrl: ApiUrl,
}

type Response struct {
	Ok          bool
	Result      json.RawMessage
	ErrorCode   int `json:"error_code"`
	Description string
	Parameters  json.RawMessage
}

type TgBotGetter struct {
	Client *http.Client
	ApiUrl string
}

type TgBotGetterInterface interface {
	Get(bot Bot, method string, params url.Values) (*Response, error)
	Post(bot Bot, fileType string, method string, params url.Values, file io.Reader, filename string) (*Response, error)
}

func Get(bot Bot, method string, params url.Values) (*Response, error) {
	return DefaultTgBotGetter.Get(bot, method, params)
}

func (tbg *TgBotGetter) Get(bot Bot, method string, params url.Values) (*Response, error) {
	req, err := http.NewRequest("GET", tbg.ApiUrl+bot.Token+"/"+method, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to build get request to %v", method)
	}
	req.URL.RawQuery = params.Encode()

	resp, err := tbg.Client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to execute get request to %v", method)
	}
	defer resp.Body.Close()

	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrapf(err, "could not decode in Get %v call", method)
	}
	return &r, nil
}

func (tbg *TgBotGetter) Post(bot Bot, fileType string, method string, params url.Values, file io.Reader, filename string) (*Response, error) {
	if filename == "" {
		filename = "unnamed_file"
	}
	b := bytes.Buffer{}
	w := multipart.NewWriter(&b)
	part, err := w.CreateFormFile(fileType, filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", tbg.ApiUrl+bot.Token+"/"+method, &b)
	if err != nil {
		logrus.WithError(err).Errorf("failed to send to %v func", method)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := tbg.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
