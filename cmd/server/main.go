// main.go
package main

import (
	"song_library/configs"
	_ "song_library/docs"
	"song_library/internal/app"
	"song_library/internal/utils"
)

// @title           Song Library API
// @version         1.0
// @description     API для управления онлайн-библиотекой песен

// @host      localhost:8080
// @BasePath  /
func main() {

	utils.LoadEnv()

	cfg, err := configs.LoadConfig()
	if err != nil {
		utils.GetLogger().Fatalf("Не удалось загрузить конфигурацию: %v", err)

	}

	application := app.NewApp(cfg)

	if err := application.Run(); err != nil {
		utils.GetLogger().Fatalf("Не удалось запустить приложение: %v", err)
	}
}
