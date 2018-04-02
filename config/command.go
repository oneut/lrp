package config

type Command struct {
	Executes     []string
	NeedsRestart bool     `yaml:"needs_restart"`
	WatchStdouts []string `yaml:"watch_stdouts"`
}

func (c *Command) IsValid() bool {
	for _, execute := range c.Executes {
		if execute != "" {
			return true
		}
	}

	return false
}
