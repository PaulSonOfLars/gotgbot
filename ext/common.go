package ext

import (
	"log"
	"log/slog"
)

func logError(l *slog.Logger, text string, err error) {
	if l == nil {
		log.Printf("ERROR: %s: %s", text, err)
		return
	}
	l.Error(text, "error", err)
}

// ternary operator approximation.
func iftrue[T any](b bool, t T, f T) T {
	if b {
		return t
	}
	return f
}
