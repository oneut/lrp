package config

type Command struct {
	Execute      string
	NeedsRestart bool     `yaml:"needs_restart"`
	WatchStdouts []string `yaml:"watch_stdouts"`
}

func (c *Command) IsValid() bool {
	return c.Execute != ""
}
