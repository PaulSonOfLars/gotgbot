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
	// DefaultAPIURL is the default telegram API URL.
	DefaultAPIURL = "https://api.telegram.org/bot"
	// DefaultGetTimeout is the default timeout to be set for a GET request.
	DefaultGetTimeout = time.Second * 3
	// DefaultPostTimeout is the default timeout to be set for a POST request.
	DefaultPostTimeout = time.Second * 10
)

type Response struct {
	// Ok: if true, request was successful, and result can be found in the Result field.
	// If false, error can be explained in the Description.
	Ok bool `json:"ok"`
	// Result: result of requests (if Ok)
	Result json.RawMessage `json:"result"`
	// ErrorCode: Integer error code of request. Subject to change in the future.
	ErrorCode int `json:"error_code"`
	// Description: contains a human readable description of the error result.
	Description string `json:"description"`
	// Parameters: Optional extra data which can help automatically handle the error.
	Parameters *ResponseParameters `json:"parameters"`
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

// GetWithContext allows sending a Get request with an existing context.
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

// PostWithContext allows sending a Post request with an existing context.
func (bot *Bot) PostWithContext(ctx context.Context, method string, params url.Values, data map[string]NamedReader) (json.RawMessage, error) {
	b := &bytes.Buffer{}
	contentType := "application/json"

	if len(data) > 0 {
		var err error
		contentType, err = fillBuffer(b, data)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", bot.endpoint(method), b)
	if err != nil {
		return nil, fmt.Errorf("failed to build POST request to %s: %w", method, err)
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Set("Content-Type", contentType)

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

func fillBuffer(b *bytes.Buffer, data map[string]NamedReader) (string, error) {
	w := multipart.NewWriter(b)

	for field, file := range data {
		fileName := file.Name()
		if fileName == "" {
			fileName = field
		}

		part, err := w.CreateFormFile(field, fileName)
		if err != nil {
			return "", fmt.Errorf("failed to create form file for field %s and fileName %s: %w", field, fileName, err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return "", fmt.Errorf("failed to copy file contents of field %s to form: %w", field, err)
		}
	}

	err := w.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close multipart form writer: %w", err)
	}

	return w.FormDataContentType(), nil
}

func (bot *Bot) endpoint(method string) string {
	if bot.APIURL == "" {
		return DefaultAPIURL + bot.Token + "/" + method
	}
	return bot.APIURL + bot.Token + "/" + method
}
