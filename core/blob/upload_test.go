package blob_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadInvokeUsesBlobCapabilityAndSpecialDataKeys(t *testing.T) {
	req := &jmap.Request{}
	asBase64 := "aGVsbG8="

	id := req.Invoke(&blob.Upload{
		Account: "u1",
		Create: map[jmap.ID]*blob.UploadObject{
			"c1": {
				Data: []blob.DataSource{{AsBase64: &asBase64}},
				Type: "text/plain",
			},
		},
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"using":["urn:ietf:params:jmap:blob"]`)
	assert.Contains(t, string(data), `"Blob/upload"`)
	assert.Contains(t, string(data), `"data:asBase64":"aGVsbG8="`)
	assert.Contains(t, string(data), `"type":"text/plain"`)
}

func TestDataSourceUnmarshalSpecialKeys(t *testing.T) {
	var source blob.DataSource

	require.NoError(t, json.Unmarshal([]byte(`{"data:asText":"hello"}`), &source))
	require.NotNil(t, source.AsText)
	assert.Equal(t, "hello", *source.AsText)
	assert.Nil(t, source.AsBase64)

	require.NoError(t, json.Unmarshal([]byte(`{"data:asBase64":"aGVsbG8="}`), &source))
	require.NotNil(t, source.AsBase64)
	assert.Equal(t, "aGVsbG8=", *source.AsBase64)
}

func TestUploadResponseUnmarshal(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["Blob/upload",{"accountId":"u1","created":{"c1":{"id":"blob1","size":5,"type":"text/plain"}},"notCreated":{"c2":{"type":"invalidProperties"}}},"0"]]}`)

	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)

	inv := resp.Responses[0]
	assert.Equal(t, "Blob/upload", inv.Name)
	assert.Equal(t, "0", inv.CallID)

	methodResp, ok := inv.Args.(*blob.UploadResponse)
	require.True(t, ok)
	require.Contains(t, methodResp.Created, jmap.ID("c1"))
	assert.Equal(t, jmap.ID("blob1"), methodResp.Created["c1"].ID)
	assert.Equal(t, uint64(5), methodResp.Created["c1"].Size)
	assert.Equal(t, "text/plain", methodResp.Created["c1"].Type)
	require.Contains(t, methodResp.NotCreated, jmap.ID("c2"))
	assert.Equal(t, "invalidProperties", methodResp.NotCreated["c2"].Type)
}
