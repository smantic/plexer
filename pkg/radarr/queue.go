package radarr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Item struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Quality struct {
	Id      int
	Quality struct {
		Id         int    `json:"id"`
		Name       string `json:"name"`
		Resolution int    `json:"resolution"`
		Modifier   string `json:"modifier"`
	} `json:"quality"`
	Revision struct {
		Version  int  `json:"version"`
		Real     int  `json:"real"`
		IsRepack bool `json:"is_repack"`
	} `json:"revision"`
}

type CustomFormat struct {
	Id   int    `json:"id"`
	Name string `json:"name"`

	IncludeCustomFormatWhenRenaming bool `json:"includeCustomFormatWhenRenaming"`

	Specifications struct {
		Name                string `json:"name"`
		Implementation      string `json:"implementation"`
		ImpelementationName string `json:"impelementationMame"`
		InfoLink            string `json:"infoLink"`
		Negate              string `json:"negate"`
		Required            bool   `json:"required"`
		Fields              []struct {
			Order    int    `json:"order"`
			Name     string `json:"name"`
			Label    string `json:"label"`
			HelpText string `json:"helpText"`
			Value    string `json:"value"`
			Type     string `json:"type"`
			Advanced bool   `json:"advanced"`
		} `json:"fields"`
	} `json:"specifications"`
}

type QueueItemInfo struct {
	Languages     []Item         `json:"languages"`
	Quality       Quality        `json:"quality"`
	CustomFormats []CustomFormat `json:"customFormats"`

	Size     int    `json:"size"`
	Title    string `json:"title"`
	SizeLeft int    `json:"size_left"`
	TimeLeft string `json:"time_left"`
	Status   string `json:"status"`

	EstimatedCompletedTime string `json:"estimated_completed_time"`
	TrackedDownloadStatus  string `json:"tracked_download_status"`
	ErrorMessage           string `json:"error_message"`
	DownloadID             string `json:"download_id"`
	Protocol               string `json:"protocol"`
	DownloadClient         string `json:"download_client"`
	Indexer                string `json:"indexer"`
	OutputPath             string `json:"output_path"`

	StatusMessages []struct {
		Title    string        `json:"title"`
		Messages []interface{} `json:"messages"`
	} `json:"status_messages"`
}

func (c *Client) QueueDetails(ctx context.Context) ([]QueueItemInfo, error) {

	result := []QueueItemInfo{}

	url := fmt.Sprintf("%s/queue/details", c.BaseURL)
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
