// utils/logger.go

package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger(logLevel string) {
	Logger = logrus.New()

	// Установка уровня логирования
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// Настройка вывода
	Logger.SetOutput(os.Stdout)

	// Формат логирования
	Logger.SetFormatter(&logrus.JSONFormatter{})
}
