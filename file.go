package gotgbot

import (
	"encoding/json"
	"errors"
	"io"
)

// InputFile (https://core.telegram.org/bots/api#inputfile)
//
// This object represents the contents of a file to be uploaded.
// Must be posted using multipart/form-data in the usual way that files are uploaded via the browser.
type InputFile interface {
	InputFileOrString
	justFiles()
}

// InputString (https://core.telegram.org/bots/api#inputfile)
//
// This object represents a publicly accessible URL to be reused, or a file_id already available on telegram servers.
type InputString interface {
	InputFileOrString
	justStrings()
}

// InputFileOrString (https://core.telegram.org/bots/api#inputfile)
//
// This object represents the contents of a file to be uploaded, or a publicly accessible URL to be reused.
// Files must be posted using multipart/form-data in the usual way that files are uploaded via the browser.
type InputFileOrString interface {
	Attach(name string, data map[string]FileReader) error
	getValue() string
}

var (
	_ InputFileOrString = &FileString{}
	_ InputString       = &FileString{}

	_ InputFileOrString = &FileReader{}
	_ InputFile         = &FileReader{}
)

type FileString struct {
	Value string
}

func (f *FileString) justStrings() {}

func (f *FileString) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.getValue())
}

func (f *FileString) Attach(_ string, _ map[string]FileReader) error {
	return nil
}

func (f *FileString) getValue() string {
	return f.Value
}

type FileReader struct {
	Name string
	Data io.Reader

	value string
}

func (f *FileReader) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.getValue())
}

var ErrAttachmentKeyAlreadyExists = errors.New("key already exists")

func (f *FileReader) justFiles() {}

func (f *FileReader) Attach(key string, data map[string]FileReader) error {
	if _, ok := data[key]; ok {
		return ErrAttachmentKeyAlreadyExists
	}
	f.value = key
	data[key] = *f
	return nil
}

// getValue returns the file attach reference for the relevant multipart form.
// Calling getValue is only valid after having called Attach on this struct.
func (f *FileReader) getValue() string {
	return "attach://" + f.value
}

// InputFileByURL is used to send a file on the internet via a publicly accessible HTTP URL.
func InputFileByURL(url string) InputFileOrString {
	return &FileString{Value: url}
}

// InputFileByID is used to send a file that is already present on telegram's servers, using its telegram file_id.
func InputFileByID(fileID string) InputFileOrString {
	return &FileString{Value: fileID}
}

// InputFileByReader is used to send a file by a reader interface; such as a filehandle from os.Open(), or from a byte
// buffer.
//
// For example:
//
//	f, err := os.Open("some_file.go")
//	if err != nil {
//		return fmt.Errorf("failed to open file: %w", err)
//	}
//
//	m, err := b.SendDocument(<chat_id>, gotgbot.InputFileByReader("source.go", f), nil)
//
// Or
//
//	m, err := b.SendDocument(<chat_id>, gotgbot.InputFileByReader("file.txt", strings.NewReader("Some file contents")), nil)
func InputFileByReader(name string, r io.Reader) InputFile {
	return &FileReader{Name: name, Data: r}
}
