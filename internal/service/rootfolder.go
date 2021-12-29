package service

import (
	"context"
	"fmt"

	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
)

type FolderInfo struct {
	ID          int
	Path        string
	FreeSpace   int64
	ContentType ContentType
}

func (s *Service) getAllUniqueFolders(ctx context.Context) (map[string]FolderInfo, error) {

	m := make(map[string]FolderInfo)
	rf, err := s.Radarr.GetRootFolders()
	if err != nil {
		return nil, fmt.Errorf("failed to find sonarr root folders: %w", err)
	}

	for _, f := range rf {
		m[f.Path] = infoFromRadarr(f)
	}

	sf, err := s.Sonarr.GetRootFolders()
	if err != nil {
		return nil, fmt.Errorf("failed to find sonarr root folders: %w", err)
	}

	for _, f := range sf {
		m[f.Path] = infoFromSonarr(f)
	}

	return m, nil
}

// GetTotalFreeSpace gets the total free space available in bytes.
func (s *Service) GetTotalFreeSpace(ctx context.Context) (int64, error) {

	m, err := s.getAllUniqueFolders(ctx)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, info := range m {
		total += info.FreeSpace
	}

	return total, nil
}

// GetRootFolders gets all of the unique root folders.
// prbly will go unused.
func (s *Service) GetRootFolderPaths(ctx context.Context) ([]string, error) {

	m, err := s.getAllUniqueFolders(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(m))
	for path := range m {
		result = append(result, path)
	}

	return result, nil
}

// GetRootFolderInfos will get the root folder info for each service.
func (s *Service) GetRootFolderInfos(ctx context.Context, kind ContentType) ([]FolderInfo, error) {

	switch kind {
	case CONTENT_MOVIE:
		rf, err := s.Radarr.GetRootFolders()
		if err != nil {
			return nil, fmt.Errorf("failed to find sonarr root folders: %w", err)
		}

		result := make([]FolderInfo, 0, len(rf))
		for _, f := range rf {
			result = append(result, infoFromRadarr(f))
		}

		return result, nil
	case CONTENT_SHOW:
		sf, err := s.Sonarr.GetRootFolders()
		if err != nil {
			return nil, fmt.Errorf("failed to find sonarr root folders: %w", err)
		}

		result := make([]FolderInfo, 0, len(sf))
		for _, f := range sf {
			result = append(result, infoFromSonarr(f))
		}

		return result, nil
	default:
		rf, err := s.Radarr.GetRootFolders()
		if err != nil {
			return nil, fmt.Errorf("failed to find sonarr root folders: %w", err)
		}

		sf, err := s.Sonarr.GetRootFolders()
		if err != nil {
			return nil, fmt.Errorf("failed to find sonarr root folders: %w", err)
		}

		result := make([]FolderInfo, 0, len(rf)+len(sf))
		for _, f := range rf {
			result = append(result, infoFromRadarr(f))
		}

		for _, f := range sf {
			result = append(result, infoFromSonarr(f))
		}

		return result, nil
	}
}

func infoFromSonarr(f *sonarr.RootFolder) FolderInfo {
	return FolderInfo{
		ID:          int(f.ID),
		Path:        f.Path,
		FreeSpace:   f.FreeSpace,
		ContentType: CONTENT_SHOW,
	}
}

func infoFromRadarr(f *radarr.RootFolder) FolderInfo {
	return FolderInfo{
		ID:          int(f.ID),
		Path:        f.Path,
		FreeSpace:   f.FreeSpace,
		ContentType: CONTENT_MOVIE,
	}
}
