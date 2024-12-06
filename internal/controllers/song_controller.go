// song_controller.go
package controllers

import (
	"net/http"

	"song_library/internal/models"
	"song_library/internal/services"
	"song_library/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SongController struct {
	SongService *services.SongService
}

func NewSongController(songService *services.SongService) *SongController {
	return &SongController{
		SongService: songService,
	}
}

// GetSongs godoc
// @Summary      Получить список песен
// @Description  Получить список песен с фильтрацией и пагинацией
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        group   query     string  false  "Название группы"
// @Param        song    query     string  false  "Название песни"
// @Param        page    query     int     false  "Номер страницы"
// @Param        limit   query     int     false  "Количество элементов на странице"
// @Success      200     {array}   models.Song
// @Failure      400     {object}  utils.HTTPError
// @Failure      500     {object}  utils.HTTPError
// @Router       /api/songs [get]
func (sc *SongController) GetSongs(c *gin.Context) {
	var filter models.SongFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPError(http.StatusBadRequest, "Некорректные параметры запроса"))
		return
	}

	pagination := utils.NewPaginationFromRequest(c)

	songs, err := sc.SongService.GetSongs(filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, songs)
}

// AddSong godoc
// @Summary      Добавить новую песню
// @Description  Добавить новую песню в библиотеку
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        song  body      models.AddSongRequest  true  "Данные новой песни"
// @Success      201   {object}  models.Song
// @Failure      400   {object}  utils.HTTPError
// @Failure      500   {object}  utils.HTTPError
// @Router       /api/songs [post]
func (sc *SongController) AddSong(c *gin.Context) {
	var req models.AddSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPError(http.StatusBadRequest, "Некорректное тело запроса"))
		return
	}

	song, err := sc.SongService.AddSong(req.GroupName, req.SongTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, song)
}

// GetSongLyrics godoc
// @Summary      Получить текст песни
// @Description  Получить текст песни по ID с пагинацией по куплетам
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id      path      string  true   "ID песни"
// @Param        page    query     int     false  "Номер страницы"
// @Param        limit   query     int     false  "Количество куплетов на странице"
// @Success      200     {object}  models.SongLyricsResponse
// @Failure      400     {object}  utils.HTTPError
// @Failure      404     {object}  utils.HTTPError
// @Failure      500     {object}  utils.HTTPError
// @Router       /api/songs/{id}/lyrics [get]
func (sc *SongController) GetSongLyrics(c *gin.Context) {
	idParam := c.Param("id")
	songID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPError(http.StatusBadRequest, "Некорректный ID песни"))
		return
	}

	pagination := utils.NewPaginationFromRequest(c)

	lyrics, err := sc.SongService.GetSongLyrics(songID, pagination)
	if err != nil {
		if err == services.ErrSongNotFound {
			c.JSON(http.StatusNotFound, utils.NewHTTPError(http.StatusNotFound, "Песня не найдена"))
		} else {
			c.JSON(http.StatusInternalServerError, utils.NewHTTPError(http.StatusInternalServerError, err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, lyrics)
}

// UpdateSong godoc
// @Summary      Обновить данные песни
// @Description  Обновить данные существующей песни по ID
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id    path      string                    true  "ID песни"
// @Param        song  body      models.UpdateSongRequest  true  "Новые данные песни"
// @Success      200   {object}  models.Song
// @Failure      400   {object}  utils.HTTPError
// @Failure      404   {object}  utils.HTTPError
// @Failure      500   {object}  utils.HTTPError
// @Router       /api/songs/{id} [put]
func (sc *SongController) UpdateSong(c *gin.Context) {
	idParam := c.Param("id")
	songID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPError(http.StatusBadRequest, "Некорректный ID песни"))
		return
	}

	var req models.UpdateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPError(http.StatusBadRequest, "Некорректное тело запроса"))
		return
	}

	song, err := sc.SongService.UpdateSong(songID, req)
	if err != nil {
		if err == services.ErrSongNotFound {
			c.JSON(http.StatusNotFound, utils.NewHTTPError(http.StatusNotFound, "Песня не найдена"))
		} else {
			c.JSON(http.StatusInternalServerError, utils.NewHTTPError(http.StatusInternalServerError, err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, song)
}

// DeleteSong godoc
// @Summary      Удалить песню
// @Description  Удалить песню из библиотеки по ID
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "ID песни"
// @Success      204   "No Content"
// @Failure      400   {object}  utils.HTTPError
// @Failure      404   {object}  utils.HTTPError
// @Failure      500   {object}  utils.HTTPError
// @Router       /api/songs/{id} [delete]
func (sc *SongController) DeleteSong(c *gin.Context) {
	idParam := c.Param("id")
	songID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewHTTPError(http.StatusBadRequest, "Некорректный ID песни"))
		return
	}

	err = sc.SongService.DeleteSong(songID)
	if err != nil {
		if err == services.ErrSongNotFound {
			c.JSON(http.StatusNotFound, utils.NewHTTPError(http.StatusNotFound, "Песня не найдена"))
		} else {
			c.JSON(http.StatusInternalServerError, utils.NewHTTPError(http.StatusInternalServerError, err.Error()))
		}
		return
	}

	c.Status(http.StatusNoContent)
}
