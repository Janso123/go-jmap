package sieve_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/sieve"
	"github.com/stretchr/testify/require"
)

func TestQueryUnmarshalFilter(t *testing.T) {
	var q sieve.Query
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","filter":{"operator":"OR","conditions":[{"name":"vacation"}]},"sort":[{"property":"name"}]}`), &q))
	op := q.Filter.(*jmap.FilterOperator)
	require.Equal(t, "vacation", op.Conditions[0].(*sieve.FilterCondition).Name)
	require.Equal(t, jmap.ID("a"), q.Account)
	require.Len(t, q.Sort, 1)
}
