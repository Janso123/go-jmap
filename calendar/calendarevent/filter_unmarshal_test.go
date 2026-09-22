package calendarevent_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar/calendarevent"
	"github.com/stretchr/testify/require"
)

func TestQueryUnmarshalFilter(t *testing.T) {
	var q calendarevent.Query
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","filter":{"operator":"OR","conditions":[{"inCalendar":"c1"}]},"sort":[{"property":"start"}],"expandRecurrences":true}`), &q))
	op := q.Filter.(*jmap.FilterOperator)
	require.Equal(t, jmap.ID("c1"), op.Conditions[0].(*calendarevent.FilterCondition).InCalendar)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Len(t, q.Sort, 1)
	require.True(t, q.ExpandRecurrences)
}

func TestQueryChangesUnmarshalFilter(t *testing.T) {
	var q calendarevent.QueryChanges
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","sinceQueryState":"s1","filter":{"inCalendar":"c1"},"sort":[{"property":"start"}]}`), &q))
	require.Equal(t, jmap.ID("c1"), q.Filter.(*calendarevent.FilterCondition).InCalendar)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Equal(t, "s1", q.SinceQueryState)
	require.Len(t, q.Sort, 1)
}
