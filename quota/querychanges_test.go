package quota

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
		Account:         "u1",
		Filter:          &FilterCondition{ResourceType: "octets"},
		Sort:            []*SortComparator{{Property: "used"}},
		SinceQueryState: "s1",
		MaxChanges:      50,
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:quota"],"methodCalls":[["Quota/queryChanges",{"accountId":"u1","sinceQueryState":"s1","maxChanges":50,"filter":{"resourceType":"octets"},"sort":[{"property":"used","isAscending":false}]},"0"]]}`,
		string(data))
}
