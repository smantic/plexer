package radarr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiskSpace is the response we get from radarr disk space endpoint
// {
//    "path": "D:\\",
//    "label": "DrivePool",
//    "freeSpace": 16187217043456,
//    "totalSpace": 56009755148288
//  },
type DiskSpace struct {
	Path       string `json:"path"`
	Label      string `json:"label"`
	FreeSpace  int    `json:"freeSpace"`
	TotalSpace int    `json:"totalSpace"`
}

func (c *Client) DiskSpace(ctx context.Context) ([]DiskSpace, error) {

	result := []DiskSpace{}

	url := fmt.Sprintf("%s/diskspace", c.BaseURL)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	response, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
