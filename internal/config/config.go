package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `toml:"env"`
	HTTPServer  `toml:"http_server"`
	BaseUrl     string `toml:"base_url"`
	Storage     string `toml:"storage"`
	AliasLength int    `toml:"aliasLength"`
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

	const op = "internal.config.config.go.MustLoad"

	flag.StringVar(&configPath, "path", "configs/api.toml", "config path")
	flag.Parse()
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("op: %v \nconfig file does not exists: %s", op, configPath)
	}

	var cfg Config

	_, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Printf("op: %v \ncan not decode config file, %v", op, err)
		os.Exit(1)
	}

	return &cfg
}
