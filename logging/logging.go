package logging

import (
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var currentLevel = LogLevelInfo
var logger = log.New(os.Stderr, "", log.LstdFlags)

func SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		currentLevel = LogLevelDebug
	case "info":
		currentLevel = LogLevelInfo
	case "warn", "warning":
		currentLevel = LogLevelWarn
	case "error":
		currentLevel = LogLevelError
	}
}

func LogDebug(format string, args ...interface{}) {
	if currentLevel <= LogLevelDebug {
		logger.Printf("[DEBUG] "+format, args...)
	}
}

func LogInfo(format string, args ...interface{}) {
	if currentLevel <= LogLevelInfo {
		logger.Printf("[INFO] "+format, args...)
	}
}

func LogWarn(format string, args ...interface{}) {
	if currentLevel <= LogLevelWarn {
		logger.Printf("[WARN] "+format, args...)
	}
}

func LogError(format string, args ...interface{}) {
	if currentLevel <= LogLevelError {
		logger.Printf("[ERROR] "+format, args...)
	}
}
