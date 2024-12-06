// song.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey" example:"123e4567-e89b-12d3-a456-426614174000"`
	GroupName   string    `json:"group" gorm:"not null;column:group_name" example:"Muse"`
	SongTitle   string    `json:"song" gorm:"not null;column:song_title" example:"Supermassive Black Hole"`
	ReleaseDate time.Time `json:"releaseDate" example:"2006-07-16T00:00:00Z"`
	Text        string    `json:"text" example:"Ooh baby, don't you know I suffer?..."`
	Link        string    `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type AddSongRequest struct {
	GroupName string `json:"group" binding:"required" example:"Muse"`
	SongTitle string `json:"song" binding:"required" example:"Supermassive Black Hole"`
}

type UpdateSongRequest struct {
	GroupName string `json:"group" example:"Muse"`
	SongTitle string `json:"song" example:"Uprising"`
}

type SongFilter struct {
	GroupName string `form:"group"`
	SongTitle string `form:"song"`
}

type SongLyricsResponse struct {
	Verses []string `json:"verses"`
	Page   int      `json:"page"`
	Limit  int      `json:"limit"`
	Total  int      `json:"total"`
}
