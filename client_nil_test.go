package jmap_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestNilHTTPClientDoesNotPanic(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"capabilities":{},"apiUrl":"/api","downloadUrl":"/d","uploadUrl":"/u","eventSourceUrl":"/e","state":"s"}`)
	}))
	defer srv.Close()
	c := &jmap.Client{SessionEndpoint: srv.URL} // HttpClient nil
	require.NotPanics(t, func() { _ = c.Authenticate(t.Context()) })
}

func TestPrimaryAccount(t *testing.T) {
	c := &jmap.Client{}
	_, err := c.PrimaryAccount(jmap.CoreURI)
	require.Error(t, err)
	require.Contains(t, err.Error(), "session not loaded")

	c.Session = &jmap.Session{
		PrimaryAccounts: map[jmap.URI]jmap.ID{
			"urn:ietf:params:jmap:mail": "a1",
		},
	}
	id, err := c.PrimaryAccount("urn:ietf:params:jmap:mail")
	require.NoError(t, err)
	require.Equal(t, jmap.ID("a1"), id)

	_, err = c.PrimaryAccount("urn:ietf:params:jmap:calendar")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no primary account")
}

func TestTrustedHosts(t *testing.T) {
	var hops int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hops++
		switch hops {
		case 1:
			http.Redirect(w, r, "http://evil.example/session", http.StatusFound)
		default:
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"capabilities":{},"apiUrl":"/api","downloadUrl":"/d","uploadUrl":"/u","eventSourceUrl":"/e","state":"s"}`)
		}
	}))
	defer srv.Close()

	c := jmap.NewClient(srv.URL, jmap.WithTrustedHosts(hostOf(srv.URL)))
	err := c.Authenticate(t.Context())
	require.Error(t, err)
	require.Contains(t, err.Error(), "untrusted host")
}

func TestTrustedHostsAllowsSameHost(t *testing.T) {
	var hops int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hops++
		switch {
		case hops == 1 && r.URL.Path == "/start":
			http.Redirect(w, r, "/session", http.StatusFound)
		default:
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"capabilities":{},"apiUrl":"/api","downloadUrl":"/d","uploadUrl":"/u","eventSourceUrl":"/e","state":"s"}`)
		}
	}))
	defer srv.Close()

	c := jmap.NewClient(srv.URL+"/start", jmap.WithTrustedHosts(hostOf(srv.URL)))
	require.NoError(t, c.Authenticate(t.Context()))
	require.NotNil(t, c.Session)
}

func hostOf(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u.Host
}
