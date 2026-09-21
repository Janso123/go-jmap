package sieve_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/sieve"
	"github.com/stretchr/testify/require"
)

func TestSieveScriptIsObject(t *testing.T) {
	var _ jmap.Object = sieve.SieveScript{}
	require.Equal(t, "SieveScript", sieve.SieveScript{}.JMAPType())
	require.Equal(t, []jmap.URI{sieve.URI}, sieve.SieveScript{}.Requires())
}

func TestSieveScriptGetName(t *testing.T) {
	require.Equal(t, "SieveScript/get", (&sieve.Get{}).Name())
}

func TestSieveScriptQueryName(t *testing.T) {
	require.Equal(t, "SieveScript/query", (&sieve.Query{}).Name())
}

func TestSieveScriptSetName(t *testing.T) {
	require.Equal(t, "SieveScript/set", (&sieve.Set{}).Name())
}

func TestSieveScriptValidateRemainsManual(t *testing.T) {
	require.Equal(t, "SieveScript/validate", (&sieve.Validate{}).Name())
	require.Equal(t, []jmap.URI{sieve.URI}, (&sieve.Validate{}).Requires())
}
