// controllers/song_controller.go

package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blocckhaindeveloper/music_library/models"
	"github.com/blocckhaindeveloper/music_library/repository"
	"github.com/blocckhaindeveloper/music_library/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type SongController struct {
	repo   repository.SongRepository
	logger *logrus.Logger
	apiURL string
}

func NewSongController(repo repository.SongRepository, logger *logrus.Logger, apiURL string) *SongController {
	return &SongController{
		repo:   repo,
		logger: logger,
		apiURL: apiURL,
	}
}

// GetSongs godoc
// @Summary Get list of songs
// @Description Get list of songs with optional filtering and pagination
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group query string false "Group name"
// @Param song query string false "Song name"
// @Param skip query int false "Number of records to skip"
// @Param limit query int false "Number of records to limit"
// @Success 200 {object} utils.PaginatedResponse{data=[]models.Song}
// @Failure 400 {object} gin.H{"error": "Invalid skip parameter"}
// @Failure 400 {object} gin.H{"error": "Invalid limit parameter"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /songs [get]
func (sc *SongController) GetSongs(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "10")

	skip, err := strconv.Atoi(skipStr)
	if err != nil {
		sc.logger.Error("Invalid skip parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skip parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		sc.logger.Error("Invalid limit parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	filters := make(map[string]interface{})
	if group != "" {
		filters["group"] = group
	}
	if song != "" {
		filters["song"] = song
	}

	total, songs, err := sc.repo.GetSongs(filters, skip, limit)
	if err != nil {
		sc.logger.Error("Error fetching songs: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	sc.logger.Info("Fetched songs list")
	c.JSON(http.StatusOK, utils.PaginatedResponse{
		Total: total,
		Data:  songs,
	})
}

// GetSongByID godoc
// @Summary Get song by ID
// @Description Get song details by ID
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Success 200 {object} models.Song
// @Failure 400 {object} gin.H{"error": "Invalid song ID"}
// @Failure 404 {object} gin.H{"error": "Song not found"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /songs/{id} [get]
func (sc *SongController) GetSongByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sc.logger.Error("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	song, err := sc.repo.GetSongByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			sc.logger.Warn("Song not found with ID: ", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		sc.logger.Error("Error fetching song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	sc.logger.Info("Fetched song with ID: ", id)
	c.JSON(http.StatusOK, song)
}

// CreateSong godoc
// @Summary Create a new song
// @Description Add a new song to the library
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body models.Song true "Song data"
// @Success 201 {object} models.Song
// @Failure 400 {object} gin.H{"error": "Invalid input"}
// @Failure 400 {object} gin.H{"error": "Missing required fields"}
// @Failure 502 {object} gin.H{"error": "External API request failed"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /songs [post]
func (sc *SongController) CreateSong(c *gin.Context) {
	var input models.Song
	if err := c.ShouldBindJSON(&input); err != nil {
		sc.logger.Error("Invalid input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Валидация полей
	if input.Group == "" || input.Song == "" || input.ReleaseDate.IsZero() || input.Text == "" || input.Link == "" {
		sc.logger.Warn("Missing required fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Создание песни
	if err := sc.repo.CreateSong(&input); err != nil {
		sc.logger.Error("Error creating song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	sc.logger.Info("Created new song with ID: ", input.ID)

	// Запрос к внешнему API для обогащения данных
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"group": input.Group,
			"song":  input.Song,
		}).
		SetHeader("Content-Type", "application/json").
		Get(sc.apiURL)

	if err != nil {
		sc.logger.Error("External API request failed: ", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "External API request failed"})
		return
	}

	if resp.StatusCode() != http.StatusOK {
		sc.logger.Warn("External API returned non-200 status: ", resp.Status())
		c.JSON(http.StatusBadGateway, gin.H{"error": "External API request failed"})
		return
	}

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &apiResponse); err != nil {
		sc.logger.Error("Error parsing external API response: ", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "External API response invalid"})
		return
	}

	// Обогащение данных песни
	if releaseDateStr, ok := apiResponse["releaseDate"].(string); ok {
		releaseDate, err := time.Parse("02.01.2006", releaseDateStr)
		if err == nil {
			input.ReleaseDate = releaseDate
		}
	}
	if text, ok := apiResponse["text"].(string); ok {
		input.Text = text
	}
	if link, ok := apiResponse["link"].(string); ok {
		input.Link = link
	}

	// Обновление записи в базе данных
	if err := sc.repo.UpdateSong(&input); err != nil {
		sc.logger.Error("Error updating song with enriched data: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	sc.logger.Info("Updated song with enriched data, ID: ", input.ID)
	c.JSON(http.StatusCreated, input)
}

// UpdateSong godoc
// @Summary Update an existing song
// @Description Update song details by ID
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song data"
// @Success 200 {object} models.Song
// @Failure 400 {object} gin.H{"error": "Invalid song ID"}
// @Failure 400 {object} gin.H{"error": "Invalid input"}
// @Failure 404 {object} gin.H{"error": "Song not found"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /songs/{id} [put]
func (sc *SongController) UpdateSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sc.logger.Error("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var input models.Song
	if err := c.ShouldBindJSON(&input); err != nil {
		sc.logger.Error("Invalid input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	song, err := sc.repo.GetSongByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			sc.logger.Warn("Song not found with ID: ", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		sc.logger.Error("Error fetching song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Обновление полей
	if input.Group != "" {
		song.Group = input.Group
	}
	if input.Song != "" {
		song.Song = input.Song
	}
	if !input.ReleaseDate.IsZero() {
		song.ReleaseDate = input.ReleaseDate
	}
	if input.Text != "" {
		song.Text = input.Text
	}
	if input.Link != "" {
		song.Link = input.Link
	}

	if err := sc.repo.UpdateSong(song); err != nil {
		sc.logger.Error("Error updating song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	sc.logger.Info("Updated song with ID: ", id)
	c.JSON(http.StatusOK, song)
}

// DeleteSong godoc
// @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Success 204 {object} nil
// @Failure 400 {object} gin.H{"error": "Invalid song ID"}
// @Failure 404 {object} gin.H{"error": "Song not found"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /songs/{id} [delete]
func (sc *SongController) DeleteSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sc.logger.Error("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	song, err := sc.repo.GetSongByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			sc.logger.Warn("Song not found with ID: ", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		sc.logger.Error("Error fetching song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := sc.repo.DeleteSong(uint(id)); err != nil {
		sc.logger.Error("Error deleting song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	sc.logger.Info("Deleted song with ID: ", id)
	c.Status(http.StatusNoContent)
}

// GetLyrics godoc
// @Summary Get song lyrics with pagination
// @Description Get lyrics of a song by ID with pagination by verses
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param page query int false "Page number"
// @Param per_page query int false "Verses per page"
// @Success 200 {object} models.Song
// @Failure 400 {object} gin.H{"error": "Invalid song ID"}
// @Failure 400 {object} gin.H{"error": "Page out of range"}
// @Failure 404 {object} gin.H{"error": "Song not found"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /songs/{id}/lyrics [get]
func (sc *SongController) GetLyrics(c *gin.Context) {
	idStr := c.Param("id")
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "4")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sc.logger.Error("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		perPage = 4
	}

	song, err := sc.repo.GetSongByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			sc.logger.Warn("Song not found with ID: ", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		sc.logger.Error("Error fetching song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	verses := strings.Split(song.Text, "\n\n")
	start := (page - 1) * perPage
	end := start + perPage

	if start >= len(verses) {
		sc.logger.Warn("Page out of range for song ID: ", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page out of range"})
		return
	}

	if end > len(verses) {
		end = len(verses)
	}

	paginatedText := strings.Join(verses[start:end], "\n\n")

	// Создание ответа
	response := models.Song{
		ID:          song.ID,
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: song.ReleaseDate,
		Text:        paginatedText,
		Link:        song.Link,
	}

	sc.logger.Info("Fetched lyrics page ", page, " for song ID: ", id)
	c.JSON(http.StatusOK, response)
}
