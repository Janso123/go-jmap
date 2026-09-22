package jmap_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json/jsontext"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

const redirectSessionJSON = `{"capabilities":{},"apiUrl":"/api","downloadUrl":"/d","uploadUrl":"/u","eventSourceUrl":"/e","state":"s"}`

func TestTrustedHostsRejectsSchemeDowngrade(t *testing.T) {
	t.Parallel()
	tlsSrv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://"+r.Host+"/session", http.StatusFound)
	}))
	t.Cleanup(tlsSrv.Close)

	c := jmap.NewClient(tlsSrv.URL, jmap.WithHTTPClient(tlsSrv.Client()), jmap.WithTrustedHosts(hostOf(tlsSrv.URL)))
	err := c.Authenticate(t.Context())
	require.Error(t, err)
	require.Nil(t, c.Session)
	require.Contains(t, err.Error(), "untrusted")
}

func TestAuthenticateRejectsSchemeDowngradeWithoutTrustedHosts(t *testing.T) {
	t.Parallel()
	tlsSrv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://"+r.Host+"/session", http.StatusFound)
	}))
	t.Cleanup(tlsSrv.Close)

	c := jmap.NewClient(tlsSrv.URL, jmap.WithHTTPClient(tlsSrv.Client()))
	err := c.Authenticate(t.Context())
	require.Error(t, err)
	require.Nil(t, c.Session)
	require.Contains(t, err.Error(), "untrusted")
}

func TestTrustedHostsRejectsDifferentPort(t *testing.T) {
	t.Parallel()
	var saw bool
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		saw = true
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, redirectSessionJSON)
	}))
	t.Cleanup(other.Close)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL+"/session", http.StatusFound)
	}))
	t.Cleanup(srv.Close)

	// Hostname matches; port must not be ignored.
	c := jmap.NewClient(srv.URL, jmap.WithTrustedHosts("127.0.0.1"))
	err := c.Authenticate(t.Context())
	require.Error(t, err)
	require.Nil(t, c.Session)
	require.False(t, saw)
	require.Contains(t, err.Error(), "untrusted")
}

func TestDoRefusesSameOrigin302POST(t *testing.T) {
	t.Parallel()
	var methods []string
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		if r.Method == http.MethodPost {
			http.Redirect(w, r, srv.URL+"/next", http.StatusFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"methodResponses":[],"sessionState":"A"}`)
	}))
	t.Cleanup(srv.Close)

	c := &jmap.Client{
		HttpClient: srv.Client(),
		Session: &jmap.Session{
			APIURL: srv.URL + "/api",
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: jsontext.Value(`{}`),
			},
		},
	}
	_, err := c.Do(t.Context(), &jmap.Request{})
	require.Error(t, err)
	require.Equal(t, []string{http.MethodPost}, methods)
	require.NotContains(t, methods, http.MethodGet)
}

func TestDoFollowsSameOrigin307WithBody(t *testing.T) {
	t.Parallel()
	var methods []string
	var bodies [][]byte
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		b, _ := io.ReadAll(r.Body)
		bodies = append(bodies, b)
		if r.URL.Path == "/api" {
			http.Redirect(w, r, srv.URL+"/next", http.StatusTemporaryRedirect)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"methodResponses":[],"sessionState":"A"}`)
	}))
	t.Cleanup(srv.Close)

	c := &jmap.Client{
		HttpClient: srv.Client(),
		Session: &jmap.Session{
			APIURL: srv.URL + "/api",
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: jsontext.Value(`{}`),
			},
		},
	}
	resp, err := c.Do(t.Context(), &jmap.Request{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []string{http.MethodPost, http.MethodPost}, methods)
	require.NotEmpty(t, bodies)
	require.Len(t, bodies, 2)
	require.JSONEq(t, string(bodies[0]), string(bodies[1]))
}
