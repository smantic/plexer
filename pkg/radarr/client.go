package radarr

import "net/http"

type Client struct {
	BaseURL string
	Apikey  string
	Http    http.Client

	// Debug will print response body
	Debug bool
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {

	r.Header.Add("Accept", "application/json")
	r.Header.Add("X-Api-key", c.Apikey)

	return c.Http.Do(r)
}
