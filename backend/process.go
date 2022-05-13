package main

import (
	"context"
	"encoding/json"
	"fmt"
	"git.kuschku.de/justjanne/imghost/shared"
	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func trackTimeSince(counter prometheus.Counter, start time.Time) time.Time {
	now := time.Now().UTC()
	counter.Add(float64(now.Sub(start).Milliseconds()) / 1000.0)
	return now
}

func ProcessImageHandler(config *shared.Config) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		log.Printf("received image resize task")
		task := shared.ImageTaskPayload{}
		if err := json.Unmarshal(t.Payload(), &task); err != nil {
			return err
		}

		log.Printf("starting image resize task %s", task.ImageId)
		errors := ResizeImage(config, task.ImageId)
		log.Printf("deleting cached image for image resize task %s", task.ImageId)
		_ = os.Remove(filepath.Join(config.SourceFolder, task.ImageId))

		errorMessages := make([]string, len(errors))
		for i, err := range errors {
			errorMessages[i] = err.Error()
		}

		if len(errors) != 0 {
			log.Printf("errors occured while processing image resize task %s: %s", task.ImageId, strings.Join(errorMessages, "\n"))
			if err := json.NewEncoder(t.ResultWriter()).Encode(shared.Result{
				Id:      task.ImageId,
				Success: true,
				Errors:  errorMessages,
			}); err != nil {
				return err
			}
			return fmt.Errorf(
				"errors occured while processing task %s (%s): %s",
				t.Type(),
				task.ImageId,
				strings.Join(errorMessages, "\n"),
			)
		}

		if err := json.NewEncoder(t.ResultWriter()).Encode(shared.Result{
			Id:      task.ImageId,
			Success: true,
			Errors:  []string{},
		}); err != nil {
			return err
		}
		return nil
	}
}
