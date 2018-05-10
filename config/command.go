package config

type Command struct {
	Executes     []string
	NeedsRestart bool     `yaml:"needsRestart"`
	WatchStdouts []string `yaml:"watchStdouts"`
}

func (c *Command) IsValid() bool {
	for _, execute := range c.Executes {
		if execute != "" {
			return true
		}
	}

	return false
}
