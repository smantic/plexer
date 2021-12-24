package service

import (
	"context"

	"golift.io/starr"
	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
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
	s := starr.New(c.SonarrKey, c.SonarrURL, 0)
	return Service{
		Radarr: radarr.New(r),
		Sonarr: sonarr.New(s),
	}
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
