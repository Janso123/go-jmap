package blob_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&blob.Lookup{
		Account:   "u1",
		TypeNames: []string{"Mailbox", "Thread", "Email"},
		IDs:       []jmap.ID{"blob1", "missing"},
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:blob"],"methodCalls":[["Blob/lookup",{"accountId":"u1","typeNames":["Mailbox","Thread","Email"],"ids":["blob1","missing"]},"0"]]}`,
		string(data))
}

func TestLookupResponseUnmarshal(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["Blob/lookup",{"accountId":"u1","list":[{"id":"blob1","matchedIds":{"Mailbox":["m1"],"Thread":["t1"],"Email":["e1","e2"]}}],"notFound":["missing"]},"0"]]}`)

	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)

	inv := resp.Responses[0]
	assert.Equal(t, "Blob/lookup", inv.Name)
	assert.Equal(t, "0", inv.CallID)

	methodResp, ok := inv.Args.(*blob.LookupResponse)
	require.True(t, ok)
	require.Len(t, methodResp.List, 1)
	assert.Equal(t, jmap.ID("blob1"), methodResp.List[0].ID)
	assert.Equal(t, []jmap.ID{"m1"}, methodResp.List[0].MatchedIDs["Mailbox"])
	assert.Equal(t, []jmap.ID{"t1"}, methodResp.List[0].MatchedIDs["Thread"])
	assert.Equal(t, []jmap.ID{"e1", "e2"}, methodResp.List[0].MatchedIDs["Email"])
	assert.Equal(t, []jmap.ID{"missing"}, methodResp.NotFound)
}
