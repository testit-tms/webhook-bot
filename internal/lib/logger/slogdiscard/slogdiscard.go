package slogdiscard

import (
	"context"

	"golang.org/x/exp/slog"
)

// NewDiscardLogger creates a new logger that discards all log records.
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler is a logger handler that discards all log records.
type DiscardHandler struct{}

// NewDiscardHandler creates a new DiscardHandler instance.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Handle ignores the log record.
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs returns the handler itself since DiscardHandler does not use attributes.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns the handler itself since DiscardHandler does not use groups.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

// Enabled returns false since DiscardHandler does not enable any log levels.
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
