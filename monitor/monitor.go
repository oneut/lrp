package monitor

import (
	"github.com/oneut/lrp/config"
)

func NewMonitor(name string, monitorConfig config.Monitor) MonitorInterface {
	if !(monitorConfig.IsValid()) {
		return &NilMonitor{}
	}

	return NewFsnotifyMonitor(name, monitorConfig)
}

type MonitorInterface interface {
	Run(func(string))
	Stop()
}
