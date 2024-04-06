package shim

import (
	"fmt"
	"log/slog"
)

// Log is a simple logging interface that is satisfied by the standard library logger amongst other idiomatic loggers
type Log interface {
	Printf(format string, v ...interface{})
}

type slogAdapter struct {
	Logger slog.Logger
}

func (sa slogAdapter) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	sa.Logger.Debug(msg)
}
