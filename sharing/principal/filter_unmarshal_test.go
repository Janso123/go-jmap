package principal_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/sharing/principal"
	"github.com/stretchr/testify/require"
)

func TestQueryUnmarshalFilter(t *testing.T) {
	var q principal.Query
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","filter":{"operator":"OR","conditions":[{"type":"individual"}]},"sort":[{"property":"name"}]}`), &q))
	op := q.Filter.(*jmap.FilterOperator)
	require.Equal(t, principal.PrincipalType("individual"), op.Conditions[0].(*principal.FilterCondition).Type)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Len(t, q.Sort, 1)
}

func TestQueryChangesUnmarshalFilter(t *testing.T) {
	var q principal.QueryChanges
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","sinceQueryState":"s1","filter":{"type":"individual"},"sort":[{"property":"name"}]}`), &q))
	require.Equal(t, principal.PrincipalType("individual"), q.Filter.(*principal.FilterCondition).Type)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Equal(t, "s1", q.SinceQueryState)
	require.Len(t, q.Sort, 1)
}
