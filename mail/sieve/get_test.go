package sieve

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Get{
		Account:    "u1",
		IDs:        []jmap.ID{"script1"},
		Properties: []string{"name", "blobId", "isActive"},
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:sieve"],"methodCalls":[["SieveScript/get",{"accountId":"u1","ids":["script1"],"properties":["name","blobId","isActive"]},"0"]]}`,
		string(data))
}
