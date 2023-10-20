package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type (
	ServerConfiguration struct {
		BrokerAddr string
		BrokerPort string
	}
)

// NewConfig initialize configuration
func NewConfig(path string) (*ServerConfiguration, error) {
	var configuration *ServerConfiguration

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("error reading config file")
		return nil, err
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to decode into struct")
		return nil, err
	}

	return configuration, nil
}
