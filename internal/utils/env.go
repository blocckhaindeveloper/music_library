// env.go
package utils

import (
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		GetLogger().Info("Не удалось загрузить .env файл, используем переменные окружения")
	}
}
