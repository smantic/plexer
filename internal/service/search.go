package service

import (
	"context"
	"time"

	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
)

type ContentInfo struct {
	Title       string
	Year        int
	Genre       []string
	Overview    string
	ImdbID      string
	Added       time.Time
	Monitored   bool
	Size        int64
	ContentType ContentType

	// raw is the struct we originally got from *starr
	raw interface{}
}

type searchResult struct {
	data []ContentInfo
	err  error
}

// Search will serach for content and return some information.
func (s *Service) Search(ctx context.Context, kind ContentType, query string, limit int) ([]ContentInfo, error) {

	if query == "" {
		return nil, nil
	}

	if limit == 0 { // default page
		limit = 10
	}

	var (
		content = make([]ContentInfo, 0, limit)
		err     error
	)

	switch kind {
	case CONTENT_SHOW:
		content, err = s.serachSonar(ctx, query)
	case CONTENT_MOVIE:
		content, err = s.serachRadarr(ctx, query)
	default:

		resultChan := make(chan searchResult, 2)
		go func() {
			movies, err := s.serachRadarr(ctx, query)
			resultChan <- searchResult{movies, err}
		}()

		go func() {
			serries, err := s.serachSonar(ctx, query)
			resultChan <- searchResult{serries, err}
		}()

		x := <-resultChan
		y := <-resultChan
		close(resultChan)

		if x.err != nil {
			return nil, x.err
		}
		if y.err != nil {
			return nil, y.err
		}

		// TODO log number of results for each type.
		var i, j int
		for { // zip merge
			if (i >= len(x.data) && j >= len(y.data)) || i+j > limit {
				break
			}
			if i < len(x.data) {
				content = append(content, x.data[i])
				i++
			}
			if j < len(y.data) {
				content = append(content, y.data[j])
				j++
			}
		}
	}

	if limit > len(content) {
		limit = len(content)
	}

	return content[:limit], err
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
		Size:      s.Statistics.SizeOnDisk,

		raw: s,
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
		Size:      m.SizeOnDisk,

		raw: m,
	}
}
func (s *Service) serachSonar(ctx context.Context, query string) ([]ContentInfo, error) {

	serries, err := s.Sonarr.GetSeriesLookup(query, 0)
	if err != nil {
		return nil, err
	}

	content := make([]ContentInfo, 0, len(serries))
	for _, s := range serries {
		if s == nil {
			continue
		}
		content = append(content, contentInfoFromSerries(*s))
	}
	return content, nil
}

func (s *Service) serachRadarr(ctx context.Context, query string) ([]ContentInfo, error) {

	movies, err := s.Radarr.Lookup(query)
	if err != nil {
		return nil, err
	}

	content := make([]ContentInfo, 0, len(movies))
	for _, m := range movies {
		content = append(content, contentInfoFromMovie(m))
	}
	return content, nil
}
