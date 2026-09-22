package jmap

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/stretchr/testify/require"
)

func TestOptionalUnsetMarshalErrors(t *testing.T) {
	_, err := jsonv2.Marshal(Optional[string]{})
	require.Error(t, err)
}

func TestOptionalNilSliceIsEmptyArray(t *testing.T) {
	b, err := jsonv2.Marshal(Some[[]string](nil))
	require.NoError(t, err)
	require.Equal(t, `[]`, string(b))
	b, err = jsonv2.Marshal(Null[[]string]())
	require.NoError(t, err)
	require.Equal(t, `null`, string(b))
}
