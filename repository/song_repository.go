// repository/song_repository.go

package repository

import (
	"errors"

	"github.com/blocckhaindeveloper/music_library/models"
	"gorm.io/gorm"
)

type SongRepository interface {
	GetSongs(filters map[string]interface{}, offset, limit int) (int64, []models.Song, error)
	GetSongByID(id uint) (*models.Song, error)
	CreateSong(song *models.Song) error
	UpdateSong(song *models.Song) error
	DeleteSong(id uint) error
}

type songRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) SongRepository {
	return &songRepository{db}
}

func (r *songRepository) GetSongs(filters map[string]interface{}, offset, limit int) (int64, []models.Song, error) {
	var songs []models.Song
	var total int64

	query := r.db.Model(&models.Song{})

	for key, value := range filters {
		query = query.Where(key+" ILIKE ?", "%"+value.(string)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&songs).Error; err != nil {
		return 0, nil, err
	}

	return total, songs, nil
}

func (r *songRepository) GetSongByID(id uint) (*models.Song, error) {
	var song models.Song
	if err := r.db.First(&song, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &song, nil
}

func (r *songRepository) CreateSong(song *models.Song) error {
	return r.db.Create(song).Error
}

func (r *songRepository) UpdateSong(song *models.Song) error {
	return r.db.Save(song).Error
}

func (r *songRepository) DeleteSong(id uint) error {
	return r.db.Delete(&models.Song{}, id).Error
}
