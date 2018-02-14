package monitor

type Task struct {
	MonitorPath string `yaml:"monitor_path"`
	Command     string
	Compile     bool
}
