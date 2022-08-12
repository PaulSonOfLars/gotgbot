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

type BotClient interface {
	// RequestWithContext submits a POST HTTP request a bot API instance.
	RequestWithContext(ctx context.Context, method string, params map[string]string, data map[string]NamedReader, opts *RequestOpts) (json.RawMessage, error)
	// TimeoutContext calculates the required timeout contect required given the passed RequestOpts, and any default opts defined by the BotClient.
	TimeoutContext(opts *RequestOpts) (context.Context, context.CancelFunc)
	// GetAPIURL gets the URL of the API the bot is interacting with.
	GetAPIURL() string
	// GetToken gets the current bots' token.
	GetToken() string
}

type BaseBotClient struct {
	// Token stores the bot's secret token obtained from t.me/BotFather, and used to interact with telegram's API.
	Token string
	// Client is the HTTP Client used for all HTTP requests made for this bot.
	Client http.Client
	// UseTestEnvironment defines whether this bot was created to run on telegram's test environment.
	// Enabling this uses a slightly different API path.
	// See https://core.telegram.org/bots/webapps#using-bots-in-the-test-environment for more details.
	UseTestEnvironment bool
	// The default request opts for this bot instance.
	DefaultRequestOpts *RequestOpts
}

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

// TimeoutContext returns the appropriate context for the current settings.
func (bot *BaseBotClient) TimeoutContext(opts *RequestOpts) (context.Context, context.CancelFunc) {
	if opts != nil {
		ctx, cancelFunc := timeoutFromOpts(opts)
		if ctx != nil {
			return ctx, cancelFunc
		}
	}

	if bot.DefaultRequestOpts != nil {
		ctx, cancelFunc := timeoutFromOpts(bot.DefaultRequestOpts)
		if ctx != nil {
			return ctx, cancelFunc
		}
	}

	return context.WithTimeout(context.Background(), DefaultTimeout)
}

func timeoutFromOpts(opts *RequestOpts) (context.Context, context.CancelFunc) {
	// nothing? no timeout.
	if opts == nil {
		return nil, nil
	}

	if opts.Timeout > 0 {
		// > 0 timeout defined.
		return context.WithTimeout(context.Background(), opts.Timeout)

	} else if opts.Timeout < 0 {
		// < 0  no timeout; infinite.
		return context.Background(), func() {}
	}
	// 0 == nothing defined, use defaults.
	return nil, nil
}

// PostWithContext allows sending a POST request to the telegram bot API with an existing context.
//   - ctx: the timeout contexts to be used.
//   - method: the telegram API method to call.
//   - params: map of parameters to be sending to the telegram API. eg: chat_id, user_id, etc.
//   - data: map of any files to be sending to the telegram API.
//   - opts: request opts to use. Note: Timeout opts are ignored when used in PostWithContext. Timeout handling is the
//     responsibility of the caller/context owner.
func (bot *BaseBotClient) RequestWithContext(ctx context.Context, method string, params map[string]string, data map[string]NamedReader, opts *RequestOpts) (json.RawMessage, error) {
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
func (bot *BaseBotClient) GetAPIURL() string {
	return bot.getAPIURL(nil)
}

// GetToken returns the currently used token.
func (bot *BaseBotClient) GetToken() string {
	return bot.Token
}

// getAPIURL returns the currently used API endpoint.
func (bot *BaseBotClient) getAPIURL(opts *RequestOpts) string {
	if opts != nil && opts.APIURL != "" {
		return getCleanAPIURL(opts.APIURL)
	}
	if bot.DefaultRequestOpts != nil && bot.DefaultRequestOpts.APIURL != "" {
		return getCleanAPIURL(bot.DefaultRequestOpts.APIURL)
	}
	return DefaultAPIURL
}

func (bot *BaseBotClient) methodEnpoint(method string, opts *RequestOpts) string {
	if bot.UseTestEnvironment {
		return fmt.Sprintf("%s/bot%s/test/%s", bot.getAPIURL(opts), bot.Token, method)
	}
	return fmt.Sprintf("%s/bot%s/%s", bot.getAPIURL(opts), bot.Token, method)
}
