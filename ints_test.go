package jmap

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/stretchr/testify/require"
)

func TestUnsignedIntRange(t *testing.T) {
	_, err := jsonv2.Marshal(UnsignedInt(1 << 53))
	require.Error(t, err)
	var u UnsignedInt
	require.Error(t, jsonv2.Unmarshal([]byte(`9007199254740992`), &u))
	require.Error(t, jsonv2.Unmarshal([]byte(`1.5`), &u))
	require.Error(t, jsonv2.Unmarshal([]byte(`null`), &u))
	require.NoError(t, jsonv2.Unmarshal([]byte(`0`), &u))
	b, err := jsonv2.Marshal(UnsignedInt(0))
	require.NoError(t, err)
	require.Equal(t, `0`, string(b))
}

func TestIntRange(t *testing.T) {
	_, err := jsonv2.Marshal(Int(-(1 << 53)))
	require.Error(t, err)
	var i Int
	require.NoError(t, jsonv2.Unmarshal([]byte(`-9007199254740991`), &i))
}
