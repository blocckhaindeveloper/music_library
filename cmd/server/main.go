// cmd/server/main.go

package main

import (
	"fmt"

	"github.com/blocckhaindeveloper/music_library/config"
	"github.com/blocckhaindeveloper/music_library/controllers"
	"github.com/blocckhaindeveloper/music_library/repository"
	"github.com/blocckhaindeveloper/music_library/routes"
	"github.com/blocckhaindeveloper/music_library/utils"
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Music Library API
// @version 1.0
// @description API for managing an online music library.

// @contact.name API Support
// @contact.email support@musiclibrary.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация логирования
	utils.InitLogger(cfg.LogLevel)
	logger := utils.Logger

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}

	// Получение *sql.DB из *gorm.DB для использования с Goose
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get sql.DB from gorm.DB: ", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error("Error closing database connection: ", err)
		}
	}()

	// Применение миграций
	goose.SetLogger(logger)
	goose.SetBaseFS(nil)

	// Установка диалекта для Goose
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatal("Failed to set goose dialect: ", err)
	}

	// Применение всех миграций из директории "migrations"
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		logger.Fatal("Failed to apply migrations: ", err)
	}
	logger.Info("Database migrations applied successfully")

	// Инициализация репозиториев и контроллеров
	songRepo := repository.NewSongRepository(db)
	songController := controllers.NewSongController(songRepo, logger, cfg.APIURL)

	// Настройка маршрутов
	router := gin.Default()
	routes.SetupRoutes(router, songController)

	// Настройка Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("Starting server on ", addr)
	if err := router.Run(addr); err != nil {
		logger.Fatal("Failed to run server: ", err)
	}
}
