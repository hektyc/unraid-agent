package logger

import (
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

func Init(level, format string) {
	once.Do(func() {
		logger = logrus.New()
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			lvl = logrus.InfoLevel
		}
		logger.SetLevel(lvl)

		if format == "json" {
			logger.SetFormatter(&logrus.JSONFormatter{})
		} else {
			logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		}

		logger.SetOutput(os.Stdout)
	})
}

func Get() *logrus.Logger {
	if logger == nil {
		Init("info", "text")
	}
	return logger
}

func Redact(s string) string {
	s = strings.ReplaceAll(s, os.Getenv("UNRAID_API_KEY"), "[REDACTED]")
	if bearer := os.Getenv("UNRAID_MCP_BEARER_TOKEN"); bearer != "" {
		s = strings.ReplaceAll(s, bearer, "[REDACTED]")
	}
	return s
}
