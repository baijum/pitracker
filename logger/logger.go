/* Package logger implements levels for logging

Usage:

	if logger.Level <= logger.INFO {
		log.Printf("Some log message")
	}

	if logger.Level <= logger.WARNING {
		log.Printf("Some log message")
	}


*/
package logger

import (
	"os"
	"strings"
)

const (
	DEBUG = 1 << iota
	INFO
	WARNING
	ERROR
)

// Log level
var Level int8 = INFO

// Set log level
func SetLevel(l int8) {
	Level = l
}

func InitLogLevel() {
	logLevel := os.Getenv("PITRACKER_LOG_LEVEL")
	logLevel = strings.ToUpper(logLevel)
	switch {
	case logLevel == "DEBUG":
		Level = DEBUG
	case logLevel == "INFO":
		Level = INFO
	case logLevel == "WARNING":
		Level = WARNING
	case logLevel == "ERROR":
		Level = ERROR
	}
}

func init() {
	InitLogLevel()
}
