package jmap_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestTimeoutSurvivesBearer(t *testing.T) {
	t.Parallel()
	c := jmap.NewClient("https://example.com/jmap/session",
		jmap.WithTimeout(5*time.Second),
		jmap.WithBearer("tok"),
	)
	require.NotNil(t, c.HttpClient)
	require.Equal(t, 5*time.Second, c.HttpClient.Timeout)
	require.NotEqual(t, http.DefaultClient, c.HttpClient)
}

func TestTimeoutSurvivesBearerReverseOrder(t *testing.T) {
	t.Parallel()
	c := jmap.NewClient("https://example.com/jmap/session",
		jmap.WithBearer("tok"),
		jmap.WithTimeout(5*time.Second),
	)
	require.NotNil(t, c.HttpClient)
	require.Equal(t, 5*time.Second, c.HttpClient.Timeout)
}
