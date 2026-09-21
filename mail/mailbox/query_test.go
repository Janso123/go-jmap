package mailbox

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	query := &Query{
		Account: "xyz",
		Filter: &FilterCondition{
			Name: "Inbox",
		},
		Sort: []*jmap.Comparator{
			{
				Property: "name",
			},
		},
		Limit: jmap.Uint64Ptr(10),
	}
	data, err := jsonv2.Marshal(query)
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"xyz","limit":10,"filter":{"name":"Inbox"},"sort":[{"property":"name","isAscending":false}]}`, string(data))
}
