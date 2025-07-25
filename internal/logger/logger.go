package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

type LogConfig struct {
	LogLevel string `env:"LOG_LEVEL,required"`
}

type Logger struct {
	Logger  *zerolog.Logger
	baseDir string
}

func NewLogger(config *LogConfig) *Logger {
	logLevel := zerolog.InfoLevel
	if config.LogLevel != "" {
		logLevel = getLogLevel(config.LogLevel)
	}

	zerolog.TimestampFieldName = "T"
	zerolog.LevelFieldName = "L"
	zerolog.MessageFieldName = "M"
	zerolog.ErrorFieldName = "E"

	var writers []io.Writer

	writers = append(writers, zerolog.ConsoleWriter{
		Out: os.Stdout,
	})

	logWriters := io.MultiWriter(writers...)

	logger := zerolog.New(logWriters).
		Level(logLevel).
		With().
		Str("P", "review-system-svc").
		Str("C", "review-system").
		Timestamp().Logger()

	log.Println("Successfully configured logger")

	workingDir, err := os.Getwd()
	if err != nil {
		log.Panic(fmt.Sprintf("Failure in getting the current working directory: %v", err))
		panic("Failure in getting the current working directory")
	}

	return &Logger{
		Logger:  &logger,
		baseDir: workingDir,
	}
}

func getLogLevel(logLevel string) zerolog.Level {
	switch strings.ToLower(logLevel) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func parserFileName(baseDir string, file string, line int) string {
	filename := strings.Replace(file, baseDir, "", 1)
	return fmt.Sprintf("%s:%v", filename, line)
}

func (l *Logger) Trace(message string) {
	_, file, line, ok := runtime.Caller(1)

	if ok {
		l.Logger.Trace().Str("f", parserFileName(l.baseDir, file, line)).Msg(message)
	} else {
		l.Logger.Trace().Msg(message)
	}
}

func (l *Logger) Debug(message string) {
	_, file, line, ok := runtime.Caller(1)

	if ok {
		l.Logger.Debug().Str("f", parserFileName(l.baseDir, file, line)).Msg(message)
	} else {
		l.Logger.Debug().Msg(message)
	}
}

func (l *Logger) Info(message string) {
	_, file, line, ok := runtime.Caller(1)

	if ok {
		l.Logger.Info().Str("f", parserFileName(l.baseDir, file, line)).Msg(message)
	} else {
		l.Logger.Info().Msg(message)
	}
}

func (l *Logger) Warn(err error, message string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		l.Logger.Warn().Str("f", parserFileName(l.baseDir, file, line)).Err(err).Msg(message)
	} else {
		l.Logger.Warn().Msg(message)
	}
}

func (l *Logger) Error(err error, message string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		l.Logger.Error().Str("f", parserFileName(l.baseDir, file, line)).Err(err).Msg(message)
	} else {
		l.Logger.Error().Msg(message)
	}
}

func (l *Logger) Fatal(err error, message string) {
	_, file, line, ok := runtime.Caller(1)

	if ok {
		l.Logger.Fatal().Str("f", parserFileName(l.baseDir, file, line)).Err(err).Msg(message)
	} else {
		l.Logger.Fatal().Msg(message)
	}
}

func (l *Logger) Panic(err error, message string) {
	_, file, line, ok := runtime.Caller(1)

	if ok {
		l.Logger.Panic().Str("f", parserFileName(l.baseDir, file, line)).Err(err).Msg(message)
	} else {
		l.Logger.Panic().Msg(message)
	}
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func (l *Logger) Ctx(ctx context.Context) *Logger {
	return &Logger{Logger: zerolog.Ctx(ctx)}
}

type RequestLogger struct {
	Logger *zerolog.Logger
}

func NewRequestLogger(config *LogConfig) *RequestLogger {
	var writers []io.Writer

	writers = append(writers, zerolog.ConsoleWriter{
		Out: os.Stdout,
	})

	logWriters := io.MultiWriter(writers...)

	logger := zerolog.New(logWriters).
		With().
		Str("P", "review-system-svc").
		Str("C", "review-system").
		Timestamp().Logger()

	log.Println("Successfully configured request logger")

	return &RequestLogger{
		Logger: &logger,
	}
}
