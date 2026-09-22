package quota_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/quota"
	"github.com/stretchr/testify/require"
)

func TestQueryUnmarshalFilter(t *testing.T) {
	var q quota.Query
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","filter":{"operator":"OR","conditions":[{"name":"mailbox"}]},"sort":[{"property":"name"}]}`), &q))
	op := q.Filter.(*jmap.FilterOperator)
	require.Equal(t, "mailbox", op.Conditions[0].(*quota.FilterCondition).Name)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Len(t, q.Sort, 1)
}

func TestQueryChangesUnmarshalFilter(t *testing.T) {
	var q quota.QueryChanges
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","sinceQueryState":"s1","filter":{"name":"mailbox"},"sort":[{"property":"name"}]}`), &q))
	require.Equal(t, "mailbox", q.Filter.(*quota.FilterCondition).Name)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Equal(t, "s1", q.SinceQueryState)
	require.Len(t, q.Sort, 1)
}
