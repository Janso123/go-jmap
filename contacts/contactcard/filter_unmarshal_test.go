package contactcard_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts/contactcard"
	"github.com/stretchr/testify/require"
)

func TestQueryUnmarshalFilter(t *testing.T) {
	var q contactcard.Query
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","filter":{"operator":"OR","conditions":[{"name":"Ada"}]},"sort":[{"property":"name"}]}`), &q))
	op := q.Filter.(*jmap.FilterOperator)
	require.Equal(t, "Ada", op.Conditions[0].(*contactcard.FilterCondition).Name)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Len(t, q.Sort, 1)
}

func TestQueryChangesUnmarshalFilter(t *testing.T) {
	var q contactcard.QueryChanges
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","sinceQueryState":"s1","filter":{"name":"Ada"},"sort":[{"property":"name"}]}`), &q))
	require.Equal(t, "Ada", q.Filter.(*contactcard.FilterCondition).Name)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Equal(t, "s1", q.SinceQueryState)
	require.Len(t, q.Sort, 1)
}
