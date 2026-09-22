package emailsubmission_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/emailsubmission"
	"github.com/stretchr/testify/require"
)

func TestQueryUnmarshalFilter(t *testing.T) {
	var q emailsubmission.Query
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","filter":{"operator":"OR","conditions":[{"undoStatus":"pending"}]},"sort":[{"property":"sendAt"}]}`), &q))
	op := q.Filter.(*jmap.FilterOperator)
	require.Equal(t, emailsubmission.UndoPending, op.Conditions[0].(*emailsubmission.FilterCondition).UndoStatus)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Len(t, q.Sort, 1)
}

func TestQueryChangesUnmarshalFilter(t *testing.T) {
	var q emailsubmission.QueryChanges
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","sinceQueryState":"s1","filter":{"undoStatus":"pending"},"sort":[{"property":"sendAt"}]}`), &q))
	require.Equal(t, emailsubmission.UndoPending, q.Filter.(*emailsubmission.FilterCondition).UndoStatus)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Equal(t, "s1", q.SinceQueryState)
	require.Len(t, q.Sort, 1)
}
