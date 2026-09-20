package quota

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"github.com/stretchr/testify/assert"
)

func TestChangesInvoke(t *testing.T) {
	m := &Changes{
		Account:    "u1",
		SinceState: "1234",
	}
	req := &jmap.Request{}

	id := req.Invoke(m)
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
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

		data, err := json.Marshal(req)
		assert.NoError(t, err)
		expected := `{"using":["urn:ietf:params:jmap:quota"],"methodCalls":[["Quota/changes",{"accountId":"u1","sinceState":"1234"},"manual"]]}`
		assert.Equal(t, expected, string(data))
	})
}
