package config

type Monitor struct {
	Paths   []string
	Ignores []string `yaml:"ignores"`
}

func (m *Monitor) IsValid() bool {
	return len(m.Paths) != 0
}
