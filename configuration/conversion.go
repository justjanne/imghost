package configuration

import (
	"git.kuschku.de/justjanne/imghost-frontend/configuration/types"
	"github.com/justjanne/imgconv"
)

type ConversionConfiguration struct {
	TaskId        string          `json:"task_id" yaml:"task-id"`
	MaxRetry      int             `json:"max_retry" yaml:"max-retry"`
	Timeout       types.Timeout   `json:"timeout" yaml:"timeout"`
	Queue         string          `json:"queue" yaml:"queue"`
	UniqueTimeout types.Timeout   `json:"unique_timeout" yaml:"unique-timeout"`
	Quality       imgconv.Quality `json:"quality" yaml:"quality"`
	Sizes         []imgconv.Size  `json:"sizes" yaml:"sizes"`
}
