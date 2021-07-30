package api

import (
	"database/sql"
	"encoding/json"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"github.com/gorilla/mux"
	"net/http"
)

func ReorderAlbum(env environment.FrontendEnvironment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var err error

		vars := mux.Vars(request)
		album, err := env.Repositories.Albums.Get(vars["albumId"])
		if err == sql.ErrNoRows {
			http.NotFound(writer, request)
			return
		} else if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		var changes []model.AlbumImage
		err = json.NewDecoder(request.Body).Decode(&changes)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		for index, image := range changes {
			image.Album = album.Id

			err = env.Repositories.AlbumImages.Reorder(image, index)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		writer.WriteHeader(http.StatusNoContent)
	})
}
