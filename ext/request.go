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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

type response struct {
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
	return t.Description
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
		bot.Logger.Debugw("failed to execute GET request",
			zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: method,
			},
			zap.Error(err))
		return nil, errors.Wrapf(err, "unable to execute GET request to %v", method)
	}
	defer resp.Body.Close()
	bot.Logger.Debugf("successful GET request: %+v", resp)

	var r response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		bot.Logger.Debugw("failed to deserialize GET response body",
			zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: method,
			},
			zap.Error(err))
		return nil, errors.Wrapf(err, "could not decode in GET %v call", method)
	}
	if !r.Ok {
		bot.Logger.Debugw("error from GET",
			zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: method,
			},
			zapcore.Field{
				Key:     "error_code",
				Type:    zapcore.Int64Type,
				Integer: int64(r.ErrorCode),
			},
			zapcore.Field{
				Key:    "description",
				Type:   zapcore.StringType,
				String: r.Description,
			},
		)
		return nil, &TelegramError{
			Code:        r.ErrorCode,
			Description: r.Description,
		}
	}

	bot.Logger.Debugw("obtained GET result",
		zapcore.Field{
			Key:    "method",
			Type:   zapcore.StringType,
			String: method,
		},
		zapcore.Field{
			Key:    "result",
			Type:   zapcore.StringType,
			String: string(r.Result),
		},
	)

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
		bot.Logger.Debugw("failed to execute POST request",
			zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: method,
			},
			zap.Error(err))
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

	var r response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		bot.Logger.Debugw("failed to deserialize POST response body",
			zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: method,
			},
			zap.Error(err))
		return nil, errors.Wrapf(err, "could not decode in POST %v call", method)
	}
	if !r.Ok {
		bot.Logger.Debugw("error from POST",
			zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: method,
			},
			zapcore.Field{
				Key:     "error_code",
				Type:    zapcore.Int64Type,
				Integer: int64(r.ErrorCode),
			},
			zapcore.Field{
				Key:    "description",
				Type:   zapcore.StringType,
				String: r.Description,
			},
		)
		return nil, &TelegramError{
			Code:        r.ErrorCode,
			Description: r.Description,
		}
	}

	bot.Logger.Debugw("obtained POST result",
		zapcore.Field{
			Key:    "method",
			Type:   zapcore.StringType,
			String: method,
		},
		zapcore.Field{
			Key:    "result",
			Type:   zapcore.StringType,
			String: string(r.Result),
		},
	)

	return r.Result, nil
}
