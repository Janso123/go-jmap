package identity_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentityJSONRoundTrip(t *testing.T) {
	id := identity.Identity{
		ID:    "id1",
		Name:  "Work",
		Email: "me@example.com",
		ReplyTo: []*mail.Address{
			{Name: "Me", Email: "me@example.com"},
		},
		Bcc: []*mail.Address{
			{Email: "archive@example.com"},
		},
		TextSignature: "Regards",
		HTMLSignature: "<p>Regards</p>",
		MayDelete:     true,
	}

	data, err := json.Marshal(id)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"id": "id1",
		"name": "Work",
		"email": "me@example.com",
		"replyTo": [{"name": "Me", "email": "me@example.com"}],
		"bcc": [{"email": "archive@example.com"}],
		"textSignature": "Regards",
		"htmlSignature": "<p>Regards</p>",
		"mayDelete": true
	}`, string(data))

	var got identity.Identity
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, id.ID, got.ID)
	require.Equal(t, id.Name, got.Name)
	require.Equal(t, id.Email, got.Email)
	require.Equal(t, id.ReplyTo[0].Email, got.ReplyTo[0].Email)
	require.Equal(t, id.Bcc[0].Email, got.Bcc[0].Email)
	require.Equal(t, id.TextSignature, got.TextSignature)
	require.Equal(t, id.HTMLSignature, got.HTMLSignature)
	require.Equal(t, id.MayDelete, got.MayDelete)
}
