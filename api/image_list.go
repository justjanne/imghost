package api

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost-frontend/auth"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"net/http"
)

func ListImages(env environment.Environment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user, err := auth.ParseUser(request, env)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
		}
		images, err := env.Repositories.Images.List(user)
		if err == sql.ErrNoRows {
			http.NotFound(writer, request)
			return
		} else if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		util.ReturnJson(writer, images)
	})
}