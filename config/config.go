package config

import (
	"gopkg.in/yaml.v2"
)

func CreateConfig(yamlBuf []byte) Config {
	config := Config{}
	yaml.Unmarshal(yamlBuf, &config)
	return config
}

type Config struct {
	Proxy  Proxy
	Source Source
	Tasks  map[string]Task
}
