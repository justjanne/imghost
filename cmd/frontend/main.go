package main

import (
	"git.kuschku.de/justjanne/imghost-frontend/api"
	"git.kuschku.de/justjanne/imghost-frontend/configuration"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
)

func main() {
	var config configuration.Configuration
	configFile, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		panic(err)
	}

	env, err := environment.NewClientEnvironment(config)
	if err != nil {
		panic(err)
	}
	defer env.Destroy()

	router := mux.NewRouter()
	// Image API
	router.Handle(
		"/api/v1/images",
		api.ListImages(env)).Methods(http.MethodGet)
	router.Handle(
		"/api/v1/images/{imageId}",
		api.GetImage(env)).Methods(http.MethodGet)

	// Album API
	router.Handle(
		"/api/v1/albums",
		api.ListAlbums(env)).Methods(http.MethodGet)
	router.Handle(
		"/api/v1/albums/{imageId}",
		api.GetAlbum(env)).Methods(http.MethodGet)

	// Album Image API
	router.Handle(
		"/api/v1/albums/{albumId}/images",
		api.ListAlbumImages(env)).Methods(http.MethodGet)
	router.Handle(
		"/api/v1/albums/{albumId}/images/{imageId}",
		api.GetAlbumImage(env)).Methods(http.MethodGet)

	// TODO: Implement mutating API methods

	if err = http.ListenAndServe(":8080", util.MethodOverride(router)); err != nil {
		panic(err)
	}
}
