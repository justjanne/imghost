package api

import (
	"database/sql"
	"encoding/json"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"github.com/gorilla/mux"
	"net/http"
)

func UpdateAlbumImage(env environment.FrontendEnvironment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var err error

		vars := mux.Vars(request)
		image, err := env.Repositories.AlbumImages.Get(
			vars["albumId"],
			vars["imageId"])
		if err == sql.ErrNoRows {
			http.NotFound(writer, request)
			return
		} else if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		var changes model.AlbumImage
		err = json.NewDecoder(request.Body).Decode(&changes)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		image.Title = changes.Title
		image.Description = changes.Description

		err = env.Repositories.AlbumImages.Update(image)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	})
}
