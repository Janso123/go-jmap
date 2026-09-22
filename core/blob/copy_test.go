package blob_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/blob"
	"github.com/stretchr/testify/require"
)

func TestCopyResponseCopiedKey(t *testing.T) {
	var resp blob.CopyResponse
	err := jsonv2.Unmarshal([]byte(`{"fromAccountId":"A","accountId":"B","copied":{"bFrom":"bTo"}}`), &resp)
	require.NoError(t, err)
	copied, ok := resp.Copied.Value()
	require.True(t, ok)
	require.Equal(t, jmap.ID("bTo"), copied["bFrom"])
}

func TestCopyResponseNullCopiedAndNotCopied(t *testing.T) {
	t.Parallel()
	var resp blob.CopyResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"copied":null,"notCopied":null}`), &resp))
	b, err := jsonv2.Marshal(&resp)
	require.NoError(t, err)
	require.Contains(t, string(b), `"copied":null`)
	require.Contains(t, string(b), `"notCopied":null`)
}

func TestUploadObjectEmptyDataEmitted(t *testing.T) {
	b, err := jsonv2.Marshal(blob.UploadObject{Data: []blob.DataSource{}})
	require.NoError(t, err)
	require.Contains(t, string(b), `"data":[]`)
}
