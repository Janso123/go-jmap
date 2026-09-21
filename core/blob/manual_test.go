package blob_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/blob"
	"github.com/stretchr/testify/require"
)

// Blob management methods (RFC 9404 + core Blob/copy) are non-standard relative
// to the generic Object kit; they remain manual in go-jmap v1.
func TestBlobMethodsRemainManual(t *testing.T) {
	require.Equal(t, "Blob/get", (&blob.Get{}).Name())
	require.Equal(t, "Blob/upload", (&blob.Upload{}).Name())
	require.Equal(t, "Blob/lookup", (&blob.Lookup{}).Name())
	require.Equal(t, "Blob/copy", (&blob.Copy{}).Name())

	require.Equal(t, []jmap.URI{blob.URI}, (&blob.Get{}).Requires())
	require.Equal(t, []jmap.URI{blob.URI}, (&blob.Upload{}).Requires())
	require.Equal(t, []jmap.URI{blob.URI}, (&blob.Lookup{}).Requires())
	require.Nil(t, (&blob.Copy{}).Requires())
}

func TestBlobCapabilityURIUnchanged(t *testing.T) {
	require.Equal(t, jmap.URI("urn:ietf:params:jmap:blob"), blob.URI)
}
