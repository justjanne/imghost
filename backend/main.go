package main

import (
	"context"
	"database/sql"
	"git.kuschku.de/justjanne/imghost/shared"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/gographics/imagick.v2/imagick"
	"log"
	"net/http"
	"os"
)

var imageProcessDuration = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "imghost_process_duration",
	Help: "The amount of time spent processing images",
}, []string{"task"})
var imageProcessDurationRead = imageProcessDuration.WithLabelValues("read")
var imageProcessDurationClone = imageProcessDuration.WithLabelValues("clone")
var imageProcessDurationCrop = imageProcessDuration.WithLabelValues("crop")
var imageProcessDurationResize = imageProcessDuration.WithLabelValues("resize")
var imageProcessDurationWrite = imageProcessDuration.WithLabelValues("write")

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

	env := ProcessingEnvironment{
		Config:   &config,
		Database: db,
	}

	imagick.Initialize()
	defer imagick.Terminate()

	srv := asynq.NewServer(
		config.AsynqOpts(),
		asynq.Config{Concurrency: config.Concurrency},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(shared.TypeImageResize, ProcessImageHandler(env))

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})
	metrics := &http.Server{
		Addr:    ":2112",
		Handler: metricsMux,
	}

	runner := shared.Runner{}
	runner.RunParallel(func() {
		log.Printf("starting metrics server")
		if err := metrics.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error in metrics server: %s", err.Error())
		}
		log.Printf("metrics shut down, shutting down asynq as well")
		srv.Shutdown()
	})
	runner.RunParallel(func() {
		log.Printf("starting asynq server")
		if err := srv.Run(mux); err != nil {
			log.Printf("error in asynq server: %s", err.Error())
		}
		log.Printf("asynq shut down, shutting down metrics as well")
		if err := metrics.Shutdown(context.Background()); err != nil {
			log.Printf("error shutting down metrics server: %s", err.Error())
		}
	})
	runner.Wait()
}
