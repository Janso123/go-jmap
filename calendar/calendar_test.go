package calendar_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapabilityUnmarshalFromSessionAndAccount(t *testing.T) {
	raw := `{
	  "capabilities": {
	    "urn:ietf:params:jmap:calendars": {},
	    "urn:ietf:params:jmap:calendars:parse": {},
	    "urn:ietf:params:jmap:principals:availability": {}
	  },
	  "accounts": {
	    "u1": {
	      "name": "user@example.com",
	      "isPersonal": true,
	      "isReadOnly": false,
	      "accountCapabilities": {
	        "urn:ietf:params:jmap:calendars": {
	          "maxCalendarsPerEvent": 7,
	          "minDateTime": "1900-01-01T00:00:00Z",
	          "maxDateTime": "2200-01-01T00:00:00Z",
	          "maxExpandedQueryDuration": "P365D",
	          "maxParticipantsPerEvent": 250,
	          "mayCreateCalendar": true
	        },
	        "urn:ietf:params:jmap:calendars:parse": {},
	        "urn:ietf:params:jmap:principals:availability": {
	          "maxAvailabilityDuration": "P31D"
	        }
	      }
	    }
	  },
	  "primaryAccounts": {},
	  "username": "u",
	  "apiUrl": "https://server.example.com/jmap/api/",
	  "downloadUrl": "https://server.example.com/dl/{accountId}/{blobId}/{name}?accept={type}",
	  "uploadUrl": "https://server.example.com/upload/{accountId}/",
	  "eventSourceUrl": "https://server.example.com/es/",
	  "state": "s1"
	}`

	var session jmap.Session
	require.NoError(t, json.Unmarshal([]byte(raw), &session))

	sessionCalendars, ok := session.Capabilities[calendar.URI].(*calendar.Capability)
	require.True(t, ok)
	assert.Nil(t, sessionCalendars.MaxCalendarsPerEvent)
	assert.Empty(t, sessionCalendars.MinDateTime)
	assert.Empty(t, sessionCalendars.MaxDateTime)
	assert.Empty(t, sessionCalendars.MaxExpandedQueryDuration)
	assert.Nil(t, sessionCalendars.MaxParticipantsPerEvent)
	assert.Nil(t, sessionCalendars.MayCreateCalendar)

	_, ok = session.Capabilities[calendar.ParseURI].(*calendar.ParseCapability)
	require.True(t, ok)

	sessionAvailability, ok := session.Capabilities[calendar.AvailabilityURI].(*calendar.AvailabilityCapability)
	require.True(t, ok)
	assert.Empty(t, sessionAvailability.MaxAvailabilityDuration)

	accountCalendars, ok := session.Accounts["u1"].Capabilities[calendar.URI].(*calendar.Capability)
	require.True(t, ok)
	require.NotNil(t, accountCalendars.MaxCalendarsPerEvent)
	assert.Equal(t, uint64(7), *accountCalendars.MaxCalendarsPerEvent)
	assert.Equal(t, jscalendar.UTCDateTime("1900-01-01T00:00:00Z"), accountCalendars.MinDateTime)
	assert.Equal(t, jscalendar.UTCDateTime("2200-01-01T00:00:00Z"), accountCalendars.MaxDateTime)
	assert.Equal(t, jscalendar.Duration("P365D"), accountCalendars.MaxExpandedQueryDuration)
	require.NotNil(t, accountCalendars.MaxParticipantsPerEvent)
	assert.Equal(t, uint64(250), *accountCalendars.MaxParticipantsPerEvent)
	require.NotNil(t, accountCalendars.MayCreateCalendar)
	assert.True(t, *accountCalendars.MayCreateCalendar)

	_, ok = session.Accounts["u1"].Capabilities[calendar.ParseURI].(*calendar.ParseCapability)
	require.True(t, ok)

	accountAvailability, ok := session.Accounts["u1"].Capabilities[calendar.AvailabilityURI].(*calendar.AvailabilityCapability)
	require.True(t, ok)
	assert.Equal(t, jscalendar.Duration("P31D"), accountAvailability.MaxAvailabilityDuration)
}
