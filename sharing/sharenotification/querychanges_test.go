package sharenotification

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryChangesInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&QueryChanges{
		Account: "u1",
		Filter: &FilterCondition{
			ObjectAccountID: "a1",
		},
		Sort:            []*jmap.Comparator{{Property: "created"}},
		SinceQueryState: "s1",
		MaxChanges:      50,
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:principals"],"methodCalls":[["ShareNotification/queryChanges",{"accountId":"u1","sinceQueryState":"s1","maxChanges":50,"filter":{"objectAccountId":"a1"},"sort":[{"property":"created","isAscending":false}]},"0"]]}`,
		string(data))
}
