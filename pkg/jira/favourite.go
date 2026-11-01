package jira

import (
	"context"
	"encoding/json"
	"net/http"
)

// Filter is a saved Jira filter.
type Filter struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	JQL   string `json:"jql"`
	Owner struct {
		Name        string `json:"name,omitempty"`
		DisplayName string `json:"displayName"`
	} `json:"owner"`
}

// FavouriteFilters fetches the user's starred filters from v2 /filter/favourite,
// which Jira Cloud, Server and Data Center all serve.
func (c *Client) FavouriteFilters() ([]*Filter, error) {
	res, err := c.GetV2(context.Background(), "/filter/favourite", nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, formatUnexpectedResponse(res)
	}

	var out []*Filter

	err = json.NewDecoder(res.Body).Decode(&out)

	return out, err
}
