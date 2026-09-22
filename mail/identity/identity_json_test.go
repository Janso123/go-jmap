package identity_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
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
		ReplyTo: jmap.Some([]*mail.Address{
			{Name: jmap.Some("Me"), Email: "me@example.com"},
		}),
		Bcc: jmap.Some([]*mail.Address{
			{Email: "archive@example.com"},
		}),
		TextSignature: "Regards",
		HTMLSignature: "<p>Regards</p>",
		MayDelete:     new(true),
	}

	data, err := jsonv2.Marshal(id)
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
	require.NoError(t, jsonv2.Unmarshal(data, &got))
	require.Equal(t, id.ID, got.ID)
	require.Equal(t, id.Name, got.Name)
	require.Equal(t, id.Email, got.Email)
	replyTo, ok := id.ReplyTo.Value()
	require.True(t, ok)
	gotReplyTo, ok := got.ReplyTo.Value()
	require.True(t, ok)
	require.Equal(t, replyTo[0].Email, gotReplyTo[0].Email)
	bcc, ok := id.Bcc.Value()
	require.True(t, ok)
	gotBcc, ok := got.Bcc.Value()
	require.True(t, ok)
	require.Equal(t, bcc[0].Email, gotBcc[0].Email)
	require.Equal(t, id.TextSignature, got.TextSignature)
	require.Equal(t, id.HTMLSignature, got.HTMLSignature)
	require.Equal(t, id.MayDelete, got.MayDelete)
}

func TestIdentityCreateOmitsMayDelete(t *testing.T) {
	id := identity.Identity{Name: "Me", Email: "me@x"}
	b, err := jsonv2.Marshal(id)
	require.NoError(t, err)
	require.JSONEq(t, `{"name":"Me","email":"me@x"}`, string(b))
	require.NotContains(t, string(b), "mayDelete")
}
