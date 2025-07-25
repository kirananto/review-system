package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	t.Run("creates logger with info level by default", func(t *testing.T) {
		config := &LogConfig{}
		logger := NewLogger(config)
		assert.NotNil(t, logger)
		assert.Equal(t, zerolog.InfoLevel, logger.Logger.GetLevel())
	})

	t.Run("creates logger with specified level", func(t *testing.T) {
		config := &LogConfig{LogLevel: "debug"}
		logger := NewLogger(config)
		assert.NotNil(t, logger)
		assert.Equal(t, zerolog.DebugLevel, logger.Logger.GetLevel())
	})
}

func TestLogLevels(t *testing.T) {
	levels := map[string]zerolog.Level{
		"trace": zerolog.TraceLevel,
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"panic": zerolog.PanicLevel,
	}

	for levelName, level := range levels {
		t.Run(levelName, func(t *testing.T) {
			assert.Equal(t, level, getLogLevel(levelName))
		})
	}
}

func TestLogMethods(t *testing.T) {
	var buf bytes.Buffer
	config := &LogConfig{LogLevel: "trace"}
	logger := NewLogger(config)
	newLogger := logger.Logger.Output(&buf)
	logger.Logger = &newLogger

	t.Run("Info", func(t *testing.T) {
		buf.Reset()
		logger.Info("test info")
		assert.Contains(t, buf.String(), `"M":"test info"`)
		assert.Contains(t, buf.String(), `"L":"info"`)
	})

	t.Run("Warn", func(t *testing.T) {
		buf.Reset()
		logger.Warn(errors.New("test warn error"), "test warn")
		assert.Contains(t, buf.String(), `"M":"test warn"`)
		assert.Contains(t, buf.String(), `"L":"warn"`)
		assert.Contains(t, buf.String(), `"E":"test warn error"`)
	})

	t.Run("Error", func(t *testing.T) {
		buf.Reset()
		logger.Error(errors.New("test error"), "test error message")
		var log map[string]interface{}
		err := json.Unmarshal(buf.Bytes(), &log)
		assert.NoError(t, err)
		assert.Equal(t, "error", log["L"])
		assert.Equal(t, "test error message", log["M"])
		assert.Equal(t, "test error", log["E"])
	})
}

func TestParserFileName(t *testing.T) {
	workingDir, _ := os.Getwd()
	fileName := parserFileName(workingDir, "/app/internal/logger/logger.go", 25)
	assert.Equal(t, strings.Replace("/app/internal/logger/logger.go", workingDir, "", 1)+":25", fileName)
}
