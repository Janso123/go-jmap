package mailbox_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/mailbox"
	"github.com/stretchr/testify/require"
)

func TestQueryUnmarshalFilter(t *testing.T) {
	var q mailbox.Query
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","filter":{"operator":"OR","conditions":[{"name":"Inbox"}]},"sort":[{"property":"name"}],"sortAsTree":true}`), &q))
	op := q.Filter.(*jmap.FilterOperator)
	require.Equal(t, "Inbox", op.Conditions[0].(*mailbox.FilterCondition).Name)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Len(t, q.Sort, 1)
	require.True(t, q.SortAsTree)
}

func TestQueryChangesUnmarshalFilter(t *testing.T) {
	var q mailbox.QueryChanges
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","sinceQueryState":"s1","filter":{"name":"Inbox"},"sort":[{"property":"name"}]}`), &q))
	require.Equal(t, "Inbox", q.Filter.(*mailbox.FilterCondition).Name)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Equal(t, "s1", q.SinceQueryState)
	require.Len(t, q.Sort, 1)
}
