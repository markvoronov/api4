package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"time"
)

type Config struct {
	Env        string `toml:"env"`
	HTTPServer `toml:"http_server"`
}

type HTTPServer struct {
	Address     string        `toml:"address" env-default:"localhost:8083"`
	Timeout     time.Duration `toml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `toml:"idle_timeout" env-default:"60s"`
}

var (
	configPath string
)

func MustLoad() *Config {

	flag.StringVar(&configPath, "path", "configs/api.toml", "config path")
	flag.Parse()
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	var cfg Config

	_, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Println("can not find config file", err)
	}

	return &cfg
}
