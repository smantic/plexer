package service

import "context"

type DiskSpaceResult struct {
	// Total size of available disk space in bytes
	TotalCapacity int
	// UsedCapacity is the total amount of disk space in bytes used already
	UsedCapacity int
	// FreeSpace is the total amount of disk space in bytes left over
	FreeSpace int

	// Paths contains all of the paths were we are storing content
	Paths []string

	Err error
}

type QueuResult struct {
	Queue []QueueItem
	Err   error
}

type QueueItem struct {
	ContentType string
	Title       string
	Size        int
	Quality     int

	Status                  string
	TimeLeft                string
	EstimatedCompletionTime string

	Indexer        string
	DownloadClient string
}

func (s *Service) GetDiskSpaceInfo(ctx context.Context) <-chan DiskSpaceResult {

	c := make(chan DiskSpaceResult)

	go func() {
		disks, err := s.Radarr.DiskSpace(ctx)
		result := DiskSpaceResult{Err: err}

		for _, d := range disks {
			result.FreeSpace += d.FreeSpace
			result.TotalCapacity += d.TotalSpace
			result.UsedCapacity = result.UsedCapacity + (d.TotalSpace - d.FreeSpace)
			result.Paths = append(result.Paths, d.Path)
		}
		c <- result
		close(c)
	}()

	return c
}

func (s *Service) GetQueue(ctx context.Context) <-chan QueuResult {
	c := make(chan QueuResult)

	go func() {
		details, err := s.Radarr.QueueDetails(ctx)
		result := QueuResult{Err: err}

		items := make([]QueueItem, 0, len(details))
		for _, d := range details {
			item := QueueItem{
				ContentType:             CONTENT_MOVIE,
				Title:                   d.Title,
				Size:                    d.Size,
				Quality:                 d.Quality.Quality.Resolution,
				Status:                  d.Status,
				TimeLeft:                d.TimeLeft,
				EstimatedCompletionTime: d.EstimatedCompletedTime,
				Indexer:                 d.Indexer,
				DownloadClient:          d.DownloadClient,
			}
			items = append(items, item)
		}
		result.Queue = items
		c <- result
		close(c)
	}()

	return c
}
