package logger

import (
	"go.uber.org/zap"
)

// Log is the globally accessible sugared logger
var Log *zap.SugaredLogger

// Init initializes the logger. You call this once in main.go.
func Init() {
	// 1. Create a base development logger
	baseLogger, _ := zap.NewDevelopment()

	// 2. Convert it to a SugaredLogger and assign to the global variable
	Log = baseLogger.Sugar()
}

// Sync flushes any buffered log entries
func Sync() {
	_ = Log.Sync()
}
