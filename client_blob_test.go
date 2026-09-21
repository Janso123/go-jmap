package jmap_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestUploadUsesContentType(t *testing.T) {
	var gotCT string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCT = r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"accountId":"a","blobId":"b","type":"text/plain","size":1}`)
	}))
	defer srv.Close()

	c := &jmap.Client{
		HttpClient:      srv.Client(),
		SessionEndpoint: "http://unused",
		Session: &jmap.Session{
			UploadURL: srv.URL + "/upload/{accountId}/",
		},
	}
	_, err := c.Upload(t.Context(), "a", strings.NewReader("x"), "text/plain")
	require.NoError(t, err)
	require.Equal(t, "text/plain", gotCT)
}

func TestDownloadOptions(t *testing.T) {
	var gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "blob-data")
	}))
	defer srv.Close()

	c := &jmap.Client{
		HttpClient:      srv.Client(),
		SessionEndpoint: "http://unused",
		Session: &jmap.Session{
			DownloadURL: srv.URL + "/dl/{accountId}/{blobId}/{name}?accept={type}",
		},
	}
	rc, err := c.Download(t.Context(), "acc", "blob1", jmap.DownloadOptions{
		Name: "file.txt",
		Type: "text/plain",
	})
	require.NoError(t, err)
	defer rc.Close()

	require.Equal(t, "/dl/acc/blob1/file.txt", gotPath)
	require.Equal(t, "accept=text%2Fplain", gotQuery)
}

func TestDownloadEncodesTemplateVars(t *testing.T) {
	var gotPath, gotRawQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		gotRawQuery = r.URL.RawQuery
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
	rc, err := c.Download(t.Context(), "acc", "blob/1", jmap.DownloadOptions{
		Name: "a b?x#y",
		Type: "application/octet-stream",
	})
	require.NoError(t, err)
	rc.Close()
	require.NotContains(t, gotPath, "?")
	require.NotContains(t, gotPath, "#")
	require.Contains(t, gotPath, "a%20b")
	require.Contains(t, gotPath, "%3F")
	require.Contains(t, gotPath, "%23")
	require.NotContains(t, gotRawQuery, "application/octet-stream")
	require.Contains(t, gotRawQuery, "application%2Foctet-stream")
}
