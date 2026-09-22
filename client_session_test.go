package jmap_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestDoMarksSessionStale(t *testing.T) {
	var apiURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/session":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{
				"capabilities": {"urn:ietf:params:jmap:core": {}},
				"accounts": {},
				"primaryAccounts": {},
				"username": "u",
				"apiUrl": %q,
				"downloadUrl": "/download/{accountId}/{blobId}/{name}",
				"uploadUrl": "/upload/{accountId}",
				"eventSourceUrl": "/event",
				"state": "A"
			}`, apiURL)
		case r.Method == http.MethodPost && r.URL.Path == "/api":
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"methodResponses":[],"sessionState":"B"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	apiURL = srv.URL + "/api"

	c := &jmap.Client{
		HttpClient:      srv.Client(),
		SessionEndpoint: srv.URL + "/session",
	}
	require.NoError(t, c.Authenticate(t.Context()))
	require.Equal(t, "A", c.Session.State)
	require.False(t, c.SessionStale())

	_, err := c.Do(t.Context(), &jmap.Request{})
	require.NoError(t, err)
	require.True(t, c.SessionStale())
}
