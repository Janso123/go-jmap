package searchsnippet_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/Janso123/go-jmap/mail/searchsnippet"
	"github.com/stretchr/testify/require"
)

func TestGetUnmarshalFilter(t *testing.T) {
	var g searchsnippet.Get
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","filter":{"operator":"OR","conditions":[{"subject":"hello"}]},"emailIds":["e1"]}`), &g))
	op := g.Filter.(*jmap.FilterOperator)
	require.Equal(t, "hello", op.Conditions[0].(*email.FilterCondition).Subject)
	require.Equal(t, jmap.ID("a"), g.Account)
	require.Equal(t, []jmap.ID{"e1"}, g.EmailIDs)
}
