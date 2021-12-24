package service

import "context"

type FreeSpace struct {
	// FreeSpace is the total amount of disk space in bytes left over
	FreeSpace int64

	// Paths contains all of the paths were we are storing content
	Paths []string
}

type QueueItem struct {
	ContentType ContentType
	Title       string
	Size        float64
	Quality     int

	Status   string
	TimeLeft string

	Indexer        string
	DownloadClient string
}

func (s *Service) GetFreeSpace(ctx context.Context) (FreeSpace, error) {

	result := FreeSpace{}
	folders, err := s.Radarr.GetRootFolders()
	if err != nil {
		return result, err
	}

	for _, f := range folders {
		result.FreeSpace += f.FreeSpace
		result.Paths = append(result.Paths, f.Path)
	}
	//TODO: how does this work with other services.

	return result, nil
}

func (s *Service) GetQueue(ctx context.Context) ([]QueueItem, error) {

	result := make([]QueueItem, 0, 20)

	rQ, err := s.Radarr.GetQueue(10, 0)
	if err != nil {
		return result, err
	}

	sQ, err := s.Sonarr.GetQueue(10, 0)
	if err != nil {
		return result, err
	}

	for _, s := range sQ.Records {
		item := QueueItem{
			ContentType:    CONTENT_SHOW,
			Title:          s.Title,
			Size:           s.Size,
			Quality:        s.Quality.ID,
			Status:         s.Status,
			TimeLeft:       s.Timeleft,
			Indexer:        s.Indexer,
			DownloadClient: s.DownloadClient,
		}
		result = append(result, item)
	}

	for _, r := range rQ.Records {
		item := QueueItem{
			ContentType:    CONTENT_MOVIE,
			Title:          r.Title,
			Size:           r.Size,
			Quality:        r.Quality.Quality.Resolution,
			Status:         r.Status,
			TimeLeft:       r.Timeleft,
			Indexer:        r.Indexer,
			DownloadClient: r.DownloadClient,
		}
		result = append(result, item)
	}

	return result, nil
}
