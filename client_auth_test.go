package jmap_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestAuthenticateSetsAcceptJSON(t *testing.T) {
	t.Parallel()
	var accept string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept = r.Header.Get("Accept")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{
			"capabilities":{"urn:ietf:params:jmap:core":{}},
			"accounts":{},"primaryAccounts":{},
			"username":"u","apiUrl":"http://example/api",
			"downloadUrl":"http://example/d","uploadUrl":"http://example/u",
			"eventSourceUrl":"http://example/e","state":"A"
		}`)
	}))
	defer srv.Close()

	c := &jmap.Client{HttpClient: srv.Client(), SessionEndpoint: srv.URL}
	require.NoError(t, c.Authenticate(t.Context()))
	require.Equal(t, "application/json", accept)
	require.Equal(t, "A", c.Session.State)
}

func TestAuthenticateProblemJSONVariants(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		body string
		typ  string
	}{
		{"notJSON", `{"type":"urn:ietf:params:jmap:error:notJSON","detail":"x","status":400}`, jmap.ErrTypeNotJSON},
		{"notRequest", `{"type":"urn:ietf:params:jmap:error:notRequest","detail":"x","status":400}`, jmap.ErrTypeNotRequest},
		{"unknownCapability", `{"type":"urn:ietf:params:jmap:error:unknownCapability","detail":"x","status":400}`, jmap.ErrTypeUnknownCapability},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(400)
				io.WriteString(w, tt.body)
			}))
			defer srv.Close()
			c := &jmap.Client{HttpClient: srv.Client(), SessionEndpoint: srv.URL}
			err := c.Authenticate(t.Context())
			require.Error(t, err)
			var re *jmap.RequestError
			require.ErrorAs(t, err, &re)
			require.Equal(t, tt.typ, re.Type)
			require.Equal(t, 400, re.Status)
		})
	}
}

func TestRefreshSessionClearsStaleAndFiresCallback(t *testing.T) {
	t.Parallel()
	var apiURL string
	var sessionHits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/session":
			sessionHits.Add(1)
			state := "A"
			if sessionHits.Load() > 1 {
				state = "B"
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{
				"capabilities":{"urn:ietf:params:jmap:core":{}},
				"accounts":{},"primaryAccounts":{},
				"username":"u","apiUrl":%q,
				"downloadUrl":"http://example/d","uploadUrl":"http://example/u",
				"eventSourceUrl":"http://example/e","state":%q
			}`, apiURL, state)
		case r.Method == http.MethodPost && r.URL.Path == "/api":
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"methodResponses":[],"sessionState":"B"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	apiURL = srv.URL + "/api"

	var called int
	c := &jmap.Client{
		HttpClient:      srv.Client(),
		SessionEndpoint: srv.URL + "/session",
		OnSessionChange: func(old, neu *jmap.Session) {
			called++
			require.Equal(t, "A", old.State)
			require.Equal(t, "B", neu.State)
		},
	}
	require.NoError(t, c.Authenticate(t.Context()))
	_, err := c.Do(t.Context(), &jmap.Request{})
	require.NoError(t, err)
	require.True(t, c.SessionStale())

	require.NoError(t, c.RefreshSession(t.Context()))
	require.False(t, c.SessionStale())
	require.Equal(t, "B", c.Session.State)
	require.Equal(t, 1, called)
	require.GreaterOrEqual(t, sessionHits.Load(), int32(2))
}
