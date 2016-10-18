package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RedisHost string `toml:"redisHost"`
	RedisPort int    `toml:"redisPort"`
	SecretKey string `toml:"secretKey"`
}

func (cfg *Config) Load(file string) {
	if _, err := toml.DecodeFile(file, &cfg); err != nil {
		fmt.Println("Cannot process TOML file.")
		os.Exit(1)
	}
}
