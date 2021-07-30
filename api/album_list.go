package api

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost-frontend/auth"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"net/http"
)

func ListAlbums(env environment.FrontendEnvironment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user, err := auth.ParseUser(request, env)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
		}
		albums, err := env.Repositories.Albums.List(user)
		if err == sql.ErrNoRows {
			http.NotFound(writer, request)
			return
		} else if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		for i, album := range albums {
			album.Images, err = env.Repositories.AlbumImages.List(album.Id)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			for j, image := range album.Images {
				err = image.LoadUrl(env.Storage, env.Configuration.Storage)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
					return
				}
				album.Images[j] = image
			}
			albums[i] = album
		}

		util.ReturnJson(writer, albums)
	})
}
