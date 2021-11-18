package service

import (
	"context"
	"fmt"

	"github.com/smantic/plexer/pkg/radarr"
)

type Service struct {
	radarr *radarr.Client
}

type Dependencies struct {
	Radarr *radarr.Client
	// torrent client
	// torrent finder
}

func NewService(deps *Dependencies) (Service, error) {

	if deps.Radarr == nil {
		return Service{}, fmt.Errorf("expected non nil radarr client")
	}

	return Service{
		radarr: deps.Radarr,
	}, nil
}

// Search will serach for a title to add
func (s *Service) Search(ctx context.Context, query string) {
	radmovies, err := s.radarr.Search(ctx, query)
	if err != nil {
		// return err
	}
	_ = radmovies
}
