// song_service.go
package services

import (
	"errors"
	"strings"
	"time"

	"song_library/internal/models"
	"song_library/internal/repositories"
	"song_library/internal/utils"
	"song_library/pkg/external_api"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrSongNotFound = errors.New("песня не найдена")
)

type SongService struct {
	SongRepo          *repositories.SongRepository
	ExternalAPIClient *external_api.MusicAPIClient
}

func NewSongService(repo *repositories.SongRepository, apiClient *external_api.MusicAPIClient) *SongService {
	return &SongService{
		SongRepo:          repo,
		ExternalAPIClient: apiClient,
	}
}

func (s *SongService) GetSongs(filter models.SongFilter, pagination *utils.Pagination) ([]models.Song, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()

	songs, _, err := s.SongRepo.GetAll(filter, offset, limit)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *SongService) AddSong(groupName, songTitle string) (*models.Song, error) {

	songDetail, err := s.ExternalAPIClient.GetSongInfo(groupName, songTitle)
	if err != nil {
		return nil, err
	}

	releaseDate, err := time.Parse("02.01.2006", songDetail.ReleaseDate)
	if err != nil {
		releaseDate = time.Time{}
	}

	newSong := &models.Song{
		ID:          uuid.New(),
		GroupName:   groupName,
		SongTitle:   songTitle,
		ReleaseDate: releaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	err = s.SongRepo.Create(newSong)
	if err != nil {
		return nil, err
	}

	return newSong, nil
}

func (s *SongService) GetSongLyrics(id uuid.UUID, pagination *utils.Pagination) (*models.SongLyricsResponse, error) {
	song, err := s.SongRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSongNotFound
		}
		return nil, err
	}

	verses := strings.Split(song.Text, "\n\n")
	totalVerses := len(verses)

	offset := pagination.GetOffset()
	limit := pagination.GetLimit()

	if offset > totalVerses {
		offset = totalVerses
	}

	end := offset + limit
	if end > totalVerses {
		end = totalVerses
	}

	paginatedVerses := verses[offset:end]

	response := &models.SongLyricsResponse{
		Verses: paginatedVerses,
		Page:   pagination.Page,
		Limit:  pagination.Limit,
		Total:  totalVerses,
	}

	return response, nil
}

func (s *SongService) UpdateSong(id uuid.UUID, req models.UpdateSongRequest) (*models.Song, error) {
	song, err := s.SongRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSongNotFound
		}
		return nil, err
	}

	if req.GroupName != "" {
		song.GroupName = req.GroupName
	}
	if req.SongTitle != "" {
		song.SongTitle = req.SongTitle
	}

	err = s.SongRepo.Update(song)
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (s *SongService) DeleteSong(id uuid.UUID) error {
	err := s.SongRepo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSongNotFound
		}
		return err
	}
	return nil
}
