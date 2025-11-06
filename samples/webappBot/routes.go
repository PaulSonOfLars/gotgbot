package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var indexTmpl = template.Must(template.ParseFiles("index.html"))

func index(webappURL string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := indexTmpl.ExecuteTemplate(writer, "index.html", struct {
			WebAppURL string
		}{
			WebAppURL: webappURL,
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}

func validate(token string) func(writer http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Our index.html sends the WebApp.initData field over the X-Auth header.
		// We parse this string as a URL query.
		authQuery, err := url.ParseQuery(r.Header.Get("X-Auth"))
		if err != nil {
			http.Error(w, "validation failed; failed to parse auth query: "+err.Error(), http.StatusBadRequest)
			return
		}

		// We validate that the query has been hashed correctly, ensuring data can be trusted.
		ok, err := ext.ValidateWebAppQuery(authQuery, token)
		if err != nil {
			http.Error(w, "validation failed; error: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if !ok {
			http.Error(w, "validation failed; data cannot be trusted.", http.StatusUnauthorized)
			return
		}

		// Once we've confirmed the data can be trusted, we unmarshal any data we may need to use.
		var u gotgbot.User
		err = json.Unmarshal([]byte(authQuery.Get("user")), &u)
		if err != nil {
			http.Error(w, "validation failed; failed to unmarshal user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// And then we can choose to either return it, or work with it.
		fmt.Fprintf(w, "validation success; user '%s' is authenticated (id: %d).", u.FirstName, u.Id)
	}
}
