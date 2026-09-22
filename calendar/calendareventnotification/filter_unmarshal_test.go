package calendareventnotification_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/calendareventnotification"
	"github.com/stretchr/testify/require"
)

func TestQueryUnmarshalFilter(t *testing.T) {
	var q calendareventnotification.Query
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","filter":{"operator":"OR","conditions":[{"type":"created"}]},"sort":[{"property":"created"}]}`), &q))
	op := q.Filter.(*jmap.FilterOperator)
	require.Equal(t, calendareventnotification.Type("created"), op.Conditions[0].(*calendareventnotification.FilterCondition).Type)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Len(t, q.Sort, 1)
}

func TestQueryChangesUnmarshalFilter(t *testing.T) {
	var q calendareventnotification.QueryChanges
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","sinceQueryState":"s1","filter":{"type":"created"},"sort":[{"property":"created"}]}`), &q))
	require.Equal(t, calendareventnotification.Type("created"), q.Filter.(*calendareventnotification.FilterCondition).Type)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Equal(t, "s1", q.SinceQueryState)
	require.Len(t, q.Sort, 1)
}
