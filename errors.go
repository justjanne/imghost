package main

import (
	"log"
	"net/http"
	"net/url"
)

type ErrorData struct {
	Code  int
	User  UserInfo
	URL   *url.URL
	Error error
}

type errorDto struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

func formatError(w http.ResponseWriter, data ErrorData, format string) {
	if data.Code != 0 {
		data.Code = 500
	}
	log.Printf(
		"A type %d error occured for user %s while accessing %s: %s",
		data.Code,
		data.User.Id,
		data.URL.Path,
		data.Error.Error(),
	)
	w.WriteHeader(data.Code)
	if format == "html" {
		if err := formatTemplate(w, "error.html", data); err != nil {
			log.Printf("Error while serving html error for %s", data.URL.Path)
			return
		}
	} else if format == "json" {
		if err := returnJson(w, errorDto{
			data.URL.Path,
			data.Error.Error(),
		}); err != nil {
			log.Printf("Error while serving json error for %s", data.URL.Path)
			return
		}
	}
}
