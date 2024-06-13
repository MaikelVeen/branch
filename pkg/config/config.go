package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/spf13/viper"
)

var (
	configuration *viper.Viper
)

const (
	KeyPattern = "pattern"

	defaultConfigFilename = "config"
	path                  = "$HOME/.config/branch/"
	envPrefix             = "BRANCH"
)

// Config represents the configuration of the application.
type Config struct {
	Pattern *string `yaml:"pattern"`
}

func (c *Config) Save() error {
	return configuration.WriteConfig()
}

// Option represents a configuration option that can be displayed to the user.
type Option struct {
	Key          string
	Description  string
	CurrentValue func(cfg Config) *string
	SetValue     func(cfg *Config, value string) error
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
var Options = map[string]*Option{
	KeyPattern: {
		Key:          KeyPattern,
		Description:  "The pattern to use for branch names",
		CurrentValue: func(cfg Config) *string { return cfg.Pattern },
		SetValue: func(cfg *Config, value string) error {
			cfg.Pattern = &value
			configuration.Set(KeyPattern, value)
			return nil
		},
	},
}

func Init() (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(defaultConfigFilename)
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if errors.As(err, &e) {
			// If the configuration file is not found, create it.
			if err = createDefaultConfigFile(); err != nil {
				return nil, fmt.Errorf("failed to write default configuration: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read configuration: %w", err)
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	configuration = v
	return v, nil
}

func createDefaultConfigFile() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/.config/branch", usr.HomeDir)
	if err = os.MkdirAll(path, 0755); err != nil {
		return err
	}

	fn := fmt.Sprintf("%s/%s.yaml", path, "defaultConfigFilename")
	return os.WriteFile(fn, []byte(""), 0600)
}
