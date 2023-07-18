// Logger/logger.go
package Logger

import (
	"github.com/opencoff/go-logger"
	"os"
	"strings"
)

type App struct {
	Log *logger.Logger
}

func (a *App) InitLogger(prefix string) error { // note the "string" type for prefix
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

	a.Log, err = logger.New(os.Stdout, level, prefix, logger.LstdFlags)
	return err
}
