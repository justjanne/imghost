package main

import (
	"context"
	"database/sql"
	"git.kuschku.de/justjanne/imghost/shared"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	configFile, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("error opening config file: %s", err.Error())
	}
	config := shared.LoadConfigFromFile(configFile)

	db, err := sql.Open(config.Database.Format, config.Database.Url)
	if err != nil {
		log.Fatalf("error connecting to database: %s", err.Error())
	}

	pageContext := PageContext{
		context.Background(),
		&config,
		asynq.NewClient(config.AsynqOpts()),
		config.UploadTimeoutDuration(),
		db,
		http.FileServer(http.Dir(config.TargetFolder)),
		http.FileServer(http.Dir("assets")),
	}

	http.Handle("/upload/", pageUpload(pageContext))

	http.Handle("/i/", http.StripPrefix("/i/", pageImageDetail(pageContext)))
	http.Handle("/a/", http.StripPrefix("/a/", pageAlbumDetail(pageContext)))

	http.Handle("/me/images/", pageImageList(pageContext))
	http.Handle("/assets/", http.StripPrefix("/assets/", pageContext.AssetServer))
	http.Handle("/", pageIndex(pageContext))

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error in http server: %s", err.Error())
	}
}
