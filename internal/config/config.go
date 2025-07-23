package config

import (
	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	Database struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"database"`
}

// LoadConfig loads the configuration from the given path.
func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	// Bind the DATABASE_DSN environment variable to the config struct
	viper.BindEnv("database.dsn", "DATABASE_DSN")

	if err := viper.ReadInConfig(); err != nil {
		// If running in Lambda, we might not have a config file, which is fine.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
