package config

type Command struct {
	Execute      string
	NeedsRestart bool     `yaml:"needs_restart"`
	WatchStdout  []string `yaml:"watch_stdout"`
}
