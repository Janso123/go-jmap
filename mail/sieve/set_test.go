package sieve

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
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

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:sieve"],"methodCalls":[["SieveScript/set",{"accountId":"u1","onSuccessActivateScript":"#create-1"},"0"]]}`,
		string(data))
}
