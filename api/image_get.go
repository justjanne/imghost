package api

import (
	"context"
	"database/sql"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"github.com/gorilla/mux"
	"net/http"
)

func EnrichImageInfo(env environment.FrontendEnvironment, image model.Image) (info model.ImageInfo, err error) {
	info.Image = image
	info.State, err = env.Repositories.ImageStates.Get(image.Id)
	if err != nil {
		return
	}
	imageUrl, err := env.Storage.UrlFor(context.Background(), env.Configuration.Storage.ImageBucket, info.Image.Id)
	if err != nil {
		return
	}
	info.Url = imageUrl.String()
	return
}

func GetImage(env environment.FrontendEnvironment) http.Handler {
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
		info, err := EnrichImageInfo(env, image)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		util.ReturnJson(writer, info)
	})
}
