package service

import (
	"context"
	"time"

	"github.com/smantic/starr/radarr"
	"github.com/smantic/starr/sonarr"
)

type ContentInfo struct {
	Title     string
	Year      int
	Genre     []string
	Overview  string
	ImdbID    string
	Added     time.Time
	Monitored bool
}

func contentInfoFromSerries(s sonarr.SeriesLookup) ContentInfo {
	return ContentInfo{
		Title:     s.Title,
		Year:      s.Year,
		Genre:     s.Genres,
		Overview:  s.Overview,
		ImdbID:    s.ImdbID,
		Added:     s.Added,
		Monitored: s.Monitored,
	}
}

func contentInfoFromMovie(m radarr.Movie) ContentInfo {
	return ContentInfo{
		Title:     m.Title,
		Year:      m.Year,
		Genre:     m.Genres,
		Overview:  m.Overview,
		ImdbID:    m.ImdbID,
		Added:     m.Added,
		Monitored: m.Monitored,
	}
}

// Search will serach for content and return some information.
func (s *Service) Search(ctx context.Context, kind ContentType, query string, limit int) ([]ContentInfo, error) {

	var (
		content = make([]ContentInfo, 0, limit)
		err     error
	)

	switch kind {
	case CONTENT_SHOW:
		var serries []*sonarr.SeriesLookup

		serries, err = s.Sonarr.GetSeriesLookup(query, 0)
		for _, s := range serries {
			if s == nil {
				continue
			}
			content = append(content, contentInfoFromSerries(*s))
			if len(content) > limit {
				break
			}
		}
	case CONTENT_MOVIE:
		var movies []radarr.Movie

		movies, err = s.Radarr.Lookup(query)
		for i, m := range movies {
			if i >= limit {
				break
			}
			content = append(content, contentInfoFromMovie(m))
		}
	}

	return content, err
}
