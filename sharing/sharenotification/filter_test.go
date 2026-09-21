package sharenotification

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterConditionMarshal(t *testing.T) {
	after := jmap.UTCDate(time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC))
	before := jmap.UTCDate(time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC))

	data, err := json.Marshal(&FilterCondition{
		After:           &after,
		Before:          &before,
		ObjectType:      "Mailbox",
		ObjectAccountID: "a1",
	})
	require.NoError(t, err)
	assert.Equal(t,
		`{"after":"2026-09-01T00:00:00Z","before":"2026-10-01T00:00:00Z","objectType":"Mailbox","objectAccountId":"a1"}`,
		string(data))
}

func TestFilterOperatorMarshal(t *testing.T) {
	data, err := json.Marshal(And(
		&FilterCondition{ObjectType: "Mailbox"},
		&FilterCondition{ObjectAccountID: "a1"},
	))
	require.NoError(t, err)
	assert.Equal(t,
		`{"operator":"AND","conditions":[{"objectType":"Mailbox"},{"objectAccountId":"a1"}]}`,
		string(data))
}
