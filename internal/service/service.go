package service

import (
	"context"

	"github.com/smantic/plexer/pkg/radarr"
)

type Service struct {
	Radarr *radarr.Client
}

// Search will serach for a title to add.
func (s *Service) Search(ctx context.Context, query string) ([]radarr.Movie, error) {
	movies, err := s.Radarr.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	return movies, nil
}
