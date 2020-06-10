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
	"go.uber.org/zap/zapcore"
)

const ApiUrl = "https://api.telegram.org/bot"

var DefaultTgBotRequester = BaseRequester{
	Client: http.Client{
		Timeout: time.Millisecond * 1500,
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

type BaseRequester struct {
	Client http.Client
	ApiUrl string
}

type TelegramError struct {
	Method      string
	Values      url.Values
	Code        int
	Description string
}

func (t *TelegramError) Error() string {
	return fmt.Sprintf("unable to %s: %s", t.Method, t.Description)
}

type Requester interface {
	Get(l *zap.SugaredLogger, token string, method string, params url.Values) (json.RawMessage, error)
	Post(l *zap.SugaredLogger, token string, method string, params url.Values, fileType string, file io.Reader, filename string) (json.RawMessage, error)
}

func (tbg BaseRequester) Get(l *zap.SugaredLogger, token string, method string, params url.Values) (json.RawMessage, error) {
	endpoint := tbg.ApiUrl
	if endpoint == "" {
		endpoint = ApiUrl
	}

	req, err := http.NewRequest("GET", endpoint+token+"/"+method, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to build GET request to %v", method)
	}
	req.URL.RawQuery = params.Encode()

	l.Debugf("executing GET: %+v", req)
	resp, err := tbg.Client.Do(req)
	if err != nil {
		l.Debugw("failed to execute GET request",
			zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: method,
			},
			zap.Error(err))
		return nil, errors.Wrapf(err, "client error executing GET request to %v", method)
	}
	defer resp.Body.Close()
	l.Debugf("successful GET request: %+v", resp)

	var r response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		l.Debugw("failed to deserialize GET response body",
			zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: method,
			},
			zap.Error(err))
		return nil, errors.Wrapf(err, "decoding error while reading GET request to %v", method)
	}
	if !r.Ok {
		l.Debugw("error from GET",
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
			Method:      method,
			Values:      params,
			Code:        r.ErrorCode,
			Description: r.Description,
		}
	}

	l.Debugw("obtained GET result",
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

func (tbg BaseRequester) Post(l *zap.SugaredLogger, token string, method string, params url.Values, fileType string, file io.Reader, filename string) (json.RawMessage, error) {
	endpoint := tbg.ApiUrl
	if endpoint == "" {
		endpoint = ApiUrl
	}

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

	req, err := http.NewRequest("POST", endpoint+token+"/"+method, &b)
	if err != nil {
		l.Debugw("failed to execute POST request",
			zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: method,
			},
			zap.Error(err))
		return nil, errors.Wrapf(err, "client error executing POST request to %v", method)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Content-Type", w.FormDataContentType())

	l.Debugf("POST request with body: %+v", b)
	l.Debugf("executing POST: %+v", req)
	resp, err := tbg.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	l.Debugf("successful POST request: %+v", resp)

	var r response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		l.Debugw("failed to deserialize POST response body",
			zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: method,
			},
			zap.Error(err))
		return nil, errors.Wrapf(err, "decoding error while reading POST request to %v", method)
	}
	if !r.Ok {
		l.Debugw("error from POST",
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
			Method:      method,
			Values:      params,
			Code:        r.ErrorCode,
			Description: r.Description,
		}
	}

	l.Debugw("obtained POST result",
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
