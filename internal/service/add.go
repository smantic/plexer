package service

import (
	"context"

	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
)

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
