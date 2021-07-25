package api

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"github.com/gorilla/mux"
	"net/http"
)

func GetImage(env environment.Environment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		image, err := env.Repositories.Images.Get(vars["imageId"])
		if err == sql.ErrNoRows {
			http.NotFound(writer, request)
			return
		} else if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		util.ReturnJson(writer, image)
	})
}