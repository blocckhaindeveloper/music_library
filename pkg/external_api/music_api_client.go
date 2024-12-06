// music_api_client.go
package external_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type MusicAPIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewMusicAPIClient(baseURL string) *MusicAPIClient {
	return &MusicAPIClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

type SongDetailResponse struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (client *MusicAPIClient) GetSongInfo(group, song string) (*SongDetailResponse, error) {
	endpoint, err := url.Parse(client.BaseURL)
	if err != nil {
		return nil, err
	}

	queryParams := endpoint.Query()
	queryParams.Set("group", group)
	queryParams.Set("song", song)
	endpoint.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("внешний API вернул статус %d", resp.StatusCode)
	}

	var songDetail SongDetailResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&songDetail); err != nil {
		return nil, err
	}

	return &songDetail, nil
}
