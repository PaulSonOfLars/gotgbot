package ext

import "log"

// Logger is an interface to use non standard loggers.
type Logger interface {
	Printf(format string, v ...interface{})
}

var _ Logger = log.Default()
