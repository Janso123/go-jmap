package jmap_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestOptionalNullAndSome(t *testing.T) {
	t.Parallel()
	type wrap struct {
		ID jmap.Optional[jmap.ID] `json:"id,omitzero"`
	}
	b, err := jsonv2.Marshal(wrap{})
	require.NoError(t, err)
	require.JSONEq(t, `{}`, string(b))

	b, err = jsonv2.Marshal(wrap{ID: jmap.Null[jmap.ID]()})
	require.NoError(t, err)
	require.JSONEq(t, `{"id":null}`, string(b))

	b, err = jsonv2.Marshal(wrap{ID: jmap.Some(jmap.ID("x"))})
	require.NoError(t, err)
	require.JSONEq(t, `{"id":"x"}`, string(b))

	var got wrap
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":null}`), &got))
	require.True(t, got.ID.IsNull())
}
