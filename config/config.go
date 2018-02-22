package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func GetConfig() *Config {
	buf, err := ioutil.ReadFile("lrp.yml")
	if err != nil {
		panic(err)
	}

	config := &Config{}
	yaml.Unmarshal(buf, config)

	return config
}

type Config struct {
	ProxyHost  string `yaml:"proxy_host"`
	SourceHost string `yaml:"source_host"`
	Tasks      map[string]Task
}
