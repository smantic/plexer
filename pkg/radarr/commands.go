package radarr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Command string

const (
	//Trigger an update of Radarr
	CommandApplicationUpdate Command = "ApplicationUpdate"
	//Trigger a backup routine
	CommandBackup Command = "Backup"
	//Trigger a system health check
	CommandCheckHealth Command = "CheckHealth"
	// Triggers the removal of all blocklisted movies
	CommandClearBlocklist Command = "ClearBlocklist"
	// Trigger a recycle bin cleanup check
	CommandCleanUpRecycleBin Command = "CleanUprecycleBin"
	// Triggers the removal of all Info/Debug/Trace log files
	CommandDeleteLogFiles Command = "DeleteLogFiles"
	// Triggers the removal of all Update log files
	CommandDeleteUpdateLogFiles Command = "DeleteUpdateLogFiles"
	// Triggers the scan of downloaded movies
	CommandDownloadedMoviesScan Command = "DownloadMovieScan"
	// Triggers a search of all missing movies
	CommandMissingMoviesSearch Command = "MissingMoviesSearch"
	// Triggers the scan of monitored downloads
	CommandRefreshMonitoredDownloads Command = "RefreshMonitoredDownloads"
	// Trigger a refresh / scan of library
	CommandRefreshMovie Command = "RefreshMovie"
)

type CommandPayload struct {
	Name Command `json:"name"`

	// MovieIds are specified for the Refresh movie command
	MovieIds []int `json:"movies,omitempty"`
}

func (c *Client) PostCommand(ctx context.Context, req CommandPayload) error {

	url := fmt.Sprintf("%s/command", c.BaseURL)
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("X-Api-key", c.Apikey)

	response, err := c.Http.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return nil
	default:
		if c.Debug {
			bytes, _ := ioutil.ReadAll(response.Body)
			fmt.Printf("response: %s\n", string(bytes))
		}
		return fmt.Errorf("received %d from radarr", response.StatusCode)
	}
}
