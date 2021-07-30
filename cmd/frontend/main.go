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
	var config configuration.FrontendConfiguration
	configFile, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		panic(err)
	}

	env, err := environment.NewFrontendEnvironment(config)
	if err != nil {
		panic(err)
	}
	defer env.Destroy()

	router := mux.NewRouter()
	// Image API
	router.Handle(
		"/api/v1/images",
		api.ListImages(env)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle(
		"/api/v1/images",
		api.UploadImage(env)).Methods(http.MethodPost, http.MethodOptions)
	router.Handle(
		"/api/v1/images/{imageId}",
		api.GetImage(env)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle(
		"/api/v1/images/{imageId}",
		api.UpdateImage(env)).Methods(http.MethodPost, http.MethodOptions)
	router.Handle(
		"/api/v1/images/{imageId}",
		api.DeleteImage(env)).Methods(http.MethodDelete, http.MethodOptions)

	// Album API
	router.Handle(
		"/api/v1/albums",
		api.ListAlbums(env)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle(
		"/api/v1/albums/{albumId}",
		api.GetAlbum(env)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle(
		"/api/v1/albums/{albumId}",
		api.UpdateAlbum(env)).Methods(http.MethodPost, http.MethodOptions)
	router.Handle(
		"/api/v1/albums/{albumId}/reorder",
		api.ReorderAlbum(env)).Methods(http.MethodPost, http.MethodOptions)
	router.Handle(
		"/api/v1/albums/{albumId}",
		api.DeleteAlbum(env)).Methods(http.MethodDelete, http.MethodOptions)

	// Album Image API
	router.Handle(
		"/api/v1/albums/{albumId}/images",
		api.ListAlbumImages(env)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle(
		"/api/v1/albums/{albumId}/images/{imageId}",
		api.GetAlbumImage(env)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle(
		"/api/v1/albums/{albumId}/images/{imageId}",
		api.UpdateAlbumImage(env)).Methods(http.MethodPost, http.MethodOptions)
	router.Handle(
		"/api/v1/albums/{albumId}/images/{imageId}",
		api.DeleteAlbumImage(env)).Methods(http.MethodDelete, http.MethodOptions)

	if err = http.ListenAndServe(":8080", util.CorsWrapper(router)); err != nil {
		panic(err)
	}
}
