package gotgbot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_getFieldContents(t *testing.T) {
	var testString = "test"
	var testInt = 42

	type args struct {
		v any
		k string
		w *multipart.Writer
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "regular string",
			args: args{
				v: testString,
			},
			want:    testString,
			wantErr: false,
		}, {
			name: "string pointer",
			args: args{
				v: &testString,
			},
			want:    testString,
			wantErr: false,
		}, {
			name: "integer",
			args: args{
				v: testInt,
			},
			want:    strconv.Itoa(testInt),
			wantErr: false,
		}, {
			name: "integer pointer",
			args: args{
				v: &testInt,
			},
			want:    strconv.Itoa(testInt),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFieldContents(tt.args.v, tt.args.k, tt.args.w)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFieldContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFieldContents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type countingReader struct {
	r         io.ReadSeeker
	bytesRead int64
}

func (c *countingReader) Read(p []byte) (n int, err error) {
	n, err = c.r.Read(p)
	c.bytesRead += int64(n)
	return
}

func (c *countingReader) Seek(offset int64, whence int) (int64, error) {
	return c.r.Seek(offset, whence)
}

func TestFileNotBufferedIntoMemory(t *testing.T) {
	fileContents := []byte("hello, this is some file content")

	cr := &countingReader{r: bytes.NewReader(fileContents)}

	// Wire up a FileReader with the counting reader as Data.
	// Intentionally NOT an io.Seeker - we want to catch greedy reads.
	params := map[string]any{
		"document": &FileReader{
			Name: "test.txt",
			Data: cr,
		},
	}

	r, _ := buildMultipart(context.Background(), params)

	// Before draining: file should not have been read yet (streaming, not buffered).
	if cr.bytesRead != 0 {
		t.Errorf("file was read during multipart construction: %d bytes read, expected 0", cr.bytesRead)
	}

	// Drain the multipart body as the HTTP transport would.
	if _, err := io.Copy(io.Discard, r); err != nil {
		t.Fatalf("unexpected error reading multipart body: %v", err)
	}

	// Which means that we should now have read the file.
	if cr.bytesRead != int64(len(fileContents)) {
		t.Errorf("file was read %d bytes after drain, expected exactly %d", cr.bytesRead, len(fileContents))
	}
}

// We want to make sure that retriable requests can in fact be retried - this simulates the HTTP2 GOAWAY state.
func TestGetBodyReturnsCorrectRetryReader(t *testing.T) {
	fileContents := []byte("hello, this is some file content")
	cr := &countingReader{r: bytes.NewReader(fileContents)}

	params := map[string]any{
		"document": &FileReader{
			Name: "test.txt",
			Data: cr,
		},
	}

	req, err := (&BaseBotClient{}).buildRequest(context.Background(), params, "test-token", "sendDocument", nil)
	if err != nil {
		t.Fatalf("unexpected error building request: %v", err)
	}

	if req.GetBody == nil {
		t.Fatal("GetBody should be set for seekable files")
	}

	// Read the request body, simulating our "first request" hitting an HTTP2 GOAWAY
	// GOAWAY would discard the body, but we store it to check equality.
	originalBody, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("failed to read original body: %v", err)
	}
	origContentType := req.Header.Get("Content-Type")

	// We should now have read the file once
	if cr.bytesRead != int64(len(fileContents)) {
		t.Errorf("expected file to be read exactly once after first attempt, got %d bytes read", cr.bytesRead)
	}

	retryBody, err := req.GetBody()
	if err != nil {
		t.Fatalf("GetBody returned error: %v", err)
	}

	retryBodyBytes, err := io.ReadAll(retryBody)
	if err != nil {
		t.Fatalf("failed to read retry body: %v", err)
	}
	retryContentType := req.Header.Get("Content-Type")

	// Retry - we have now read the file twice
	if cr.bytesRead != int64(len(fileContents))*2 {
		t.Errorf("expected file to be read exactly twice after retry, got %d bytes read", cr.bytesRead)
	}

	originalFile := extractFileFromMultipart(t, origContentType, originalBody, "document")
	if !bytes.Equal(fileContents, originalFile) {
		t.Errorf("local file contents do not match original file")
	}

	retryFile := extractFileFromMultipart(t, retryContentType, retryBodyBytes, "document")
	if !bytes.Equal(originalFile, retryFile) {
		t.Errorf("retry file do not match original file")
	}
}

func extractFileFromMultipart(t *testing.T, contentType string, body []byte, formName string) []byte {
	t.Helper()

	_, mediaParams, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Fatalf("failed to parse content-type: %v", err)
	}

	mr := multipart.NewReader(bytes.NewReader(body), mediaParams["boundary"])
	for {
		part, err := mr.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("failed to read multipart part: %v", err)
		}

		if part.FormName() == formName {
			contents, err := io.ReadAll(part)
			if err != nil {
				t.Fatalf("failed to read file part: %v", err)
			}
			return contents
		}
	}
	t.Fatalf("%s part not found in multipart body", formName)
	return nil
}

