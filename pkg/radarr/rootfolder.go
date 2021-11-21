package radarr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RootFolderInfo struct {
	// Fully qualified path
	Path string `json:"path"`
	// FreeSpace in the root folder
	FreeSpace int64 `json:"freeSpace"`
	// idk
	UnmmapedFolders []interface{} `json:"unmmapedFolders"`
	Id              int           `json:"id"`
}

// GetRootFolder gets radarr's root folder
func (c *Client) GetRootFolder(ctx context.Context) ([]RootFolderInfo, error) {

	result := []RootFolderInfo{}
	url := fmt.Sprintf("%s/rootfolder", c.BaseURL)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("X-Api-key", c.Apikey)

	response, err := c.Http.Do(r)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
