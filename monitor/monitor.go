package monitor

import (
	"github.com/oneut/lrp/config"
)

func NewMonitor(name string, monitorConfig config.Monitor) Monitorer {
	if !(monitorConfig.IsValid()) {
		return &NilMonitor{}
	}

	return NewFsnotifyMonitor(name, monitorConfig)
}

type Monitorer interface {
	Run(func(string))
	Stop()
}
