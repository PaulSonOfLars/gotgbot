package ext

import (
	"net/http"
	"encoding/json"
	"net/url"
	"time"
	"github.com/pkg/errors"
)

var apiUrl = "https://api.telegram.org/bot"

var client = &http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout: time.Second * 5,
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
		return nil, errors.Wrapf(err, "unable to build get request")
	}
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to execute get request")
	}
	defer resp.Body.Close()

	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrapf(err, "could not decode in Get call")
	}
	return &r, nil
}
