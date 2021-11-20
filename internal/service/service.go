package service

import (
	"context"

	"github.com/smantic/plexer/pkg/radarr"
	"github.com/webtor-io/go-jackett"
)

type Service struct {
	Radarr  *radarr.Client
	Jackett *jackett.Jackett
}

// Search will serach for a title to add
func (s *Service) Search(ctx context.Context, query string) ([]string, error) {
	movies, err := s.Radarr.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	titles := make([]string, 0, len(movies))
	for _, m := range movies {
		titles = append(titles, m.Title)
	}
	return titles, nil
}
