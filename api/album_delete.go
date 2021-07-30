package api

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"github.com/gorilla/mux"
	"net/http"
)

func DeleteAlbum(env environment.FrontendEnvironment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		album, err := env.Repositories.Albums.Get(vars["albumId"])
		if err == sql.ErrNoRows {
			http.NotFound(writer, request)
			return
		} else if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		err = env.Repositories.Albums.Delete(album)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	})
}
