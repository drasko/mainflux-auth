package config

import (
	"fmt"
	"log"
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
		log.Fatalf("Cannot load config due to %s", err)
		os.Exit(1)
	}
}

func (cfg *Config) CacheURL() string {
	return fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
}
