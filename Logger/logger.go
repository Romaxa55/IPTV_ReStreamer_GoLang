package Logger

import (
	"github.com/opencoff/go-logger"
	"os"
	"strings"
)

var Log *logger.Logger

func InitLogger() error {
	var err error
	logLevel := os.Getenv("LOG_LEVEL")

	var level logger.Priority
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = logger.LOG_DEBUG
	case "INFO":
		level = logger.LOG_INFO
	case "WARN":
		level = logger.LOG_WARNING
	case "ERROR":
		level = logger.LOG_ERR
	default:
		level = logger.LOG_INFO
	}

	Log, err = logger.New(os.Stdout, level, "Webserver", logger.LstdFlags)
	return err
}
