package sieve

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Validate{
		Account: "u1",
		BlobID:  "blob1",
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:sieve"],"methodCalls":[["SieveScript/validate",{"accountId":"u1","blobId":"blob1"},"0"]]}`,
		string(data))
}

func TestValidateResponseUnmarshal(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["SieveScript/validate",{"accountId":"u1","error":{"type":"invalidSieve","description":"line 3: syntax error"}},"0"]]}`)

	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)

	inv := resp.Responses[0]
	assert.Equal(t, "SieveScript/validate", inv.Name)
	assert.Equal(t, "0", inv.CallID)

	methodResp, ok := inv.Args.(*ValidateResponse)
	require.True(t, ok)
	require.NotNil(t, methodResp.Error)
	assert.Equal(t, "u1", string(methodResp.Account))
	assert.Equal(t, "invalidSieve", methodResp.Error.Type)
	require.NotNil(t, methodResp.Error.Description)
	assert.Equal(t, "line 3: syntax error", *methodResp.Error.Description)
}
