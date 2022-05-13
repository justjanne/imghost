package main

import (
	"context"
	"encoding/json"
	"fmt"
	"git.kuschku.de/justjanne/imghost/shared"
	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus"
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
		task := shared.ImageTaskPayload{}
		if err := json.Unmarshal(t.Payload(), &task); err != nil {
			return err
		}

		errors := ResizeImage(config, task.ImageId)
		_ = os.Remove(filepath.Join(config.SourceFolder, task.ImageId))

		errorMessages := make([]string, len(errors))
		for i, err := range errors {
			errorMessages[i] = err.Error()
		}

		if len(errors) != 0 {
			return fmt.Errorf(
				"errors occured while processing task %s (%s): %s",
				t.Type(),
				task.ImageId,
				strings.Join(errorMessages, "\n"),
			)
		}

		return nil
	}
}
