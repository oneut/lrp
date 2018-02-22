package monitor

import (
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/log"
)

func NewMonitor(name string, monitorConfig config.Monitor) Monitorer {
	log.Info(name, "Start Monitor")
	return NewFsnotifyMonitor(name, monitorConfig)
}

type Monitorer interface {
	Run(fn func(string))
}
