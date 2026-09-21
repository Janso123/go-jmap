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
	require.Equal(t, jmap.ID("bTo"), resp.Copied["bFrom"])
}

func TestUploadObjectEmptyDataEmitted(t *testing.T) {
	b, err := jsonv2.Marshal(blob.UploadObject{Data: []blob.DataSource{}})
	require.NoError(t, err)
	require.Contains(t, string(b), `"data":[]`)
}
