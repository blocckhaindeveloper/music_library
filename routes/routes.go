// routes/routes.go

package routes

import (
	"github.com/blocckhaindeveloper/music_library/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, songController *controllers.SongController) {
	api := router.Group("/api")
	{
		songs := api.Group("/songs")
		{
			songs.GET("", songController.GetSongs)
			songs.GET("/:id", songController.GetSongByID)
			songs.POST("", songController.CreateSong)
			songs.PUT("/:id", songController.UpdateSong)
			songs.DELETE("/:id", songController.DeleteSong)
			songs.GET("/:id/lyrics", songController.GetLyrics)
		}
	}
}
