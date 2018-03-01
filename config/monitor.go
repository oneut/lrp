package config

type Monitor struct {
	Paths  []string
	Ignore []string `yaml:"ignore"`
}

func (m *Monitor) IsValid() bool {
	return len(m.Paths) != 0
}
