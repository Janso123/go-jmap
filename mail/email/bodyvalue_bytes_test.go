package email_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func TestMaxBodyValueBytesExplicitZeroOnWire(t *testing.T) {
	t.Parallel()
	zero := jmap.Uint64Ptr(0)

	b, err := jsonv2.Marshal(&email.Get{MaxBodyValueBytes: zero})
	require.NoError(t, err)
	require.Contains(t, string(b), `"maxBodyValueBytes":0`)

	b, err = jsonv2.Marshal(&email.Parse{MaxBodyValueBytes: zero})
	require.NoError(t, err)
	require.Contains(t, string(b), `"maxBodyValueBytes":0`)
}

func TestMaxBodyValueBytesUnsetOmitted(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(&email.Get{})
	require.NoError(t, err)
	require.NotContains(t, string(b), "maxBodyValueBytes")

	b, err = jsonv2.Marshal(&email.Parse{})
	require.NoError(t, err)
	require.NotContains(t, string(b), "maxBodyValueBytes")
}
