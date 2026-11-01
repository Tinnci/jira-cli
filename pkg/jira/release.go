package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Release fetches response from v3 /project/{projectIdOrKey}/versions endpoint.
func (c *Client) Release(project string) ([]*ProjectVersion, error) {
	return c.release(project, apiVersion3)
}

// ReleaseV2 fetches response from v2 /project/{projectIdOrKey}/versions endpoint.
// Jira Server and Data Center serve only v2.
func (c *Client) ReleaseV2(project string) ([]*ProjectVersion, error) {
	return c.release(project, apiVersion2)
}

func (c *Client) release(project, ver string) ([]*ProjectVersion, error) {
	path := fmt.Sprintf("/project/%s/versions", project)

	var (
		res *http.Response
		err error
	)
	if ver == apiVersion2 {
		res, err = c.GetV2(context.Background(), path, nil)
	} else {
		res, err = c.Get(context.Background(), path, nil)
	}
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

	var out []*ProjectVersion

	err = json.NewDecoder(res.Body).Decode(&out)

	return out, err
}
