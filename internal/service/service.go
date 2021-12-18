package service

import (
	"context"

	"github.com/smantic/starr"
	"github.com/smantic/starr/radarr"
	"github.com/smantic/starr/sonarr"
)

type ContentType string

const (
	CONTENT_MOVIE ContentType = "movie"
	CONTENT_SHOW  ContentType = "show"
)

type Service struct {
	Radarr *radarr.Radarr
	Sonarr *sonarr.Sonarr
}

type Config struct {
	RadarrURL   string
	RadarrKey   string
	RadarrDebug bool

	SonarrKey string
	SonarrURL string
}

func New(c *Config) Service {
	r := starr.New(c.RadarrKey, c.RadarrURL, 0)
	s := starr.New(c.SonarrURL, c.SonarrURL, 0)
	return Service{
		Radarr: radarr.New(r),
		Sonarr: sonarr.New(s),
	}
}

// AddMovie will add a movie to radarr
func (s *Service) AddMovie(ctx context.Context, m radarr.AddMovieInput) error {

	_, err := s.Radarr.AddMovie(&m)
	if err != nil {
		return err
	}

	return nil
}

// AddShow will add a show to sonarr
func (s *Service) AddShow(ctx context.Context, in sonarr.AddSeriesInput) error {

	_, err := s.Sonarr.AddSeries(&in)
	if err != nil {
		return err
	}
	return nil
}

// Send a command that start a search for missing.
func (s *Service) SearchForMissing(ctx context.Context, kind ContentType) error {

	var err error
	switch kind {
	case CONTENT_MOVIE:
		_, err = s.Sonarr.SendCommand(&sonarr.CommandRequest{Name: "missingEpisodeSearch"})
	case CONTENT_SHOW:
		_, err = s.Radarr.SendCommand(&radarr.CommandRequest{Name: "MissingMoviesSearch"})
	default:
		_, err = s.Radarr.SendCommand(&radarr.CommandRequest{Name: "MissingMoviesSearch"})
		if err != nil {
			break
		}
		_, err = s.Sonarr.SendCommand(&sonarr.CommandRequest{Name: "missingEpisodeSearch"})
	}

	return err
}
