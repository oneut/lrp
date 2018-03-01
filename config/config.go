package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var defaultProxyHost string = ":9000"
var defaultSourceHost string = ":8080"

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

func (c *Config) GetProxyHost() string {
	if c.ProxyHost == "" {
		return defaultProxyHost
	}

	return c.ProxyHost
}

func (c *Config) GetSourceHost() string {
	if c.SourceHost == "" {
		return defaultSourceHost
	}

	return c.SourceHost
}
