package gotgbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const (
	// DefaultAPIURL is the default telegram API URL.
	DefaultAPIURL = "https://api.telegram.org"
	// DefaultTimeout is the default timeout to be set for all requests.
	DefaultTimeout = time.Second * 5
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
	Params      map[string]string
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

// RequestOpts defines any request-specific options used to interact with the telegram API.
type RequestOpts struct {
	// Timeout for the HTTP request to the telegram API.
	Timeout time.Duration
	// Custom API URL to use for requests.
	APIURL string
}

// Post sends a POST request to the telegram bot API.
func (bot *Bot) Post(method string, params map[string]string, data map[string]NamedReader, opts *RequestOpts) (json.RawMessage, error) {
	ctx, cancel := bot.getTimeoutContext(opts)
	defer cancel()

	return bot.PostWithContext(ctx, method, params, data, opts)
}

// getTimeoutContext returns the appropriate context for the current settings.
func (bot *Bot) getTimeoutContext(opts *RequestOpts) (context.Context, context.CancelFunc) {
	if opts != nil {
		if opts.Timeout > 0 {
			// custom timeout defined
			return context.WithTimeout(context.Background(), opts.Timeout)

		} else if opts.Timeout < 0 {
			return context.Background(), func() {}
		}
	}

	if bot.DefaultRequestOpts != nil {
		if bot.DefaultRequestOpts.Timeout > 0 {
			return context.WithTimeout(context.Background(), bot.DefaultRequestOpts.Timeout)
		} else if bot.DefaultRequestOpts.Timeout < 0 {
			return context.Background(), func() {}
		}
	}

	return context.WithTimeout(context.Background(), DefaultTimeout)
}

// PostWithContext allows sending a POST request to the telegram bot API with an existing context.
// - ctx: the timeout contexts to be used.
// - method: the telegram API method to call.
// - params: map of parameters to be sending to the telegram API. eg: chat_id, user_id, etc.
// - data: map of any files to be sending to the telegram API.
// - opts: request opts to use.
func (bot *Bot) PostWithContext(ctx context.Context, method string, params map[string]string, data map[string]NamedReader, opts *RequestOpts) (json.RawMessage, error) {
	b := &bytes.Buffer{}

	var contentType string
	// Check if there are any files to upload. If yes, use multipart; else, use JSON.
	if len(data) > 0 {
		var err error
		contentType, err = fillBuffer(b, params, data)
		if err != nil {
			return nil, fmt.Errorf("failed to fill buffer with parameters and file data: %w", err)
		}
	} else {
		contentType = "application/json"
		err := json.NewEncoder(b).Encode(params)
		if err != nil {
			return nil, fmt.Errorf("failed to encode parameters as JSON: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bot.methodEnpoint(method, opts), b)
	if err != nil {
		return nil, fmt.Errorf("failed to build POST request to %s: %w", method, err)
	}

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

func fillBuffer(b *bytes.Buffer, params map[string]string, data map[string]NamedReader) (string, error) {
	w := multipart.NewWriter(b)

	for k, v := range params {
		err := w.WriteField(k, v)
		if err != nil {
			return "", fmt.Errorf("failed to write multipart field %s with value %s: %w", k, v, err)
		}
	}

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

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart form writer: %w", err)
	}

	return w.FormDataContentType(), nil
}

func getCleanAPIURL(url string) string {
	if url == "" {
		return DefaultAPIURL
	}
	// Trim suffix to ensure consistent output
	return strings.TrimSuffix(url, "/")
}

// GetAPIURL returns the currently used API endpoint.
func (bot *Bot) GetAPIURL() string {
	return bot.getAPIURL(nil)
}

// getAPIURL returns the currently used API endpoint.
func (bot *Bot) getAPIURL(opts *RequestOpts) string {
	if opts != nil && opts.APIURL != "" {
		return getCleanAPIURL(opts.APIURL)
	}
	if bot.DefaultRequestOpts != nil && bot.DefaultRequestOpts.APIURL != "" {
		return getCleanAPIURL(bot.DefaultRequestOpts.APIURL)
	}
	return DefaultAPIURL
}

func (bot *Bot) methodEnpoint(method string, opts *RequestOpts) string {
	return fmt.Sprintf("%s/bot%s/%s", bot.getAPIURL(opts), bot.Token, method)
}
