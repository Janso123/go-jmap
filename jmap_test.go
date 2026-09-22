package jmap_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollationConstantsTyped(t *testing.T) {
	var c jmap.CollationAlgo = jmap.ASCIICasemap
	require.Equal(t, jmap.CollationAlgo("i;ascii-casemap"), c)
}

func TestIdLength(t *testing.T) {
	cases := []struct {
		id    string
		valid bool
	}{
		{
			id:    "",
			valid: false,
		},
		{
			// Length 256
			id:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUV",
			valid: false,
		},
	}

	for _, c := range cases {
		id := jmap.ID(c.id)

		ok, err := id.Valid()
		assert.Equal(t, ok, c.valid, "%v", err)
	}
}

func TestIDValidCharset(t *testing.T) {
	// Charset is RFC 8620 §1.2 / design: [A-Za-z0-9-_]{1,255} — no '.'
	ok, err := jmap.ID("abc-_012").Valid()
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = jmap.ID("bad!").Valid()
	require.Error(t, err)
	require.False(t, ok)

	ok, err = jmap.ID("abc-_.012").Valid()
	require.Error(t, err)
	require.False(t, ok)
}

func TestIDEmptyMarshals(t *testing.T) {
	b, err := jsonv2.Marshal(jmap.ID(""))
	require.NoError(t, err)
	require.Equal(t, `""`, string(b))
}

func TestNewCreationIDCharset(t *testing.T) {
	id := jmap.NewCreationID()
	ok, err := id.Valid()
	require.NoError(t, err)
	require.True(t, ok)
}
