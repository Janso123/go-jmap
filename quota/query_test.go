package quota

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Query{
		Account: "u1",
		Filter: &FilterCondition{
			Type: "Email",
		},
		Sort: []*SortComparator{
			{Property: "name"},
		},
		Limit: 10,
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:quota"],"methodCalls":[["Quota/query",{"accountId":"u1","filter":{"type":"Email"},"sort":[{"property":"name","isAscending":false}],"limit":10},"0"]]}`,
		string(data))
}
