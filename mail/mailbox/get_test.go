package mailbox

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	m := &Get{
		Account: "account-id",
	}
	req := &jmap.Request{}

	id := req.Invoke(m)
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	assert.NoError(t, err)
	expected := `{"using":["urn:ietf:params:jmap:mail"],"methodCalls":[["Mailbox/get",{"accountId":"account-id"},"0"]]}`
	assert.Equal(t, expected, string(data))

	t.Run("manual", func(t *testing.T) {
		m := &Get{
			Account: "account-id",
		}
		req = &jmap.Request{
			Using: []jmap.URI{mail.URI},
			Calls: []*jmap.Invocation{
				{
					Name:   "Mailbox/get",
					Args:   m,
					CallID: "manual",
				},
			},
		}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		expected := `{"using":["urn:ietf:params:jmap:mail"],"methodCalls":[["Mailbox/get",{"accountId":"account-id"},"manual"]]}`
		assert.Equal(t, expected, string(data))
	})
}
