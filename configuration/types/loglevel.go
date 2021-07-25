package types

import (
	"github.com/hibiken/asynq"
)

type Severity asynq.LogLevel

func (severity *Severity) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var spec string
	if err := unmarshal(&spec); err != nil {
		return err
	}

	var loglevel asynq.LogLevel
	err := loglevel.Set(spec)
	if err != nil {
		return err
	}

	*severity = Severity(loglevel)

	return nil
}
