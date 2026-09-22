package sieve_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/sieve"
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

func TestSieveCapabilityNullRedirectsAndExternalLists(t *testing.T) {
	t.Parallel()
	out := roundTrip[sieve.Capability](t, `{"maxNumberRedirects":null,"externalLists":null}`)
	require.Contains(t, out, `"maxNumberRedirects":null`)
	require.Contains(t, out, `"externalLists":null`)
}

func TestSieveCreateOmitsIsActive(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(&sieve.Set{
		Set: jmap.Set[sieve.SieveScript]{
			Create: jmap.Some(map[jmap.ID]*sieve.SieveScript{
				"c1": {Name: jmap.Some("s"), BlobID: "b"},
			}),
		},
	})
	require.NoError(t, err)
	require.NotContains(t, string(b), "isActive")
}
