package radarr

import (
	"context"
	"encoding/json"
	"net/http"
)

type Client struct {
	BaseURL string
	http    http.Client
}

type Image struct {
	CoverType string `json:"coverType"`
	Url       string `json:"url"`
	RemoteUrl string `json:"remoteUrl"`
}
type Ratings struct {
	Votes int `json:"votes"`
	Value int `json:"value"`
}
type Colletion struct {
	Name   string  `json:"name"`
	TmdbID int     `json:"tmdbId"`
	Images []Image `json:"images"`
}

type Movie struct {
	Id                  int      `json:"id"`
	Title               string   `json:"title"`
	SortTitle           string   `json:"sortTitle"`
	SizeOnDisk          string   `json:"sizeOnDisk"`
	Overview            string   `json:"overview"`
	InCinemas           string   `json:"inCinemas"`
	PhysicalRelease     string   `json:"physicalRelease"`
	Images              []Image  `json:"images"`
	Website             string   `json:"website"`
	Year                int      `json:"year"`
	HasFile             bool     `json:"hasFile"`
	YouTubeTrailerId    string   `json:"youTubeTrailerId"`
	Monitored           bool     `json:"monitored"`
	MinimumAvailability string   `json:"minimumAvailability"`
	IsAvailable         string   `json:"isAvailable"`
	FolderName          string   `json:"folderName"`
	Runtime             int      `json:"runtime"`
	CleanTitle          string   `json:"cleanTitle"`
	ImdbID              int      `json:"imdbId"`
	TitleSlug           string   `json:"titleSlug"`
	Certification       string   `json:"certification"`
	Genres              []string `json:"genres"`
	Tags                []int    `json:"tags"`
	Added               string   `json:"added"`
	Status              string   `json:"status"`
}

func (c *Client) Search(ctx context.Context, query string) ([]Movie, error) {

	result := []Movie{}
	url := c.BaseURL + "/movie/lookup?term=" + query
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}
	response, err := c.http.Do(r)
	if err != nil {
		return result, err
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
