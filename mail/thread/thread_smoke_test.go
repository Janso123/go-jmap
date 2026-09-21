package thread_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/thread"
	"github.com/stretchr/testify/require"
)

func TestThreadMethodNames(t *testing.T) {
	require.Equal(t, "Thread/get", (&thread.Get{}).Name())
	require.Equal(t, "Thread/changes", (&thread.Changes{}).Name())
}

func TestThreadGetRequiresMail(t *testing.T) {
	require.Equal(t, []jmap.URI{mail.URI}, (&thread.Get{}).Requires())
}

func TestThreadGetResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["Thread/get",{"accountId":"u1","state":"s0","list":[],"notFound":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*thread.GetResponse)
	require.True(t, ok)
}
