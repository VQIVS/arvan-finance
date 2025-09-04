package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadConfig(configPath string) (Config, error) {
	var c Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return c, err
	}

	return c, yaml.Unmarshal(data, &c)
}

func MustReadConfig(configPath string) Config {
	c, err := ReadConfig(configPath)
	if err != nil {
		panic(err)
	}
	return c
}
