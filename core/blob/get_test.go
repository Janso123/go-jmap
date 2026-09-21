package blob_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInvoke(t *testing.T) {
	req := &jmap.Request{}
	offset := uint64(4)
	length := uint64(9)

	id := req.Invoke(&blob.Get{
		Account:    "u1",
		IDs:        []jmap.ID{"blob1"},
		Properties: []string{"data:asText", "digest:sha-256", "size"},
		Offset:     &offset,
		Length:     &length,
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:blob"],"methodCalls":[["Blob/get",{"accountId":"u1","ids":["blob1"],"properties":["data:asText","digest:sha-256","size"],"offset":4,"length":9},"0"]]}`,
		string(data))
}

func TestGetResponseUnmarshal(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["Blob/get",{"accountId":"u1","list":[{"id":"blob1","data:asText":"hello","digest:sha-256":"abc=","size":5}],"notFound":["missing"]},"0"]]}`)

	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)

	inv := resp.Responses[0]
	assert.Equal(t, "Blob/get", inv.Name)
	assert.Equal(t, "0", inv.CallID)

	methodResp, ok := inv.Args.(*blob.GetResponse)
	require.True(t, ok)
	require.Len(t, methodResp.List, 1)
	require.NotNil(t, methodResp.List[0].DataAsText)
	assert.Equal(t, "hello", *methodResp.List[0].DataAsText)
	assert.Equal(t, "abc=", methodResp.List[0].Digests["sha-256"])
	assert.Equal(t, uint64(5), methodResp.List[0].Size)
	assert.Equal(t, []jmap.ID{"missing"}, methodResp.NotFound)
}

func TestGetResultRoundTripPreservesEmptyDataAsBase64(t *testing.T) {
	var result blob.GetResult

	require.NoError(t, json.Unmarshal([]byte(`{"id":"blob1","data:asBase64":""}`), &result))

	data, err := json.Marshal(result)
	require.NoError(t, err)
	assert.JSONEq(t, `{"id":"blob1","data:asBase64":""}`, string(data))
}
