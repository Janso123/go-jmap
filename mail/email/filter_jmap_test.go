package email_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func TestJmapAndInMailboxWire(t *testing.T) {
	t.Parallel()
	f := jmap.And(email.InMailbox("m1"))
	b, err := jsonv2.Marshal(f)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"operator":"AND",
		"conditions":[{"inMailbox":"m1"}]
	}`, string(b))
}

func TestEmailQueryDescSortWire(t *testing.T) {
	t.Parallel()
	q := &email.Query{Sort: []*jmap.Comparator{jmap.Desc(email.SortReceivedAt)}}
	b, err := jsonv2.Marshal(q)
	require.NoError(t, err)
	require.Contains(t, string(b), `"isAscending":false`)
	require.Contains(t, string(b), `"receivedAt"`)
}
