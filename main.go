package main

import (
	"database/sql"
	"fmt"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"git.kuschku.de/justjanne/imghost-frontend/repo"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

type ImageViewModel struct {
	User  model.User
	Image model.Image
}

type AlbumViewModel struct {
	User   model.User
	Album  model.Album
	Images []model.AlbumImage
}

func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			method := r.PostFormValue("_method")
			if method == "" {
				method = r.Header.Get("X-HTTP-Method-Override")
			}

			if method == http.MethodPut ||
				method == http.MethodPatch ||
				method == http.MethodDelete {
				r.Method = method
			}
		}
		next.ServeHTTP(w, r)
	})
}

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
	albumImageRepo := repo.NewAlbumImageRepository(pageContext.Database)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Home Page"))
	})
	router.HandleFunc("/i/{imageId}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := parseUser(r)
		imageId := mux.Vars(r)["imageId"]
		image, err := imageRepo.Get(imageId)
		if err != nil {
			panic(err)
		}

		err = formatTemplate(w, "image_detail.gohtml", ImageViewModel{
			user,
			image,
		})
		if err != nil {
			panic(err)
		}
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
		if err != nil {
			panic(err)
		}
		images, err := albumImageRepo.List(album)
		if err != nil {
			panic(err)
		}

		err = formatTemplate(w, "album_detail.gohtml", AlbumViewModel{
			user,
			album,
			images,
		})
		if err != nil {
			panic(err)
		}
	}).Methods("GET")
	router.HandleFunc("/a/{albumId}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := parseUser(r)
		albumId := mux.Vars(r)["albumId"]
		album, err := albumRepo.Get(albumId)
		if err != nil { panic(err) }

		err = album.VerifyOwner(user)
		if err != nil {
			panic(err)
		}

		album.Title = r.FormValue("title")
		album.Description = r.FormValue("description")

		err = albumRepo.Update(album)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/a/"+albumId, http.StatusFound)
	}).Methods("POST")
	router.HandleFunc("/a/{albumId}/{imageId}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := parseUser(r)
		vars := mux.Vars(r)
		albumId := vars["albumId"]
		imageId := vars["imageId"]
		album, err := albumRepo.Get(albumId)
		if err != nil {
			panic(err)
		}

		err = album.VerifyOwner(user)
		if err != nil {
			panic(err)
		}

		albumImage, err := albumImageRepo.Get(album, imageId)
		if err != nil {
			panic(err)
		}

		albumImage.Title = r.FormValue("title")
		albumImage.Description = r.FormValue("description")

		err = albumImageRepo.Update(album, albumImage)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/a/"+albumId, http.StatusFound)
	}).Methods("POST")
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("postgres: %v\n", pageContext.Database.Ping())))
		_, _ = w.Write([]byte(fmt.Sprintf("redis: %v\n", pageContext.Redis.Ping().Err())))
	})
	router.HandleFunc("/{imageId}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		imageId := mux.Vars(r)["imageId"]
		image, err := imageRepo.Get(imageId)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, fmt.Sprintf("/i/%s", image.Id), http.StatusFound)
	})

	err = http.ListenAndServe(":8080", MethodOverride(router))
	if err != nil {
		panic(err)
	}
}
