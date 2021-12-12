package radarr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

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
	SizeOnDisk          int      `json:"sizeOnDisk"`
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
	IsAvailable         bool     `json:"isAvailable"`
	Runtime             int      `json:"runtime"`
	CleanTitle          string   `json:"cleanTitle"`
	ImdbID              string   `json:"imdbId"`
	TmdbID              int      `json:"tmdbid"`
	TitleSlug           string   `json:"titleSlug"`
	Certification       string   `json:"certification"`
	Genres              []string `json:"genres"`
	Tags                []int    `json:"tags"`
	Added               string   `json:"added"`
	Status              string   `json:"status"`

	Path             string `json:"path"`
	RootFolderPath   string `json:"rootFolderPath"`
	QualityProfileId int    `json:"qualityProfileId"`
	FolderName       string `json:"folderName"`
}

func (c *Client) AddMovie(ctx context.Context, m *Movie) error {
	if m == nil {
		return nil
	}
	m.Id = 0

	url := fmt.Sprintf("%s/movie", c.BaseURL)
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if c.Debug {
		fmt.Printf("request: %#v\n", m)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Api-key", c.Apikey)

	response, err := c.Http.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
	default:
		if c.Debug {
			bytes, _ := ioutil.ReadAll(response.Body)
			fmt.Printf("response: %s\n", string(bytes))
		}
		return fmt.Errorf("received %d from radarr", response.StatusCode)
	}

	return nil
}

func (c *Client) Search(ctx context.Context, query string) ([]Movie, error) {

	result := []Movie{}
	q := url.QueryEscape(query)
	url := fmt.Sprintf("%s/movie/lookup?term='%s'", c.BaseURL, q)

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("X-Api-key", c.Apikey)
	response, err := c.Http.Do(r)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
