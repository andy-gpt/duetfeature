package main

import (
	"duetfeature/config"
	"duetfeature/internal/app"
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
)

func main() {
	// Config discovery
	configPath := flag.String("config", "", "Config Path")
	flag.Parse()
	if len(*configPath) == 0 {
		log.Fatal().Msg(fmt.Sprintf("path to a config file"))
	}

	// Configuration
	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to decode into struct")
	}

	// Run
	app.Run(cfg)

}
