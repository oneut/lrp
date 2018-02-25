package config

type Monitor struct {
	Paths  []string
	Ignore []string `yaml:"ignore"`
}
