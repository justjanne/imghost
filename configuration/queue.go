package configuration

import "git.kuschku.de/justjanne/imghost-frontend/configuration/types"

type QueueConfiguration struct {
	Concurrency    int            `json:"concurrency" yaml:"concurrency"`
	LogLevel       types.Severity `json:"log_level" yaml:"log-level"`
	StrictPriority bool           `json:"strict_priority" yaml:"strict-priority"`
	Queues         map[string]int `json:"queues" yaml:"queues"`
}
