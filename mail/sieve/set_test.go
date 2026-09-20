package sieve

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetInvokeUsesOnSuccessActivateScript(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Set{
		Account:                 "u1",
		OnSuccessActivateScript: "#create-1",
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:sieve"],"methodCalls":[["SieveScript/set",{"accountId":"u1","onSuccessActivateScript":"#create-1"},"0"]]}`,
		string(data))
}
