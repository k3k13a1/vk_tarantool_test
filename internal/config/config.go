package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Srv       Server    `yaml:"srv"`
	Tarantool Tarantool `yaml:"tarantool"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Tarantool struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

func SetupConfig() *Config {
	var config Config

	yamlFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		panic("Can't read config file")
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic("Can't parse config file")
	}

	return &config
}
