package types

import "time"

type Timeout time.Duration

func (timeout *Timeout) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var spec string
	if err := unmarshal(&spec); err != nil {
		return err
	}

	duration, err := time.ParseDuration(spec)
	if err != nil {
		return err
	}

	*timeout = Timeout(duration)

	return nil
}
