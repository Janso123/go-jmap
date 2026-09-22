package quota

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangesInvoke(t *testing.T) {
	m := &Changes{
		Account:    "u1",
		SinceState: "1234",
	}
	req := &jmap.Request{}

	id := req.Invoke(m)
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	assert.NoError(t, err)
	expected := `{"using":["urn:ietf:params:jmap:quota"],"methodCalls":[["Quota/changes",{"accountId":"u1","sinceState":"1234"},"0"]]}`
	assert.Equal(t, expected, string(data))

	t.Run("manual", func(t *testing.T) {
		req = &jmap.Request{
			Using: []jmap.URI{URI},
			Calls: []*jmap.Invocation{
				{
					Name:   "Quota/changes",
					Args:   m,
					CallID: "manual",
				},
			},
		}

		data, err := jsonv2.Marshal(req)
		assert.NoError(t, err)
		expected := `{"using":["urn:ietf:params:jmap:quota"],"methodCalls":[["Quota/changes",{"accountId":"u1","sinceState":"1234"},"manual"]]}`
		assert.Equal(t, expected, string(data))
	})
}

func TestChangesUpdatedPropertiesNull(t *testing.T) {
	t.Parallel()
	var resp ChangesResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"updatedProperties":null}`), &resp))
	b, err := jsonv2.Marshal(&resp)
	require.NoError(t, err)
	require.Contains(t, string(b), `"updatedProperties":null`)
}
