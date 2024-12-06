// logger.go
package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func InitLogger(level string) *logrus.Logger {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger
}

func GetLogger() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()
	}
	return logger
}
