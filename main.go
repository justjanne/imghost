package main

import (
	"database/sql"
	"fmt"
	"git.kuschku.de/justjanne/imghost-frontend/repo"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	config := NewConfigFromEnv()

	db, err := sql.Open(config.Database.Format, config.Database.Url)
	if err != nil {
		panic(err)
	}

	pageContext := PageContext{
		&config,
		redis.NewClient(&redis.Options{
			Addr:     config.Redis.Address,
			Password: config.Redis.Password,
		}),
		db,
		http.FileServer(http.Dir(config.TargetFolder)),
		http.FileServer(http.Dir("assets")),
	}

	imageRepo := repo.NewImageRepository(pageContext.Database)
	albumRepo := repo.NewAlbumRepository(pageContext.Database)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Home Page"))
	})
	router.HandleFunc("/i/{imageId}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := parseUser(r)
		imageId := mux.Vars(r)["imageId"]
		image, err := imageRepo.Get(imageId)
		if err != nil { panic(err) }

		err = formatTemplate(w, "image_detail.html", ImageDetailData{
			user,
			image,
		})
		if err != nil { panic(err) }
	}).Methods("GET")
	router.HandleFunc("/i/{imageId}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := parseUser(r)
		imageId := mux.Vars(r)["imageId"]
		image, err := imageRepo.Get(imageId)
		if err != nil { panic(err) }

		err = image.VerifyOwner(user)
		if err != nil { panic(err) }

		image.Title = r.FormValue("title")
		image.Description = r.FormValue("description")

		err = imageRepo.Update(image)
		if err != nil { panic(err) }

		http.Redirect(w, r, "/i/" + imageId, http.StatusFound)
	}).Methods("POST")
	router.HandleFunc("/a/{albumId}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := parseUser(r)
		albumId := mux.Vars(r)["albumId"]
		album, err := albumRepo.Get(albumId)
		if err != nil { panic(err) }
		images, err := albumRepo.GetImages(album)
		if err != nil { panic(err) }

		err = formatTemplate(w, "album_detail.html", AlbumDetailData{
			user,
			album,
			images,
		})
		if err != nil { panic(err) }
	}).Methods("GET")
	router.HandleFunc("/a/{albumId}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := parseUser(r)
		albumId := mux.Vars(r)["albumId"]
		album, err := albumRepo.Get(albumId)
		if err != nil { panic(err) }

		err = album.VerifyOwner(user)
		if err != nil { panic(err) }

		album.Title = r.FormValue("title")
		album.Description = r.FormValue("description")

		err = albumRepo.Update(album)
		if err != nil { panic(err) }

		http.Redirect(w, r, "/a/" + albumId, http.StatusFound)
	}).Methods("POST")
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("postgres: %v\n", pageContext.Database.Ping())))
		_, _ = w.Write([]byte(fmt.Sprintf("redis: %v\n", pageContext.Redis.Ping().Err())))
	})
	router.HandleFunc("/{imageId}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := parseUser(r)
		imageId := mux.Vars(r)["imageId"]
		image, err := imageRepo.Get(imageId)
		if err != nil { panic(err) }

		err = formatTemplate(w, "image_detail.html", ImageDetailData{
			user,
			image,
		})
		if err != nil { panic(err) }
	})

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
