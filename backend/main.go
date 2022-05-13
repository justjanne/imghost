package main

import (
	"context"
	"git.kuschku.de/justjanne/imghost-frontend/shared"
	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/gographics/imagick.v3/imagick"
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
	configFile, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("error opening config file: %s", err.Error())
	}
	config := shared.LoadConfigFromFile(configFile)

	imagick.Initialize()
	defer imagick.Terminate()

	srv := asynq.NewServer(
		config.AsynqOpts(),
		asynq.Config{Concurrency: config.Concurrency},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(shared.TypeImageResize, ProcessImageHandler(&config))

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})
	metrics := &http.Server{
		Addr:    ":2112",
		Handler: metricsMux,
	}

	go func() {
		if err := metrics.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error in metrics server: %s", err.Error())
		}
		srv.Shutdown()
	}()

	if err := srv.Run(mux); err != nil {
		log.Printf("error in asynq server: %s", err.Error())
	}
	if err := metrics.Shutdown(context.Background()); err != nil {
		log.Printf("error shutting down metrics server: %s", err.Error())
	}
}
