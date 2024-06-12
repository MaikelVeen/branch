package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var (
	configuration *viper.Viper
)

const (
	KeyPattern            = "pattern"
	defaultConfigFilename = "config"
	envPrefix             = "BRANCH"
)

// Config represents the configuration of the application.
type Config struct {
	Pattern *string `yaml:"pattern"`
}

// ConfigOption represents a configuration option that can be displayed to the user.
type ConfigOption struct {
	Key          string
	Description  string
	CurrentValue func(cfg Config) *string
}

// Load loads the configuration from the environment.
func Load() (*Config, error) {
	if configuration == nil {
		return nil, errors.New("configuration not initialized")
	}

	cfg := Config{}
	if err := configuration.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}

// Options is a list of all available configuration options.
var Options = map[string]*ConfigOption{
	KeyPattern: {
		Key:          KeyPattern,
		Description:  "The pattern to use for branch names",
		CurrentValue: func(cfg Config) *string { return cfg.Pattern },
	},
}

func Init() (*viper.Viper, error) {
	fmt.Println("initializeConfig")
	v := viper.New()

	v.SetConfigName(defaultConfigFilename)
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME/.branch/")
	v.AddConfigPath(".")

	var cfgNotFoundError viper.ConfigFileNotFoundError
	if err := v.ReadInConfig(); err != nil {
		if !errors.As(err, &cfgNotFoundError) {
			return nil, err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	configuration = v
	return v, nil
}
