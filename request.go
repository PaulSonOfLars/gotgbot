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
	DefaultTimeout = time.Second * 10
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

func (bot *Bot) Post(method string, params map[string]string, data map[string]NamedReader) (json.RawMessage, error) {
	if bot.RequestTimeout == 0 {
		bot.RequestTimeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), bot.RequestTimeout)
	defer cancel()

	return bot.PostWithContext(ctx, method, params, data)
}

// PostWithContext allows sending a Post request with an existing context.
func (bot *Bot) PostWithContext(ctx context.Context, method string, params map[string]string, data map[string]NamedReader) (json.RawMessage, error) {
	b := &bytes.Buffer{}
	contentType := "application/json"

	// Check if there are any files to upload
	if len(data) > 0 {
		var err error
		contentType, err = fillBuffer(b, params, data)
		if err != nil {
			return nil, err
		}
	} else {
		err := json.NewEncoder(b).Encode(params)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bot.methodEnpoint(method), b)
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
			return "", fmt.Errorf("failed to write multipart field %s with vale %s: %w", k, v, err)
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

// GetAPIURL returns the currently used API endpoint.
func (bot *Bot) GetAPIURL() string {
	if bot.APIURL == "" {
		return DefaultAPIURL
	}
	// Trim suffix to ensure consistent output
	return strings.TrimSuffix(bot.APIURL, "/")
}

func (bot *Bot) methodEnpoint(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", bot.GetAPIURL(), bot.Token, method)
}
