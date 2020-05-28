package ext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const ApiUrl = "https://api.telegram.org/bot"

var DefaultTgBotGetter = BotGetter{
	Client: &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Millisecond * 1500,
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

type BotGetter struct {
	Client *http.Client
	ApiUrl string
}

type TelegramError struct {
	Code        int
	Description string
}

func (t *TelegramError) Error() string {
	return fmt.Sprintf("%d: %s", t.Code, t.Description)
}

type Getter interface {
	Get(bot Bot, method string, params url.Values) (json.RawMessage, error)
	Post(bot Bot, fileType string, method string, params url.Values, file io.Reader, filename string) (json.RawMessage, error)
}

func Get(bot Bot, method string, params url.Values) (json.RawMessage, error) {
	return DefaultTgBotGetter.Get(bot, method, params)
}

func Post(bot Bot, fileType string, method string, params url.Values, file io.Reader, filename string) (json.RawMessage, error) {
	return DefaultTgBotGetter.Post(bot, fileType, method, params, file, filename)
}

func (tbg *BotGetter) Get(bot Bot, method string, params url.Values) (json.RawMessage, error) {
	req, err := http.NewRequest("GET", tbg.ApiUrl+bot.Token+"/"+method, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to build GET request to %v", method)
	}
	req.URL.RawQuery = params.Encode()

	bot.Logger.Debugf("executing GET: %+v", req)
	resp, err := tbg.Client.Do(req)
	if err != nil {
		bot.Logger.Debugw("failed to execute GET request", "method", method, zap.Error(err))
		return nil, errors.Wrapf(err, "unable to execute GET request to %v", method)
	}
	defer resp.Body.Close()
	bot.Logger.Debugf("successful GET request: %+v", resp)

	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		bot.Logger.Debugw("failed to deserialize GET response body", "method", method, zap.Error(err))
		return nil, errors.Wrapf(err, "could not decode in GET %v call", method)
	}
	if !r.Ok {
		return nil, &TelegramError{
			Code:        r.ErrorCode,
			Description: r.Description,
		}
	}

	bot.Logger.Debugf("received result: %+v", r)
	bot.Logger.Debugf("result response: %v", string(r.Result))
	return r.Result, nil
}

func (tbg *BotGetter) Post(bot Bot, fileType string, method string, params url.Values, file io.Reader, filename string) (json.RawMessage, error) {
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
		bot.Logger.Debugw("failed to execute POST request", "method", method, zap.Error(err))
		return nil, errors.Wrapf(err, "unable to execute POST request to %v", method)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Content-Type", w.FormDataContentType())

	bot.Logger.Debugf("POST request with body: %+v", b)
	bot.Logger.Debugf("executing POST: %+v", req)
	resp, err := tbg.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bot.Logger.Debugf("successful POST request: %+v", resp)

	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		bot.Logger.Debugw("failed to deserialize POST response body", "method", method, zap.Error(err))
		return nil, errors.Wrapf(err, "could not decode in POST %v call", method)
	}
	bot.Logger.Debugf("received result: %+v", r)
	bot.Logger.Debugf("result response: %v", string(r.Result))
	return r.Result, nil
}
