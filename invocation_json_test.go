package jmap_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestInvocationRoundTripV2(t *testing.T) {
	inv := &jmap.Invocation{
		Name:   "Core/echo",
		Args:   map[string]any{"Hello": "world"},
		CallID: "c0",
	}
	b, err := jsonv2.Marshal(inv)
	require.NoError(t, err)

	var out jmap.Invocation
	require.NoError(t, jsonv2.Unmarshal(b, &out))
	require.Equal(t, "Core/echo", out.Name)
	require.Equal(t, "c0", out.CallID)
}
