package calendars

import (
	"encoding/json"
	"testing"

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
		MaxChanges: 50,
	})
	assert.Equal(t, "0", id)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:calendars"],"methodCalls":[["Calendar/changes",{"accountId":"u1","sinceState":"s1","maxChanges":50},"0"]]}`,
		string(data))
}

func TestChangesRequiresCalendarsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{calendar.URI}, (&Changes{}).Requires())
}
