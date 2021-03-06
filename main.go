package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mainflux/mainflux-auth/api"
	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/config"
	"github.com/mainflux/mainflux-auth/domain"
)

const (
	defaultConfig string = "/src/github.com/mainflux/mainflux-auth/config/default.toml"
	httpPort      string = ":8180"
	help          string = `
		Usage: mainflux-auth [options]
		Options:
			-c, --config <file>         Configuration file
			-h, --help                  Prints this message end exits`
)

type options struct {
	Config string
	Help   bool
}

func main() {
	opts := options{}
	flag.StringVar(&opts.Config, "c", "", "Configuration file.")
	flag.StringVar(&opts.Config, "config", "", "Configuration file.")
	flag.BoolVar(&opts.Help, "h", false, "Show help.")
	flag.BoolVar(&opts.Help, "help", false, "Show help.")

	flag.Parse()

	if opts.Help {
		fmt.Printf("%s\n", help)
		os.Exit(0)
	}

	if opts.Config == "" {
		opts.Config = os.Getenv("GOPATH") + defaultConfig
	}

	cfg := config.Config{}
	cfg.Load(opts.Config)

	if cfg.SecretKey != "" {
		domain.SetSecretKey(cfg.SecretKey)
	}

	cache.Start(cfg.CacheURL())
	defer cache.Stop()

	http.ListenAndServe(httpPort, api.Server())
}
