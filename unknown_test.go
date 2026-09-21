package jmap_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestUnknownMethodSurvivesResponse(t *testing.T) {
	raw := `["Foo/bar",{"x":1},"c9"]`
	var inv jmap.Invocation
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &inv))
	u, ok := inv.Args.(*jmap.UnknownResponse)
	require.True(t, ok)
	require.Equal(t, "Foo/bar", u.Name)
	require.Contains(t, string(u.Raw), `"x":1`)
}

func TestUnknownResponseRemarshalInvocation(t *testing.T) {
	raw := `["Foo/bar",{"x":1},"c9"]`
	var inv jmap.Invocation
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &inv))

	out, err := jsonv2.Marshal(&inv)
	require.NoError(t, err)
	s := string(out)
	require.Contains(t, s, `"x":1`)
	require.NotContains(t, s, `"Name"`)
}
