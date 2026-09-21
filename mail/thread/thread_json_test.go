package thread_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/thread"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThreadJSONRoundTrip(t *testing.T) {
	th := thread.Thread{
		ID:       "th1",
		EmailIDs: []jmap.ID{"e1", "e2", "e3"},
	}

	data, err := json.Marshal(th)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"id": "th1",
		"emailIds": ["e1", "e2", "e3"]
	}`, string(data))

	var got thread.Thread
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, th, got)
}
