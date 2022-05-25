package main

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost/shared"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/x/metrics"
	"github.com/hibiken/asynqmon"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	defer shared.ErrorHandler()

	configFile, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("error opening config file: %s", err.Error())
	}
	config := shared.LoadConfigFromFile(configFile)

	db, err := sql.Open(config.Database.Format, config.Database.Url)
	if err != nil {
		log.Fatalf("error connecting to database: %s", err.Error())
	}

	if err := Migrate(db); err != nil {
		log.Fatalf("error migrating database: %s", err.Error())
	}

	env := PageEnvironment{
		&config,
		asynq.NewClient(config.AsynqOpts()),
		asynq.NewInspector(config.AsynqOpts()),
		config.UploadTimeoutDuration(),
		db,
		http.FileServer(http.Dir(config.TargetFolder)),
		http.FileServer(http.Dir("assets")),
	}

	monitor := asynqmon.New(asynqmon.Options{
		RootPath:         "/admin",
		RedisConnOpt:     config.AsynqOpts(),
		PayloadFormatter: asynqmon.PayloadFormatterFunc(shared.FormatPayload),
		ResultFormatter:  asynqmon.ResultFormatterFunc(shared.FormatResult),
	})
	http.Handle(monitor.RootPath()+"/static/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))
			monitor.ServeHTTP(w, r)
		}),
	)
	http.Handle(monitor.RootPath()+"/", monitor)

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(metrics.NewQueueMetricsCollector(asynq.NewInspector(config.AsynqOpts())))
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	http.Handle("/upload/", pageUpload(env))

	http.Handle("/i/", http.StripPrefix("/i/", pageImageDetail(env)))
	http.Handle("/a/", http.StripPrefix("/a/", pageAlbumDetail(env)))

	http.Handle("/me/images/", pageImageList(env))
	http.Handle("/assets/", http.StripPrefix("/assets/", env.AssetServer))
	http.Handle("/", pageIndex(env))

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error in http server: %s", err.Error())
	}
}
