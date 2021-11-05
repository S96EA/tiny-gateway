package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"sync"
)

type Config struct {
	ConfigDir string `envconfig:"CONFIG_DIR" default:""`
}

var config *Config
var once sync.Once

func Load() *Config {
	once.Do(func() {
		c := Config{}
		if err := envconfig.Process("", &c); err != nil {
			log.Panic(err)
		}
		config = &c
	})
	return config
}
