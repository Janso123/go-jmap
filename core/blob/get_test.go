package blob_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInvoke(t *testing.T) {
	req := &jmap.Request{}
	offset := jmap.UnsignedInt(4)
	length := jmap.UnsignedInt(9)

	id := req.Invoke(&blob.Get{
		Account:    "u1",
		IDs:        jmap.Some([]jmap.ID{"blob1"}),
		Properties: []string{"data:asText", "digest:sha-256", "size"},
		Offset:     &offset,
		Length:     &length,
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:blob"],"methodCalls":[["Blob/get",{"accountId":"u1","ids":["blob1"],"properties":["data:asText","digest:sha-256","size"],"offset":4,"length":9},"0"]]}`,
		string(data))
}

func TestGetResponseUnmarshal(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["Blob/get",{"accountId":"u1","list":[{"id":"blob1","data:asText":"hello","digest:sha-256":"abc=","size":5}],"notFound":["missing"]},"0"]]}`)

	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)

	inv := resp.Responses[0]
	assert.Equal(t, "Blob/get", inv.Name)
	assert.Equal(t, "0", inv.CallID)

	methodResp, ok := inv.Args.(*blob.GetResponse)
	require.True(t, ok)
	require.Len(t, methodResp.List, 1)
	text, okText := methodResp.List[0].DataAsText.Value()
	require.True(t, okText)
	assert.Equal(t, "hello", text)
	assert.Equal(t, "abc=", methodResp.List[0].Digests["sha-256"])
	require.NotNil(t, methodResp.List[0].Size)
	assert.Equal(t, jmap.UnsignedInt(5), *methodResp.List[0].Size)
	assert.Equal(t, []jmap.ID{"missing"}, methodResp.NotFound)
}

func TestGetResultRoundTripPreservesEmptyDataAsBase64(t *testing.T) {
	var result blob.GetResult

	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"blob1","data:asBase64":""}`), &result))

	data, err := jsonv2.Marshal(result)
	require.NoError(t, err)
	assert.JSONEq(t, `{"id":"blob1","data:asBase64":"","isEncodingProblem":false,"isTruncated":false}`, string(data))
}

func TestGetResultDataNullRemarshals(t *testing.T) {
	t.Parallel()
	var text blob.GetResult
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"b1","data:asText":null,"size":1}`), &text))
	require.True(t, text.DataAsText.IsNull())
	b, err := jsonv2.Marshal(&text)
	require.NoError(t, err)
	require.JSONEq(t, `{"id":"b1","data:asText":null,"size":1,"isEncodingProblem":false,"isTruncated":false}`, string(b))

	var b64 blob.GetResult
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"b1","isEncodingProblem":true,"data:asBase64":null,"size":1}`), &b64))
	require.True(t, b64.DataAsBase64.IsNull())
	b, err = jsonv2.Marshal(&b64)
	require.NoError(t, err)
	require.Contains(t, string(b), `"data:asBase64":null`)
	require.Contains(t, string(b), `"isEncodingProblem":true`)
}

func TestGetResultDataAbsentOmitted(t *testing.T) {
	t.Parallel()
	var r blob.GetResult
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"b1","size":1}`), &r))
	require.True(t, r.DataAsText.IsZero())
	require.True(t, r.DataAsBase64.IsZero())
	b, err := jsonv2.Marshal(&r)
	require.NoError(t, err)
	require.NotContains(t, string(b), "data:asText")
	require.NotContains(t, string(b), "data:asBase64")
}

func TestGetResultZeroValueOmitsSize(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(blob.GetResult{ID: "blob1"})
	require.NoError(t, err)
	require.NotContains(t, string(b), "size")
	require.Contains(t, string(b), `"isEncodingProblem":false`)
	require.Contains(t, string(b), `"isTruncated":false`)
}

func TestGetResultAbsentSizeOmitted(t *testing.T) {
	t.Parallel()
	var result blob.GetResult
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"b"}`), &result))
	b, err := jsonv2.Marshal(&result)
	require.NoError(t, err)
	require.NotContains(t, string(b), "size")
}

func TestGetResultSizeZeroKept(t *testing.T) {
	t.Parallel()
	var result blob.GetResult
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"b","size":0}`), &result))
	b, err := jsonv2.Marshal(&result)
	require.NoError(t, err)
	require.Contains(t, string(b), `"size":0`)
}

func TestGetNullIDs(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(blob.Get{
		IDs: jmap.Null[[]jmap.ID](),
		ReferenceIDs: &jmap.ResultReference{
			ResultOf: "0",
			Name:     "Foo/get",
			Path:     "/ids",
		},
		ReferenceProperties: &jmap.ResultReference{
			ResultOf: "1",
			Name:     "Foo/get",
			Path:     "/list",
		},
	})
	require.NoError(t, err)
	out := string(b)
	require.Contains(t, out, `"ids":null`)
	require.Contains(t, out, `"#ids":{"resultOf":"0","name":"Foo/get","path":"/ids"}`)
	require.Contains(t, out, `"#properties":{"resultOf":"1","name":"Foo/get","path":"/list"}`)
}

func TestGetResultDigestNullSkippedAndFieldOrder(t *testing.T) {
	t.Parallel()
	var result blob.GetResult
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"isTruncated":true,"size":1,"digest:sha":null,"digest:sha-256":"abc","digest:md5":"zzz","data:asBase64":"aGk=","data:asText":"hi","id":"b"}`), &result))
	require.Equal(t, map[string]string{"md5": "zzz", "sha-256": "abc"}, result.Digests)
	require.NotContains(t, result.Digests, "sha")

	b, err := jsonv2.Marshal(&result)
	require.NoError(t, err)
	require.Equal(t,
		`{"id":"b","data:asText":"hi","data:asBase64":"aGk=","digest:md5":"zzz","digest:sha-256":"abc","size":1,"isEncodingProblem":false,"isTruncated":true}`,
		string(b))
}
