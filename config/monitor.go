package config

type Monitor struct {
	Paths        []string
	ExcludePaths []string `yaml:"exclude_paths"`
}
