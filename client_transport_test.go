package jmap_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"encoding/json/jsontext"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core"
	"github.com/stretchr/testify/require"
)

func TestValidateURITemplateLevel1(t *testing.T) {
	require.NoError(t, jmap.ValidateURITemplateLevel1("https://x/{accountId}/{blobId}"))
	require.Error(t, jmap.ValidateURITemplateLevel1("https://x/{+path}"))
	require.Error(t, jmap.ValidateURITemplateLevel1("https://x/{a,b}"))
	require.Error(t, jmap.ValidateURITemplateLevel1("https://x/{a:3}"))
}

func TestExpandURITemplateLevel1EncodesReserved(t *testing.T) {
	got := jmap.ExpandURITemplateLevel1(
		"https://ex/dl/{accountId}/{blobId}/{name}?accept={type}",
		map[string]string{
			"accountId": "a",
			"blobId":    "b",
			"name":      "x&accountId=victim&name=evil.exe",
			"type":      "text/plain",
		},
	)
	require.NotContains(t, got, "accountId=victim")
	require.Contains(t, got, "x%26accountId%3Dvictim")
	require.Contains(t, got, "accept=text%2Fplain")
}

func TestBearerNotForwardedOnCrossOriginRedirect(t *testing.T) {
	var evilAuth string
	evil := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		evilAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"capabilities":{"urn:ietf:params:jmap:core":{}},"accounts":{},"primaryAccounts":{},"username":"u","apiUrl":"/api","downloadUrl":"/d","uploadUrl":"/u","eventSourceUrl":"/e","state":"A"}`)
	}))
	defer evil.Close()

	good := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, evil.URL, http.StatusFound)
	}))
	defer good.Close()

	c := jmap.NewClient(good.URL, jmap.WithHTTPClient(&http.Client{}), jmap.WithBearer("SUPER-SECRET-TOKEN"))
	_ = c.Authenticate(context.Background())
	require.Empty(t, evilAuth)
}

func TestAuthenticateRejectsCrossOriginSessionAndDoDoesNotSendBearer(t *testing.T) {
	const secret = "SUPER-SECRET-TOKEN"
	var evilAPIAuth string
	var evilGotPOST bool

	var evil *httptest.Server
	evil = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			evilGotPOST = true
			evilAPIAuth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"methodResponses":[],"sessionState":"A"}`)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{
			"capabilities":{"urn:ietf:params:jmap:core":{}},
			"accounts":{},"primaryAccounts":{},"username":"u",
			"apiUrl":"`+evil.URL+`/api",
			"downloadUrl":"`+evil.URL+`/d",
			"uploadUrl":"`+evil.URL+`/u",
			"eventSourceUrl":"`+evil.URL+`/e","state":"A"
		}`)
	}))
	defer evil.Close()

	good := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, evil.URL+"/session", http.StatusFound)
	}))
	defer good.Close()

	c := jmap.NewClient(good.URL, jmap.WithHTTPClient(&http.Client{}), jmap.WithBearer(secret))
	err := c.Authenticate(context.Background())
	require.Error(t, err)
	require.Nil(t, c.Session)

	_, _ = c.Do(context.Background(), &jmap.Request{})
	require.False(t, evilGotPOST)
	require.Empty(t, evilAPIAuth)
}

func TestAuthenticateRejectsCrossOriginSessionResourceURLs(t *testing.T) {
	const secret = "SUPER-SECRET-TOKEN"
	var otherAuth string
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		otherAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"methodResponses":[],"sessionState":"A"}`)
	}))
	defer other.Close()

	session := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{
			"capabilities":{"urn:ietf:params:jmap:core":{}},
			"accounts":{},"primaryAccounts":{},"username":"u",
			"apiUrl":"`+other.URL+`/api",
			"downloadUrl":"/d","uploadUrl":"/u","eventSourceUrl":"/e","state":"A"
		}`)
	}))
	defer session.Close()

	c := jmap.NewClient(session.URL, jmap.WithHTTPClient(&http.Client{}), jmap.WithBearer(secret))
	err := c.Authenticate(context.Background())
	require.Error(t, err)
	require.Nil(t, c.Session)
	require.Empty(t, otherAuth)
}

func TestAuthenticateRejectsNonHTTPSNonLoopback(t *testing.T) {
	c := jmap.NewClient("http://example.com/.well-known/jmap")
	err := c.Authenticate(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "https")
}

func TestHTTPErrorTypedStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(401)
		io.WriteString(w, "denied")
	}))
	defer srv.Close()

	c := &jmap.Client{
		HttpClient: srv.Client(),
		Session: &jmap.Session{
			APIURL: srv.URL,
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: jsontext.Value(`{}`),
			},
		},
	}
	_, err := c.Do(t.Context(), &jmap.Request{})
	require.Error(t, err)
	var he *jmap.HTTPError
	require.ErrorAs(t, err, &he)
	require.Equal(t, 401, he.Status)
	require.Contains(t, he.Body, "denied")
}

func TestProblemJSONWithoutTypeIsAboutBlank(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(429)
		io.WriteString(w, `{"status":429,"title":"Too Many Requests","detail":"slow down"}`)
	}))
	defer srv.Close()

	c := &jmap.Client{
		HttpClient: srv.Client(),
		Session: &jmap.Session{
			APIURL: srv.URL,
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: jsontext.Value(`{}`),
			},
		},
	}
	_, err := c.Do(t.Context(), &jmap.Request{})
	require.Error(t, err)
	var re *jmap.RequestError
	require.ErrorAs(t, err, &re)
	require.Equal(t, "about:blank", re.Type)
	require.Equal(t, 429, re.Status)
}

func TestDoRejectsCrossOriginRedirectResponse(t *testing.T) {
	evil := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"methodResponses":[],"sessionState":"EVIL"}`)
	}))
	defer evil.Close()

	good := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, evil.URL, http.StatusFound)
	}))
	defer good.Close()

	c := &jmap.Client{
		HttpClient: &http.Client{},
		Session: &jmap.Session{
			APIURL: good.URL,
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: jsontext.Value(`{}`),
			},
		},
	}
	resp, err := c.Do(t.Context(), &jmap.Request{})
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestDownloadRejectsQueryInjectionViaName(t *testing.T) {
	var rawQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := &jmap.Client{
		HttpClient:      srv.Client(),
		SessionEndpoint: "http://unused",
		Session: &jmap.Session{
			DownloadURL: srv.URL + "/dl/{accountId}/{blobId}/{name}?accept={type}",
		},
	}
	rc, err := c.Download(t.Context(), "acc", "b", jmap.DownloadOptions{
		Name: "x&accountId=victim",
		Type: "text/plain",
	})
	require.NoError(t, err)
	rc.Close()
	require.NotContains(t, rawQuery, "accountId=victim")
}

