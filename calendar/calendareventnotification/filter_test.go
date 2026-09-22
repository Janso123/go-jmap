package calendareventnotification

import (
	jsonv2 "encoding/json/v2"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterConditionMarshal(t *testing.T) {
	after := jmap.UTCDate(time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC))
	before := jmap.UTCDate(time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC))

	data, err := jsonv2.Marshal(&FilterCondition{
		After:            jmap.Some(after),
		Before:           jmap.Some(before),
		Type:             TypeUpdated,
		CalendarEventIDs: jmap.Some([]jmap.ID{"ev1", "ev2"}),
	})
	require.NoError(t, err)
	assert.Equal(t,
		`{"after":"2026-09-01T00:00:00Z","before":"2026-10-01T00:00:00Z","type":"updated","calendarEventIds":["ev1","ev2"]}`,
		string(data))
}

func TestFilterOperatorMarshal(t *testing.T) {
	data, err := jsonv2.Marshal(And(
		&FilterCondition{Type: TypeCreated},
		&FilterCondition{CalendarEventIDs: jmap.Some([]jmap.ID{"ev1"})},
	))
	require.NoError(t, err)
	assert.Equal(t,
		`{"operator":"AND","conditions":[{"type":"created"},{"calendarEventIds":["ev1"]}]}`,
		string(data))
}
