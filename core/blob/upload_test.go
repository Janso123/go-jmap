package blob_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadInvokeUsesBlobCapabilityAndSpecialDataKeys(t *testing.T) {
	req := &jmap.Request{}
	id := req.Invoke(&blob.Upload{
		Account: "u1",
		Create: map[jmap.ID]*blob.UploadObject{
			"c1": {
				Data: []blob.DataSource{{AsBase64: jmap.Some("aGVsbG8=")}},
				Type: jmap.Some("text/plain"),
			},
		},
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"using":["urn:ietf:params:jmap:blob"]`)
	assert.Contains(t, string(data), `"Blob/upload"`)
	assert.Contains(t, string(data), `"data:asBase64":"aGVsbG8="`)
	assert.Contains(t, string(data), `"type":"text/plain"`)
}

func TestDataSourceNullRemarshals(t *testing.T) {
	t.Parallel()
	var text blob.DataSource
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"data:asText":null}`), &text))
	require.True(t, text.AsText.IsNull())
	b, err := jsonv2.Marshal(&text)
	require.NoError(t, err)
	require.JSONEq(t, `{"data:asText":null}`, string(b))

	var b64 blob.DataSource
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"data:asBase64":null}`), &b64))
	require.True(t, b64.AsBase64.IsNull())
	b, err = jsonv2.Marshal(&b64)
	require.NoError(t, err)
	require.JSONEq(t, `{"data:asBase64":null}`, string(b))
}

func TestDataSourceUnmarshalSpecialKeys(t *testing.T) {
	var source blob.DataSource

	require.NoError(t, jsonv2.Unmarshal([]byte(`{"data:asText":"hello"}`), &source))
	text, ok := source.AsText.Value()
	require.True(t, ok)
	assert.Equal(t, "hello", text)
	assert.True(t, source.AsBase64.IsZero())

	require.NoError(t, jsonv2.Unmarshal([]byte(`{"data:asBase64":"aGVsbG8="}`), &source))
	b64, ok := source.AsBase64.Value()
	require.True(t, ok)
	assert.Equal(t, "aGVsbG8=", b64)
}

func TestUploadTypeNullRemarshals(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(blob.UploadObject{
		Data: []blob.DataSource{},
		Type: jmap.Null[string](),
	})
	require.NoError(t, err)
	require.Contains(t, string(b), `"type":null`)

	var uploaded blob.UploadBlob
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"b","type":null,"size":0}`), &uploaded))
	require.True(t, uploaded.Type.IsNull())
	out, err := jsonv2.Marshal(&uploaded)
	require.NoError(t, err)
	require.Contains(t, string(out), `"type":null`)
	require.Contains(t, string(out), `"size":0`)
}

func TestUploadBlobZeroSizeOnWire(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(blob.UploadBlob{ID: "b"})
	require.NoError(t, err)
	require.Contains(t, string(b), `"size":0`)
}

func TestUploadResponseUnmarshal(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["Blob/upload",{"accountId":"u1","created":{"c1":{"id":"blob1","size":5,"type":"text/plain"}},"notCreated":{"c2":{"type":"invalidProperties"}}},"0"]]}`)

	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)

	inv := resp.Responses[0]
	assert.Equal(t, "Blob/upload", inv.Name)
	assert.Equal(t, "0", inv.CallID)

	methodResp, ok := inv.Args.(*blob.UploadResponse)
	require.True(t, ok)
	require.Contains(t, methodResp.Created, jmap.ID("c1"))
	assert.Equal(t, jmap.ID("blob1"), methodResp.Created["c1"].ID)
	assert.Equal(t, jmap.UnsignedInt(5), methodResp.Created["c1"].Size)
	typ, okType := methodResp.Created["c1"].Type.Value()
	require.True(t, okType)
	assert.Equal(t, "text/plain", typ)
	require.Contains(t, methodResp.NotCreated, jmap.ID("c2"))
	assert.Equal(t, "invalidProperties", methodResp.NotCreated["c2"].Type)
}
