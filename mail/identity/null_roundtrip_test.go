package identity_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap/mail/identity"
	"github.com/stretchr/testify/require"
)

func roundTrip[T any](t *testing.T, in string) string {
	t.Helper()
	var v T
	require.NoError(t, jsonv2.Unmarshal([]byte(in), &v))
	b, err := jsonv2.Marshal(&v)
	require.NoError(t, err)
	return string(b)
}

func TestIdentityNullReplyToAndBccOmitEmptySignature(t *testing.T) {
	t.Parallel()
	out := roundTrip[identity.Identity](t, `{"replyTo":null,"bcc":null,"textSignature":""}`)
	require.Contains(t, out, `"replyTo":null`)
	require.Contains(t, out, `"bcc":null`)
	require.NotContains(t, out, "textSignature")
}
