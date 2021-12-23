package service

import "context"

type FolderInfo struct {
	Path      string
	FreeSpace int
}

// GetRootFolders gets all of the unique root folders.
// prbly will go unused
func (s *Service) GetRootFolderPaths(ctx context.Context) ([]string, error) {

	m := make(map[string]struct{})
	rfolders, err := s.Radarr.GetRootFolders()
	if err != nil {
		return nil, err
	}

	// dedupe folders
	for _, f := range rfolders {
		m[f.Path] = struct{}{}
	}

	sfolders, err := s.Sonarr.GetRootFolders()
	if err != nil {
		return nil, err
	}

	// dedupe folders
	for _, f := range sfolders {
		m[f.Path] = struct{}{}
	}

	result := make([]string, 0, len(m))
	for path := range m {
		result = append(result, path)
	}

	return result, nil
}

// UpdateRootFolders will update our cached root folders
func (s *Service) GetRootFolderInfos(ctx context.Context) ([]FolderInfo, error) {

	rf, err := s.Radarr.GetRootFolders()
	if err != nil {
		return nil, err
	}

	sf, err := s.Sonarr.GetRootFolders()
	if err != nil {
		return nil, err
	}

	result := make([]FolderInfo, 0, len(rf)+len(sf))
	for _, f := range rf {
		info := FolderInfo{
			Path:      f.Path,
			FreeSpace: int(f.FreeSpace),
		}
		result = append(result, info)
	}

	for _, f := range sf {
		info := FolderInfo{
			Path:      f.Path,
			FreeSpace: int(f.FreeSpace),
		}
		result = append(result, info)
	}

	return result, nil
}
