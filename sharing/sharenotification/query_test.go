package sharenotification

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryInvoke(t *testing.T) {
	req := &jmap.Request{}
	after := jmap.UTCDate(time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC))

	id := req.Invoke(&Query{
		Account: "u1",
		Filter: &FilterCondition{
			After:      &after,
			ObjectType: "Mailbox",
		},
		Sort: []*jmap.Comparator{
			{Property: "created"},
		},
		Limit: jmap.Uint64Ptr(10),
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:principals"],"methodCalls":[["ShareNotification/query",{"accountId":"u1","limit":10,"filter":{"after":"2026-09-01T00:00:00Z","objectType":"Mailbox"},"sort":[{"property":"created","isAscending":false}]},"0"]]}`,
		string(data))
}
