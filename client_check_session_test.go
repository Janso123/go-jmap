package jmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckSessionKeepsAbsoluteTemplateBraces(t *testing.T) {
	const upload = "https://api.example/upload/{accountId}/{blobId}"
	s := &Session{
		APIURL:    "https://api.example/jmap",
		UploadURL: upload,
	}
	require.NoError(t, checkSession("https://api.example/jmap/session", s))
	require.Equal(t, upload, s.UploadURL)
}

func TestCheckSessionRejectsEncodedReservedUploadTemplate(t *testing.T) {
	s := &Session{
		APIURL:    "https://api.example/jmap",
		UploadURL: "/u/%7B+path%7D",
	}
	require.Error(t, checkSession("https://api.example/jmap/session", s))
}

func TestCheckSessionResolvesRelativeTemplateBraces(t *testing.T) {
	s := &Session{
		APIURL:    "https://api.example/jmap",
		UploadURL: "/upload/{accountId}",
	}
	require.NoError(t, checkSession("https://api.example/jmap/session", s))
	require.Equal(t, "https://api.example/upload/{accountId}", s.UploadURL)
}

func TestCheckSessionRejectsReservedUploadTemplate(t *testing.T) {
	s := &Session{
		APIURL:         "https://api.example/jmap",
		DownloadURL:    "https://api.example/d/{accountId}/{blobId}/{name}",
		UploadURL:      "https://api.example/{+accountId}",
		EventSourceURL: "https://api.example/e",
	}
	err := checkSession("https://api.example/session", s)
	require.Error(t, err)
}
