package thread_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/thread"
	"github.com/stretchr/testify/require"
)

func TestThreadGetEmbedsKit(t *testing.T) {
	m := &thread.Get{Get: jmap.Get[thread.Thread]{Account: "a1"}}
	require.Equal(t, "Thread/get", m.Name())
	b, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"a1"}`, string(b))

	var resp thread.GetResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a1","state":"s0","list":[{"id":"th1","emailIds":["e1"]}],"notFound":[]}`), &resp))
	require.Equal(t, []thread.Thread{{ID: "th1", EmailIDs: []jmap.ID{"e1"}}}, resp.List)
}
