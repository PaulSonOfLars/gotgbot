package gotgbot

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"strconv"
	"testing"
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
	r         io.Reader
	bytesRead int64
}

func (c *countingReader) Read(p []byte) (n int, err error) {
	n, err = c.r.Read(p)
	c.bytesRead += int64(n)
	return
}

func TestFileNotBufferedIntoMemory(t *testing.T) {
	fileContents := []byte("hello, this is some file content")

	cr := &countingReader{r: bytes.NewReader(fileContents)}

	// Wire up a FileReader with the counting reader as Data.
	// Intentionally NOT an io.Seeker — we want to catch greedy reads.
	params := map[string]any{
		"document": &FileReader{
			Name: "test.txt",
			Data: cr,
		},
	}

	bs, _, err := buildMultipart(params)
	if err != nil {
		t.Fatalf("unexpected error building multipart: %v", err)
	}

	// Before draining: file should not have been read yet (streaming, not buffered)
	if cr.bytesRead != 0 {
		t.Errorf("file was read during multipart construction: %d bytes read, expected 0", cr.bytesRead)
	}

	// Drain the multipart body as the HTTP transport would.
	if _, err := io.Copy(io.Discard, bytes.NewBuffer(bs)); err != nil {
		t.Fatalf("unexpected error reading multipart body: %v", err)
	}

	// Sanity check: file was actually sent (and was only read once)
	if cr.bytesRead != int64(len(fileContents)) {
		t.Errorf("file was read %d bytes after drain, expected exactly %d", cr.bytesRead, len(fileContents))
	}
}

// We want to make sure that retriable requests can in fact be retried - this simulates the HTTP2 GOAWAY state.
func TestGetBodyReturnsCorrectRetryReader(t *testing.T) {
	fileContents := []byte("hello, this is some file content")

	params := map[string]any{
		"document": &FileReader{
			Name: "test.txt",
			Data: bytes.NewReader(fileContents),
		},
	}

	req, err := (&BaseBotClient{}).buildRequest(params, context.Background(), "test-token", "sendDocument", nil)
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

	retryBody, err := req.GetBody()
	if err != nil {
		t.Fatalf("GetBody returned error: %v", err)
	}

	retryBodyBytes, err := io.ReadAll(retryBody)
	if err != nil {
		t.Fatalf("failed to read retry body: %v", err)
	}

	if !bytes.Equal(originalBody, retryBodyBytes) {
		t.Errorf("retry body does not match original body")
	}
}
