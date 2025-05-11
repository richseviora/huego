package config

import (
	cconfig "github.com/JeremyLoy/config"
)

type Config struct {
	HueIP  string `config:"HUE_IP"`
	HueKey string `config:"HUE_KEY"`
}

func GetConfig() *Config {
	var c Config
	err := cconfig.FromEnv().To(&c)
	if err != nil {
		panic(err)
	}
	return &c
}
