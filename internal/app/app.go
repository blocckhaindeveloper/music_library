// app.go
package app

import (
	"fmt"

	"song_library/configs"
	"song_library/internal/controllers"
	"song_library/internal/middleware"
	"song_library/internal/repositories"
	"song_library/internal/services"
	"song_library/internal/utils"
	"song_library/pkg/external_api"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Config *configs.Config
	Router *gin.Engine
	DB     *gorm.DB
}

func NewApp(cfg *configs.Config) *App {

	logger := utils.InitLogger(cfg.LogLevel)
	gin.DefaultWriter = logger.Writer()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	if err := repositories.AutoMigrate(db); err != nil {
		logger.Fatalf("Не удалось выполнить миграцию базы данных: %v", err)
	}

	externalAPIClient := external_api.NewMusicAPIClient(cfg.ExternalAPI)

	songRepo := repositories.NewSongRepository(db)
	songService := services.NewSongService(songRepo, externalAPIClient)
	songController := controllers.NewSongController(songService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())

	RegisterRoutes(router, songController)

	return &App{
		Config: cfg,
		Router: router,
		DB:     db,
	}
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%s", a.Config.ServerPort)
	utils.GetLogger().Infof("Запуск сервера на %s", addr)
	return a.Router.Run(addr)
}

func RegisterRoutes(router *gin.Engine, songController *controllers.SongController) {
	api := router.Group("/api")
	{
		songs := api.Group("/songs")
		{
			songs.GET("", songController.GetSongs)
			songs.POST("", songController.AddSong)
			songs.GET("/:id/lyrics", songController.GetSongLyrics)
			songs.PUT("/:id", songController.UpdateSong)
			songs.DELETE("/:id", songController.DeleteSong)
		}
	}

	// Маршрут для Swagger документации
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
