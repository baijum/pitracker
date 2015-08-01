package logger

import (
	"os"
	"testing"

	"github.com/baijum/pitracker/logger"
)

func TestDefalutLevel(t *testing.T) {
	if logger.Level != logger.INFO {
		t.Log("Wrong default log level", logger.Level)
	}
}

func TestLevelValue(t *testing.T) {
	if logger.DEBUG != 1 {
		t.Log("DEBUG value changed", logger.DEBUG)
	}
	if logger.INFO != 2 {
		t.Log("INFO value changed", logger.INFO)
	}
	if logger.WARNING != 4 {
		t.Log("WARNING value changed", logger.WARNING)
	}
	if logger.ERROR != 8 {
		t.Log("ERROR value changed", logger.ERROR)
	}
}

func TestSetLevel(t *testing.T) {
	logger.SetLevel(logger.DEBUG)
	if logger.Level != logger.DEBUG {
		t.Log("Wrong log level", logger.Level)
	}
}

func TestInitLogLevel(t *testing.T) {
	os.Setenv("PITRACKER_LOG_LEVEL", "DEBUG")
	logger.InitLogLevel()
	if logger.Level != logger.DEBUG {
		t.Log("Log level not initialized", logger.Level)
	}
	os.Setenv("PITRACKER_LOG_LEVEL", "INFO")
	logger.InitLogLevel()
	if logger.Level != logger.INFO {
		t.Log("Log level not initialized", logger.Level)
	}
	os.Setenv("PITRACKER_LOG_LEVEL", "WARNING")
	logger.InitLogLevel()
	if logger.Level != logger.WARNING {
		t.Log("Log level not initialized", logger.Level)
	}
	os.Setenv("PITRACKER_LOG_LEVEL", "ERROR")
	logger.InitLogLevel()
	if logger.Level != logger.ERROR {
		t.Log("Log level not initialized", logger.Level)
	}
}
