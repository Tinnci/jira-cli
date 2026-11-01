package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReleaseV2ReadsServerVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/project/PLV/versions", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"1","name":"V12","archived":false,"released":true,"projectId":7,"releaseDate":"2026-04-09"}]`))
	}))
	defer server.Close()

	client := NewClient(Config{Server: server.URL}, WithTimeout(3*time.Second))
	versions, err := client.ReleaseV2("PLV")
	assert.NoError(t, err)
	assert.Equal(t, []*ProjectVersion{{ID: "1", Name: "V12", Released: true, ProjectID: 7, ReleaseDate: "2026-04-09"}}, versions)
}

func TestFavouriteFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/filter/favourite", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":"10100","name":"PLV 严重","jql":"project = PLV AND priority = High","owner":{"name":"a","displayName":"A"}}]`))
	}))
	defer server.Close()

	client := NewClient(Config{Server: server.URL}, WithTimeout(3*time.Second))
	filters, err := client.FavouriteFilters()
	assert.NoError(t, err)
	assert.Len(t, filters, 1)
	assert.Equal(t, "project = PLV AND priority = High", filters[0].JQL)
	assert.Equal(t, "A", filters[0].Owner.DisplayName)
}

func TestSearchRawV2KeepsTheBodyAndCountReadsTotal(t *testing.T) {
	body := `{"startAt":0,"maxResults":2,"total":1314,"issues":[{"key":"PLV-1","fields":{"customfield_15121":{"value":"严重"}}}]}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/search", r.URL.Path)
		q := r.URL.Query()
		assert.Equal(t, "project = PLV", q.Get("jql"))
		if q.Get("maxResults") == "0" {
			assert.Equal(t, "key", q.Get("fields"))
			_, _ = w.Write([]byte(`{"total":1314,"issues":[]}`))
			return
		}
		assert.Equal(t, "5", q.Get("startAt"))
		assert.Equal(t, "summary,customfield_15121", q.Get("fields"))
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	client := NewClient(Config{Server: server.URL}, WithTimeout(3*time.Second))
	raw, err := client.SearchRawV2("project = PLV", 5, 2, "summary,customfield_15121")
	assert.NoError(t, err)
	assert.JSONEq(t, body, string(raw))
	var parsed map[string]any
	assert.NoError(t, json.Unmarshal(raw, &parsed))

	total, err := client.CountV2("project = PLV")
	assert.NoError(t, err)
	assert.Equal(t, 1314, total)
}

func TestSearchRawV2ReportsUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"errorMessages":["Field 'x' does not exist"]}`))
	}))
	defer server.Close()

	client := NewClient(Config{Server: server.URL}, WithTimeout(3*time.Second))
	_, err := client.SearchRawV2("x = 1", 0, 1, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
