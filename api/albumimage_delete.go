package api

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"github.com/gorilla/mux"
	"net/http"
)

func DeleteAlbumImage(env environment.FrontendEnvironment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		albumImage, err := env.Repositories.AlbumImages.Get(vars["albumId"], vars["imageId"])
		if err == sql.ErrNoRows {
			http.NotFound(writer, request)
			return
		} else if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		err = env.Repositories.AlbumImages.Delete(albumImage)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	})
}
