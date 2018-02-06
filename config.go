package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/oneut/llp/monitor"
)

func getConfig() *Config {
	buf, err := ioutil.ReadFile("llp.yml")
	if err != nil {
		panic(err)
	}

	config := &Config{}
	yaml.Unmarshal(buf, config)

	return config
}

type Config struct {
	Port  int
	Tasks map[string]monitor.Task
}
