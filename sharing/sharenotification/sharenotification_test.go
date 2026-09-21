package sharenotification

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShareNotificationMarshal(t *testing.T) {
	created := time.Date(2026, time.September, 19, 16, 0, 0, 0, time.UTC)
	email := "jane@example.com"
	principalID := jmap.ID("p1")

	data, err := json.Marshal(&ShareNotification{
		ID:              "sn1",
		Created:         &created,
		ChangedBy:       &Entity{Name: "Jane Doe", Email: &email, PrincipalID: &principalID},
		ObjectAccountID: "a1",
		ObjectType:      "Mailbox",
		ObjectID:        "m1",
		OldRights: map[string]bool{
			"mayReadItems": true,
			"maySubmit":    false,
		},
		NewRights: map[string]bool{
			"mayReadItems": true,
			"maySubmit":    true,
		},
		Name: "Inbox",
	})
	require.NoError(t, err)
	assert.Equal(t,
		`{"id":"sn1","created":"2026-09-19T16:00:00Z","changedBy":{"name":"Jane Doe","email":"jane@example.com","principalId":"p1"},"objectAccountId":"a1","objectType":"Mailbox","objectId":"m1","oldRights":{"mayReadItems":true,"maySubmit":false},"newRights":{"mayReadItems":true,"maySubmit":true},"name":"Inbox"}`,
		string(data))
}
