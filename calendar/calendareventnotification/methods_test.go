package calendareventnotification

import (
	jsonv2 "encoding/json/v2"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangesInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Changes{
		Account:    "u1",
		SinceState: "s1",
		MaxChanges: jmap.Uint64Ptr(25),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEventNotification/changes",{"accountId":"u1","sinceState":"s1","maxChanges":25},"0"]]}`,
		string(data))
}

func TestChangesRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Changes{}).Requires())
}

func TestSetInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Set{
		Account: "u1",
		Destroy: []jmap.ID{"cn1"},
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEventNotification/set",{"accountId":"u1","destroy":["cn1"]},"0"]]}`,
		string(data))
}

func TestSetRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Set{}).Requires())
}

func TestSetDestroyOnlyOmitsCreateUpdate(t *testing.T) {
	t.Parallel()
	m := &Set{
		Account: "u1",
		Destroy: []jmap.ID{"cn1"},
	}
	b, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	require.NotContains(t, string(b), `"create"`)
	require.NotContains(t, string(b), `"update"`)
	require.JSONEq(t, `{"accountId":"u1","destroy":["cn1"]}`, string(b))
}

func TestQueryInvoke(t *testing.T) {
	req := &jmap.Request{}
	after := jmap.UTCDate(time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC))

	id := req.Invoke(&Query{
		Account: "u1",
		Filter: &FilterCondition{
			After: jmap.Some(after),
			Type:  TypeUpdated,
		},
		Sort: []*jmap.Comparator{
			{Property: "created"},
		},
		Limit: jmap.Uint64Ptr(10),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEventNotification/query",{"accountId":"u1","limit":10,"filter":{"after":"2026-09-01T00:00:00Z","type":"updated"},"sort":[{"property":"created"}]},"0"]]}`,
		string(data))
}

func TestQueryRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Query{}).Requires())
}

func TestQueryChangesInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&QueryChanges{
		Account:         "u1",
		Filter:          &FilterCondition{Type: TypeCreated},
		Sort:            []*jmap.Comparator{{Property: "created"}},
		SinceQueryState: "q1",
		MaxChanges:      jmap.Uint64Ptr(10),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["CalendarEventNotification/queryChanges",{"accountId":"u1","sinceQueryState":"q1","maxChanges":10,"filter":{"type":"created"},"sort":[{"property":"created"}]},"0"]]}`,
		string(data))
}

func TestQueryChangesRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&QueryChanges{}).Requires())
}
