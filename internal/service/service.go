package service

import (
	"context"
	"log"

	"github.com/smantic/plexer/pkg/radarr"
)

const (
	CONTENT_MOVIE string = "movie"
	CONTENT_SHOW  string = "show"
)

type Service struct {
	Radarr *radarr.Client
}

func (s *Service) Add(ctx context.Context, m radarr.Movie) error {

	m.Monitored = true // automatic monitoring
	if m.QualityProfileId == 0 {
		m.QualityProfileId = 6
	}

	err := s.Radarr.AddMovie(ctx, &m)
	if err != nil {
		return err
	}

	go func() {
		// tell radarr to search for the newly added title
		r := radarr.CommandPayload{
			Name: radarr.CommandMissingMoviesSearch,
		}
		err := s.Radarr.PostCommand(ctx, r)

		if err != nil {
			log.Printf("failed to tell radar to search for missing movies: %v\n", err)
			return
		}
	}()

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
