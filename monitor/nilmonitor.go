package monitor

type NilMonitor struct {
}

func (nm *NilMonitor) Run(fn func(string)) {
}

func (nm *NilMonitor) Stop() {
}
