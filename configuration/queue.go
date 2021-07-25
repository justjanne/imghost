package configuration

import "git.kuschku.de/justjanne/imghost-frontend/configuration/types"

type QueueConfiguration struct {
	Concurrency    int            `json:"concurrency"`
	LogLevel       types.Severity `json:"log_level"`
	StrictPriority bool           `json:"strict_priority"`
	Queues         map[string]int `json:"queues"`
}
