// Package config provides configuration loading utilities.
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents the configurable service parameters.
type Config struct {
	RedisHost string `toml:"redisHost"`
	RedisPort int    `toml:"redisPort"`
	SecretKey string `toml:"secretKey"`
}

// Load loads TOML file contents. If decoding fails, the service is aborted.
func (cfg *Config) Load(file string) {
	if _, err := toml.DecodeFile(file, &cfg); err != nil {
		fmt.Printf("Cannot load config due to %s", err)
		os.Exit(1)
	}
}

// CacheURL retrieves an URL of a redis instance that will be used by the
// service.
func (cfg *Config) CacheURL() string {
	return fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
}
