package logger

import (
	"errors"
	"os"
	"strings"

	"golang.org/x/exp/slog"
)

func New(level string) (*slog.Logger, error) {
	logLevel, err := getLogLevel(level)
	if err != nil {
		return nil, err
	}

	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}),
	), nil
}

func getLogLevel(l string) (slog.Level, error) {

	switch strings.ToLower(l) {
	case "err", "error":
		return slog.LevelError, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	}
	return 0, errors.New("invalid log level name: " + l)
}
