package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("loads config from env", func(t *testing.T) {
		viper.Reset()
		// Set the environment variable
		dsn := "postgres://user:pass@localhost:5432/db"
		os.Setenv("DATABASE_DSN", dsn)
		defer os.Unsetenv("DATABASE_DSN")

		// Load the config
		config, err := LoadConfig(".")
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, dsn, config.Database.DSN)
	})

	t.Run("loads config from file", func(t *testing.T) {
		viper.Reset()
		// Create a temporary directory
		tmpDir, err := os.MkdirTemp("", "config-test")
		assert.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Create a config file in the temp directory
		content := []byte("database:\n  dsn: \"file_dsn\"")
		// Viper looks for config.yaml by default, so we name the file config.yaml
		err = os.WriteFile(filepath.Join(tmpDir, "config.yaml"), content, 0600)
		assert.NoError(t, err)

		// Load the config from the temp directory
		config, err := LoadConfig(tmpDir)
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "file_dsn", config.Database.DSN)
	})
}
