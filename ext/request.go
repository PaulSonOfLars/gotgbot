package ext

import (
	"net/http"
	"encoding/json"
	"net/url"
	"time"
	"github.com/pkg/errors"
	"io"
	"github.com/sirupsen/logrus"
	"bytes"
	"mime/multipart"
)

var apiUrl = "https://api.telegram.org/bot"

var client = &http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       time.Second * 5,
}

type Response struct {
	Ok          bool
	Result      json.RawMessage
	ErrorCode   int `json:"error_code"`
	Description string
	Parameters  json.RawMessage
}

func Get(bot Bot, method string, params url.Values) (*Response, error) {
	req, err := http.NewRequest("GET", apiUrl+bot.Token+"/"+method, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to build get request to %v", method)
	}
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to execute get request to %v", method)
	}
	defer resp.Body.Close()

	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrapf(err, "could not decode in Get %v call", method)
	}
	return &r, nil
}

func Post(bot Bot, fileType string, method string, params url.Values, file io.Reader, filename string) (*Response, error) {
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

	req, err := http.NewRequest("POST", apiUrl+bot.Token+"/"+method, &b)
	if err != nil {
		logrus.WithError(err).Errorf("failed to send to %v func", method)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
