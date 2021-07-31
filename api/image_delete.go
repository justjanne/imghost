package api

import (
	"context"
	"database/sql"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"github.com/gorilla/mux"
	"net/http"
)

func DeleteImage(env environment.FrontendEnvironment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var err error

		vars := mux.Vars(request)
		image, err := env.Repositories.Images.Get(vars["imageId"])
		if err == sql.ErrNoRows {
			http.NotFound(writer, request)
			return
		} else if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		err = env.Repositories.Images.Delete(image)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		err = env.Storage.DeleteFiles(
			context.Background(),
			env.Configuration.Storage.ImageBucket,
			image.Id,
		)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	})
}
