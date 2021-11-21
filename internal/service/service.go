package service

import (
	"context"

	"github.com/smantic/plexer/pkg/radarr"
)

type Service struct {
	Radarr *radarr.Client
}

func (s *Service) Add(ctx context.Context, m radarr.Movie) error {

	if m.QualityProfileId == 0 {
		m.QualityProfileId = 1
	}

	err := s.Radarr.AddMovie(ctx, &m)
	if err != nil {
		return err
	}

	return nil
}

// Search will serach for a title to add.
func (s *Service) Search(ctx context.Context, query string) ([]radarr.Movie, error) {
	if query == "" {
		return nil, nil
	}
	movies, err := s.Radarr.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	return movies, nil
}
