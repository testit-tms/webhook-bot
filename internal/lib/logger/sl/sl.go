package sl

import (
	"golang.org/x/exp/slog"
)

// Err returns a slog.Attr with the given error message as its value and "error" as its key.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
