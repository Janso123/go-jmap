package email_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func TestEmailSortComparatorsWire(t *testing.T) {
	t.Parallel()

	b, err := jsonv2.Marshal(email.Desc(email.SortReceivedAt))
	require.NoError(t, err)
	require.JSONEq(t, `{"property":"receivedAt","isAscending":false}`, string(b))

	b, err = jsonv2.Marshal(email.ByKeyword(email.KeywordSeen, true))
	require.NoError(t, err)
	require.JSONEq(t, `{"property":"hasKeyword","keyword":"$seen","isAscending":true}`, string(b))

	q := &email.Query{Sort: []*email.SortComparator{email.ByKeyword(email.KeywordFlagged, false)}}
	b, err = jsonv2.Marshal(q)
	require.NoError(t, err)
	require.Contains(t, string(b), `"sort"`)
	require.Contains(t, string(b), `"hasKeyword"`)
	require.Contains(t, string(b), `"$flagged"`)
	require.Contains(t, string(b), `"isAscending":false`)
}
