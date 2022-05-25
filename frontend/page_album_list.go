package main

import (
	"net/http"
)

func pageAlbumList(env PageEnvironment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
