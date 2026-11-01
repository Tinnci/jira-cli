package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// SearchRawV2 returns the v2 /search response body unchanged, including total and
// every requested field (custom fields too). fields is a comma-separated list or empty
// for Jira's default navigable fields.
func (c *Client) SearchRawV2(jql string, from, limit uint, fields string) ([]byte, error) {
	path := fmt.Sprintf("/search?jql=%s&startAt=%d&maxResults=%d", url.QueryEscape(jql), from, limit)
	if fields != "" {
		path += "&fields=" + url.QueryEscape(fields)
	}
	res, err := c.GetV2(context.Background(), path, nil)
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
	return io.ReadAll(res.Body)
}

// CountV2 returns how many issues match jql, reading no issues.
func (c *Client) CountV2(jql string) (int, error) {
	body, err := c.SearchRawV2(jql, 0, 0, "key")
	if err != nil {
		return 0, err
	}
	var out struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return 0, err
	}
	return out.Total, nil
}
