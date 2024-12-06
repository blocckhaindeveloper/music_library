// song_repository.go
package repositories

import (
	"song_library/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) *SongRepository {
	return &SongRepository{
		db: db,
	}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Song{})
}

func (r *SongRepository) Create(song *models.Song) error {
	return r.db.Create(song).Error
}

func (r *SongRepository) GetByID(id uuid.UUID) (*models.Song, error) {
	var song models.Song
	result := r.db.First(&song, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &song, nil
}

func (r *SongRepository) Update(song *models.Song) error {
	return r.db.Save(song).Error
}

func (r *SongRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Song{}, "id = ?", id).Error
}

func (r *SongRepository) GetAll(filter models.SongFilter, offset, limit int) ([]models.Song, int64, error) {
	var songs []models.Song
	var total int64

	query := r.db.Model(&models.Song{})

	if filter.GroupName != "" {
		query = query.Where("group_name ILIKE ?", "%"+filter.GroupName+"%")
	}
	if filter.SongTitle != "" {
		query = query.Where("song_title ILIKE ?", "%"+filter.SongTitle+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Find(&songs).Error
	if err != nil {
		return nil, 0, err
	}

	return songs, total, nil
}
