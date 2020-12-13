package gotgbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

const (
	// Default telegram API URL.
	DefaultAPIURL = "https://api.telegram.org/bot"
	// Default timeout to be set for a GET request.
	DefaultGetTimeout = time.Second
	// Default timeout to be set for a POST request.
	DefaultPostTimeout = time.Second * 3
)

type Response struct {
	Ok          bool               `json:"ok"`
	Result      json.RawMessage    `json:"result"`
	ErrorCode   int                `json:"error_code"`
	Description string             `json:"description"`
	Parameters  ResponseParameters `json:"parameters"`
}

type TelegramError struct {
	Method      string
	Params      url.Values
	Code        int
	Description string
}

func (t *TelegramError) Error() string {
	return fmt.Sprintf("unable to %s: %s", t.Method, t.Description)
}

type NamedReader interface {
	Name() string
	io.Reader
}

type NamedFile struct {
	File     io.Reader
	FileName string
}

func (nf NamedFile) Read(p []byte) (n int, err error) {
	return nf.File.Read(p)
}

func (nf NamedFile) Name() string {
	return nf.FileName
}

func (bot *Bot) Get(method string, params url.Values) (json.RawMessage, error) {
	if bot.GetTimeout == 0 {
		bot.GetTimeout = DefaultGetTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), bot.GetTimeout)
	defer cancel()

	return bot.GetWithContext(ctx, method, params)
}

func (bot *Bot) GetWithContext(ctx context.Context, method string, params url.Values) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", bot.endpoint(method), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build GET request to %s: %w", method, err)
	}

	req.URL.RawQuery = params.Encode()

	resp, err := bot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GET request to %s: %w", method, err)
	}
	defer resp.Body.Close()

	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode GET request to %s: %w", method, err)
	}

	if !r.Ok {
		return nil, &TelegramError{
			Method:      method,
			Params:      params,
			Code:        r.ErrorCode,
			Description: r.Description,
		}
	}

	return r.Result, nil
}

func (bot *Bot) Post(method string, params url.Values, data map[string]NamedReader) (json.RawMessage, error) {
	if bot.PostTimeout == 0 {
		bot.PostTimeout = DefaultPostTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), bot.PostTimeout)
	defer cancel()

	return bot.PostWithContext(ctx, method, params, data)
}

func (bot *Bot) PostWithContext(ctx context.Context, method string, params url.Values, data map[string]NamedReader) (json.RawMessage, error) {
	b := bytes.Buffer{}
	w := multipart.NewWriter(&b)
	defer w.Close()

	for field, file := range data {
		fileName := file.Name()
		if fileName == "" {
			fileName = field
		}

		part, err := w.CreateFormFile(field, fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file for field %s and fileName %s: %w", field, fileName, err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, fmt.Errorf("failed to copy file contents of field %s to form: %w", field, err)
		}
	}

	err := w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart form writer: %w", err)
	}

	API := bot.APIURL
	if API == "" {
		API = DefaultAPIURL
	}

	req, err := http.NewRequestWithContext(ctx, "POST", bot.endpoint(method), &b)
	if err != nil {
		return nil, fmt.Errorf("failed to build POST request to %s: %w", method, err)
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := bot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute POST request to %s: %w", method, err)
	}
	defer resp.Body.Close()

	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode POST request to %s: %w", method, err)
	}

	if !r.Ok {
		return nil, &TelegramError{
			Method:      method,
			Params:      params,
			Code:        r.ErrorCode,
			Description: r.Description,
		}
	}

	return r.Result, nil
}

func (bot *Bot) endpoint(method string) string {
	if bot.APIURL == "" {
		return DefaultAPIURL + bot.Token + "/" + method
	}
	return bot.APIURL + bot.Token + "/" + method
}
