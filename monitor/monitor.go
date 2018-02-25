package monitor

import (
	"github.com/oneut/lrp/config"
)

func NewMonitor(name string, monitorConfig config.Monitor) Monitorer {
	return NewFsnotifyMonitor(name, monitorConfig)
}

type Monitorer interface {
	Run(fn func(string))
	Stop()
}
