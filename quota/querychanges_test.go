package quota

import (
	jsonv2 "encoding/json/v2"
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
		Sort:            []*jmap.Comparator{{Property: "used"}},
		SinceQueryState: "s1",
		MaxChanges:      jmap.Uint64Ptr(50),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:quota"],"methodCalls":[["Quota/queryChanges",{"accountId":"u1","sinceQueryState":"s1","maxChanges":50,"filter":{"resourceType":"octets"},"sort":[{"property":"used"}]},"0"]]}`,
		string(data))
}
