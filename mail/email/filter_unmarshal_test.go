package email_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func TestQueryUnmarshalFilter(t *testing.T) {
	var q email.Query
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","filter":{"operator":"OR","conditions":[{"subject":"hello"}]},"sort":[{"property":"receivedAt"}],"collapseThreads":true}`), &q))
	op := q.Filter.(*jmap.FilterOperator)
	require.Equal(t, "hello", op.Conditions[0].(*email.FilterCondition).Subject)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Len(t, q.Sort, 1)
	require.True(t, q.CollapseThreads)
}

func TestQueryChangesUnmarshalFilter(t *testing.T) {
	var q email.QueryChanges
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","sinceQueryState":"s1","filter":{"subject":"hello"},"sort":[{"property":"receivedAt"}],"collapseThreads":true}`), &q))
	require.Equal(t, "hello", q.Filter.(*email.FilterCondition).Subject)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Equal(t, "s1", q.SinceQueryState)
	require.Len(t, q.Sort, 1)
	require.True(t, q.CollapseThreads)
}
