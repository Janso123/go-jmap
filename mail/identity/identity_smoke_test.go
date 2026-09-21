package identity_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/emailsubmission"
	"github.com/Janso123/go-jmap/mail/identity"
	"github.com/stretchr/testify/require"
)

func TestIdentityMethodNames(t *testing.T) {
	require.Equal(t, "Identity/get", (&identity.Get{}).Name())
	require.Equal(t, "Identity/changes", (&identity.Changes{}).Name())
	require.Equal(t, "Identity/set", (&identity.Set{}).Name())
}

func TestIdentityGetRequiresEmailSubmission(t *testing.T) {
	require.Equal(t, []jmap.URI{emailsubmission.URI}, (&identity.Get{}).Requires())
}

func TestIdentityGetResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["Identity/get",{"accountId":"u1","state":"s0","list":[],"notFound":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*identity.GetResponse)
	require.True(t, ok)
}