func TestRefreshSessionConcurrentNoRace(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"capabilities":{"urn:ietf:params:jmap:core":{}},"accounts":{},"primaryAccounts":{},"username":"u","apiUrl":"/api","downloadUrl":"/d","uploadUrl":"/u","eventSourceUrl":"/e","state":"A"}`)
	}))
	defer srv.Close()

	c := &jmap.Client{HttpClient: srv.Client(), SessionEndpoint: srv.URL}
	require.NoError(t, c.Authenticate(t.Context()))

	var wg sync.WaitGroup
	for range 16 {
		wg.Go(func() {
			_ = c.RefreshSession(t.Context())
		})
	}
	wg.Wait()
}

func TestRelativeAPIURLResolved(t *testing.T) {
	var apiPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{
				"capabilities":{"urn:ietf:params:jmap:core":{}},
				"accounts":{},"primaryAccounts":{},
				"username":"u","apiUrl":"/jmap/api/",
				"downloadUrl":"/dl/{accountId}/{blobId}/{name}?accept={type}",
				"uploadUrl":"/upload/{accountId}/",
				"eventSourceUrl":"/es","state":"A"
			}`)
			return
		}
		apiPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"methodResponses":[],"sessionState":"A"}`)
	}))
	defer srv.Close()

	c := &jmap.Client{HttpClient: srv.Client(), SessionEndpoint: srv.URL + "/session"}
	require.NoError(t, c.Authenticate(t.Context()))
	require.Contains(t, c.Session.APIURL, "/jmap/api/")
	require.True(t, len(c.Session.APIURL) > len("/jmap/api/"))
	_, err := c.Do(t.Context(), &jmap.Request{})
	require.NoError(t, err)
	require.Equal(t, "/jmap/api/", apiPath)
}

func TestDoAccepts40MiBWhenMaxSizeRequestIs50000000(t *testing.T) {
	const payload = 40 << 20
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"methodResponses":[],"sessionState":"`)
		io.WriteString(w, strings.Repeat("a", payload))
		io.WriteString(w, `"}`)
	}))
	defer srv.Close()

	c := &jmap.Client{
		HttpClient: srv.Client(),
		Session: &jmap.Session{
			APIURL: srv.URL,
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: jsontext.Value(`{"maxSizeRequest":50000000}`),
			},
			Capabilities: map[jmap.URI]jmap.Capability{
				jmap.CoreURI: &core.Core{MaxSizeRequest: 50000000},
			},
		},
	}
	resp, err := c.Do(t.Context(), &jmap.Request{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.SessionState, payload)
}

func TestWithMaxResponseBytesRejectsOverCap(t *testing.T) {
	const payload = 2 << 20
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"methodResponses":[],"sessionState":"`)
		io.WriteString(w, strings.Repeat("a", payload))
		io.WriteString(w, `"}`)
	}))
	defer srv.Close()

	c := jmap.NewClient(srv.URL,
		jmap.WithHTTPClient(srv.Client()),
		jmap.WithMaxResponseBytes(1<<20),
	)
	c.Session = &jmap.Session{
		APIURL: srv.URL,
		RawCapabilities: map[jmap.URI]jsontext.Value{
			jmap.CoreURI: jsontext.Value(`{}`),
		},
	}
	_, err := c.Do(t.Context(), &jmap.Request{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "exceeds")
}
