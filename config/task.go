package config

import (
	"time"
)

var defaultAggregateTimeout = 300

type Task struct {
	Commands         map[string]Command
	Monitor          Monitor
	AggregateTimeout int `yaml:"aggregate_timeout"`
}

func (t *Task) GetAggregateTimeout() time.Duration {
	if t.AggregateTimeout > 0 {
		return time.Duration(t.AggregateTimeout) * time.Millisecond
	}

	return time.Duration(defaultAggregateTimeout) * time.Millisecond
}