func TestBuildMultipartContextHandling(t *testing.T) {
	params := map[string]any{
		"chat_id": "123",
		"text":    "hello",
	}

	t.Run("body is readable with active context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r, _ := buildMultipart(ctx, params)
		_, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("expected clean read, got: %v", err)
		}
	})

	t.Run("body read fails after context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		r, _ := buildMultipart(ctx, params)
		cancel()
		time.Sleep(10 * time.Millisecond)

		_, err := io.ReadAll(r)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got: %v", err)
		}
	})

	t.Run("body read fails after context deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		r, _ := buildMultipart(ctx, params)
		time.Sleep(20 * time.Millisecond)

		_, err := io.ReadAll(r)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context.DeadlineExceeded, got: %v", err)
		}
	})

	t.Run("mid-read fails after context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		r, _ := buildMultipart(ctx, params)

		buf := make([]byte, 4)
		if _, err := r.Read(buf); err != nil {
			t.Fatalf("unexpected error on first read: %v", err)
		}

		cancel()
		time.Sleep(10 * time.Millisecond)

		_, err := io.ReadAll(r)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled after mid-read cancel, got: %v", err)
		}
	})
}

func TestHasAttachments(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want bool
	}{
		{
			name: "string parameter",
			val:  "hello",
			want: false,
		},
		{
			name: "integer parameter",
			val:  42,
			want: false,
		},
		{
			name: "file reader (implements Attach)",
			val:  &FileReader{Data: nil},
			want: true,
		},
		{
			name: "media photo (implements Attach)",
			val: InputMediaPhoto{
				Media: &FileReader{Data: nil},
			},
			want: true,
		},
		{
			name: "slice of media (implements Attach)",
			val: InputMedias{
				InputMediaPhoto{
					Media: &FileReader{Data: nil},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := map[string]any{
				"test": tt.val,
			}
			got := hasAttachments(params)
			if got != tt.want {
				t.Errorf("hasAttachments(%v) = %v; want %v", tt.val, got, tt.want)
			}
		})
	}
}

func TestJSONRequestSerialization(t *testing.T) {
	params := map[string]any{
		"chat_id": 123456,
		"text":    "hello world",
	}

	bot := &BaseBotClient{}
	req, err := bot.buildRequest(context.Background(), params, "token", "sendMessage", nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", req.Header.Get("Content-Type"))
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("failed to read request body: %v", err)
	}

	var gotParams map[string]any
	if err := json.Unmarshal(bodyBytes, &gotParams); err != nil {
		t.Fatalf("failed to unmarshal request body: %v", err)
	}

	if gotParams["chat_id"].(float64) != 123456 || gotParams["text"] != "hello world" {
		t.Errorf("unexpected parameters in JSON body: %v", gotParams)
	}
}
