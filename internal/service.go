package internal

import (
	"context"

	"github.com/smantic/plexer/pkg/radarr"
)

type Service struct {
	radarr radarr.Client
}

// Search will serach for a title to add
func (s *Service) Search(ctx context.Context, query string) {
	radmovies, err := s.radarr.Search(ctx, query)
	if err != nil {
		// return err
	}
	_ = radmovies
}
