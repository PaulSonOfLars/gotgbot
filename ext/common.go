package ext

import (
	"log"
	"log/slog"
)

func logError(l *slog.Logger, text string, err error) {
	if l == nil {
		log.Printf("ERROR: %s: %s", text, err.Error())
		return
	}
	l.Error(text, "error", err.Error())
}

func logDebug(l *slog.Logger, text string, args ...any) {
	if l == nil {
		// No logger? No debug.
		return
	}
	l.Error(text, args...)
}

// ternary operator approximation.
func iftrue[T any](b bool, t T, f T) T {
	if b {
		return t
	}
	return f
}
