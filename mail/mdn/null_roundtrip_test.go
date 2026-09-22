package mdn_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/mdn"
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

func TestMDNNullSubjectAndFalseIncludeOriginalMessage(t *testing.T) {
	t.Parallel()
	m := mdn.MDN{
		Subject:                jmap.Null[string](),
		IncludeOriginalMessage: new(false),
	}
	b, err := jsonv2.Marshal(&m)
	require.NoError(t, err)
	out := string(b)
	require.Contains(t, out, `"subject":null`)
	require.Contains(t, out, `"includeOriginalMessage":false`)

	rt := roundTrip[mdn.MDN](t, `{"subject":null,"includeOriginalMessage":false}`)
	require.Contains(t, rt, `"subject":null`)
	require.Contains(t, rt, `"includeOriginalMessage":false`)
}
