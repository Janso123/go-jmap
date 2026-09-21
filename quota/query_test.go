package quota

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
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
		Sort: []*jmap.Comparator{
			{Property: "name"},
		},
		Limit: jmap.Uint64Ptr(10),
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:quota"],"methodCalls":[["Quota/query",{"accountId":"u1","limit":10,"filter":{"type":"Email"},"sort":[{"property":"name","isAscending":false}]},"0"]]}`,
		string(data))
}
