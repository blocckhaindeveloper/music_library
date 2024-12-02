// models/song.go

package models

import (
	"time"
)

type Song struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Group       string    `gorm:"type:varchar(100);not null" json:"group"`
	Song        string    `gorm:"type:varchar(100);not null" json:"song"`
	ReleaseDate time.Time `gorm:"type:date;not null" json:"release_date"`
	Text        string    `gorm:"type:text;not null" json:"text"`
	Link        string    `gorm:"type:varchar(255);not null" json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
