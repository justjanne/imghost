package api

import (
	"git.kuschku.de/justjanne/imghost-frontend/auth"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"net/http"
)

func ListImages(env environment.FrontendEnvironment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user, err := auth.ParseUser(request, env)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
		}
		images, err := env.Repositories.Images.List(user)
		var infos []model.ImageInfo
		for _, image := range images {
			info, err := EnrichImageInfo(env, image)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			infos = append(infos, info)
		}

		util.ReturnJson(writer, infos)
	})
}
