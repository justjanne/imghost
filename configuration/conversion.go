package configuration

import (
	"git.kuschku.de/justjanne/imghost-frontend/configuration/types"
	"github.com/justjanne/imgconv"
)

type ConversionConfiguration struct {
	TaskId        string          `json:"task_id"`
	MaxRetry      int             `json:"max_retry"`
	Timeout       types.Timeout   `json:"timeout"`
	Queue         string          `json:"queue"`
	UniqueTimeout types.Timeout   `json:"unique_timeout"`
	Quality       imgconv.Quality `json:"quality"`
	Sizes         []imgconv.Size  `json:"sizes"`
}
