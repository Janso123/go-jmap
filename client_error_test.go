package jmap_test

import (
	"encoding/json/jsontext"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestDecodeProblemJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(400)
		io.WriteString(w, `{"type":"urn:ietf:params:jmap:error:limit","status":400,"title":"Limit","detail":"too big","limit":"maxSizeRequest","requestId":"r1"}`)
	}))
	defer srv.Close()

	// Exercise decodeHttpError via Do with a pre-set Session.
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
	require.Equal(t, jmap.ErrTypeLimit, re.Type)
	require.Equal(t, "Limit", re.Title)
	require.Equal(t, "r1", re.RequestID)
	require.NotNil(t, re.Limit)
}
