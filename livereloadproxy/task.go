package livereloadproxy

import (
	"github.com/oneut/lrp/command"
	"github.com/oneut/lrp/monitor"
)

type Task struct {
	Commands map[string]command.CommandInterface
	Monitor  monitor.MonitorInterface
}
