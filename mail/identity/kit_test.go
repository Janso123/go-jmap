package identity_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap/mail/identity"
	"github.com/stretchr/testify/require"
)

func TestIdentityGetEmbedsKit(t *testing.T) {
	m := &identity.Get{Account: "a1"}
	require.Equal(t, "Identity/get", m.Name())
	b, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"a1"}`, string(b))

	var resp identity.GetResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a1","state":"s0","list":[{"id":"id1","email":"me@example.com"}],"notFound":[]}`), &resp))
	require.Equal(t, []identity.Identity{{ID: "id1", Email: "me@example.com"}}, resp.List)
}
