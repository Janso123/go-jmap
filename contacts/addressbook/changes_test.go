package addressbook

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangesInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Changes{
		Account:    "u1",
		SinceState: "s1",
		MaxChanges: jmap.Uint64Ptr(50),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["AddressBook/changes",{"accountId":"u1","sinceState":"s1","maxChanges":50},"0"]]}`,
		string(data))
}

func TestChangesRequiresContactsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{contacts.URI}, (&Changes{}).Requires())
}
