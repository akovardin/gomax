package core

import "github.com/akovardin/gomax/logging"

type LogLevel = logging.LogLevel

const (
	LogLevelDebug = logging.LogLevelDebug
	LogLevelInfo  = logging.LogLevelInfo
	LogLevelWarn  = logging.LogLevelWarn
	LogLevelError = logging.LogLevelError
)

var SetLogLevel = logging.SetLogLevel
var LogDebug = logging.LogDebug
var LogInfo = logging.LogInfo
var LogWarn = logging.LogWarn
var LogError = logging.LogError
