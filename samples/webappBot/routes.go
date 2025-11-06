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
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(writer, err.Error())
		}
	}
}

func validate(token string) func(writer http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Our index.html sends the WebApp.initData field over the X-Auth header.
		// We parse this string as a URL query.
		authQuery, err := url.ParseQuery(r.Header.Get("X-Auth"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "validation failed; failed to parse auth query: "+err.Error())
		}

		// We validate that the query has been hashed correctly, ensuring data can be trusted.
		ok, err := ext.ValidateWebAppQuery(authQuery, token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "validation failed; error: "+err.Error())
			return
		}
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "validation failed; data cannot be trusted.")
			return
		}

		// Once we've confirmed the data can be trusted, we unmarshal any data we may need to use.
		var u gotgbot.User
		err = json.Unmarshal([]byte(authQuery.Get("user")), &u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "validation failed; failed to unmarshal user: "+err.Error())
			return
		}

		// And then we can choose to either return it, or work with it.
		fmt.Fprintf(w, "validation success; user '%s' is authenticated (id: %d).", u.FirstName, u.Id)
	}
}
