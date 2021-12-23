package service

import (
	"context"
	"fmt"

	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
)

func (s *Service) Add(ctx context.Context, content ContentInfo) error {

	infos, err := s.GetRootFolderInfos(ctx)
	if err != nil {
		return err
	}

	var found FolderInfo
	for _, info := range infos {
		if info.ContentType == content.ContentType && info.FreeSpace > content.Size {
			found = info
			break
		}
	}

	if found.Path == "" {
		return fmt.Errorf(
			"could not find a folder with enough capcacity (%d MB)", content.Size/1000000,
		)
	}

	switch content.ContentType {
	case CONTENT_MOVIE:
		return s.AddMovie(ctx, content, found.Path)
	case CONTENT_SHOW:
		return s.AddShow(ctx, content, found.Path)
	default:
		return fmt.Errorf("unsupported content type")
	}
}

// AddMovie will add a movie to radarr
func (s *Service) AddMovie(ctx context.Context, in ContentInfo, path string) error {

	raw, ok := in.raw.(radarr.Movie)
	if !ok {
		return fmt.Errorf("failed to get raw radarr data")
	}
	_, err := s.Radarr.AddMovie(
		&radarr.AddMovieInput{
			Title:               in.Title,
			TitleSlug:           raw.TitleSlug,
			MinimumAvailability: raw.MinimumAvailability,
			RootFolderPath:      path,
			TmdbID:              raw.TmdbID,
			QualityProfileID:    0,
			ProfileID:           0,
			Year:                in.Year,
			Images:              raw.Images,
			AddOptions: &radarr.AddMovieOptions{
				SearchForMovie:             true,
				IgnoreEpisodesWithFiles:    false,
				IgnoreEpisodesWithoutFiles: false,
			},
			Tags:      raw.Tags,
			Monitored: true,
		})

	if err != nil {
		return err
	}

	return nil
}

// AddShow will add a show to sonarr
func (s *Service) AddShow(ctx context.Context, in ContentInfo, path string) error {

	raw, ok := in.raw.(sonarr.SeriesLookup)
	if !ok {
		return fmt.Errorf("failed to get raw sonarr data")
	}

	_, err := s.Sonarr.AddSeries(
		&sonarr.AddSeriesInput{
			TvdbID:            raw.TvdbID,
			QualityProfileID:  0,
			LanguageProfileID: 0,
			Tags:              raw.Tags,
			RootFolderPath:    path,
			Title:             in.Title,
			SeriesType:        raw.SeriesType,
			Seasons:           raw.Seasons,
			AddOptions: &sonarr.AddSeriesOptions{
				SearchForMissingEpisodes:     true,
				SearchForCutoffUnmetEpisodes: false,
				IgnoreEpisodesWithFiles:      false,
				IgnoreEpisodesWithoutFiles:   false,
			},
			SeasonFolder: false,
			Monitored:    true,
		})

	if err != nil {
		return err
	}

	return nil
}
