package main

import (
	"net/http"
	"text/template"

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
			writer.Write([]byte(err.Error()))
		}
	}
}

func validate(token string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ok, err := ext.ValidateWebAppQuery(request.URL.Query(), token)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("validation failed; error: " + err.Error()))
			return
		}
		if ok {
			writer.Write([]byte("validation success; user is authenticated."))
		} else {
			writer.Write([]byte("validation failed; data cannot be trusted."))
		}
	}
}
