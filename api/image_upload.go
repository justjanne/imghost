package api

import (
	"context"
	"git.kuschku.de/justjanne/imghost-frontend/auth"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"net/http"
)

func UploadImage(env environment.FrontendEnvironment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user, err := auth.ParseUser(request, env)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		var image model.Image
		image.Id = "testid"
		image.Owner = user.Id
		err = env.Repositories.Images.Create(image)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		err = env.Storage.Upload(
			context.Background(),
			env.Configuration.Storage.ConversionBucket,
			image.Id,
			request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		util.ReturnJson(writer, image)
	})
}
