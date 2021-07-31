package api

import (
	"git.kuschku.de/justjanne/imghost-frontend/auth"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
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
		for idx, image := range images {
			err = image.LoadUrl(env.Storage, env.Configuration.Storage)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			image.Metadata, err = env.Repositories.ImageMetadata.List(image)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			images[idx] = image
		}

		util.ReturnJson(writer, images)
	})
}
