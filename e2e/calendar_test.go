//go:build e2e

package e2e

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/calendarevent"
	"github.com/Janso123/go-jmap/calendar/calendareventnotification"
	"github.com/Janso123/go-jmap/calendar/calendars"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
	"github.com/Janso123/go-jmap/calendar/participantidentity"
)

const calendarRFC = "draft-ietf-jmap-calendars-29"

func TestCalendar(t *testing.T) {
	sc := &scenario{name: "Calendar"}
	if alice == nil || alice.Client == nil || bob == nil {
		t.Fatal("account session missing")
	}
	if !hasCap(alice.Client, calendar.URI) {
		skipRest(t, sc, calendarSteps, "missing "+string(calendar.URI))
		return
	}

	aliceID, err := alice.Client.PrimaryAccount(calendar.URI)
	if err != nil {
		skipRest(t, sc, calendarSteps, err.Error())
		return
	}

	call[*calendars.GetResponse](t, sc, step{
		RFC: calendarRFC, Method: "Calendar/get", Account: alice.Name, Request: "list",
	}, alice.Client, []jmap.URI{calendar.URI}, false, &calendars.Get{Account: aliceID}, func(resp *calendars.GetResponse) (string, error) {
		if resp == nil || len(resp.List) < 1 {
			return "", errString("no calendar")
		}
		return "calendars=" + strconv.Itoa(len(resp.List)), nil
	})

	var calID jmap.ID
	call[*calendars.SetResponse](t, sc, step{
		RFC: calendarRFC, Method: "Calendar/set", Account: alice.Name, Request: "create name=E2E",
	}, alice.Client, []jmap.URI{calendar.URI}, false, &calendars.Set{
		Account: aliceID,
		Create: jmap.Some(map[jmap.ID]*calendars.Calendar{
			"c1": {Name: "E2E"},
		}),
	}, func(resp *calendars.SetResponse) (string, error) {
		id, err := createdCalendarID(resp, "c1")
		if err != nil {
			return "", err
		}
		calID = id
		return "id=" + string(calID), nil
	})

	var stateBefore string
	call[*calendarevent.GetResponse](t, sc, step{
		RFC: calendarRFC, Method: "CalendarEvent/get", Account: alice.Name, Request: "ids omitted",
	}, alice.Client, []jmap.URI{calendar.URI}, false, &calendarevent.Get{Account: aliceID}, func(resp *calendarevent.GetResponse) (string, error) {
		if resp == nil || resp.State == "" {
			return "", errString("event state missing")
		}
		stateBefore = resp.State
		return "state=" + stateBefore, nil
	})

	eventIDs := make([]jmap.ID, 5)
	batch := maxInSet(alice.Client)
	for start := 1; start <= 5; start += batch {
		end := start + batch - 1
		if end > 5 {
			end = 5
		}
		create := map[jmap.ID]*calendarevent.CalendarEvent{}
		var keys []jmap.ID
		for n := start; n <= end; n++ {
			key := jmap.ID(fmt.Sprintf("e%d", n))
			keys = append(keys, key)
			var participants map[string]*jscalendar.Participant
			if n == 1 {
				participants = map[string]*jscalendar.Participant{
					"bob": {CalendarAddress: "mailto:bob@example.org", Email: bob.Email},
				}
			}
			create[key] = &calendarevent.CalendarEvent{
				CalendarIDs: jmap.Some(map[jmap.ID]bool{calID: true}),
				Event: jscalendar.Event{
					Type:         "Event",
					Title:        fmt.Sprintf("E2E Event %d", n),
					Start:        "2026-09-22T09:00:00",
					Duration:     "PT1H",
					TimeZone:     "Etc/UTC",
					Participants: participants,
				},
			}
		}
		call[*calendarevent.SetResponse](t, sc, step{
			RFC: calendarRFC, Method: "CalendarEvent/set", Account: alice.Name,
			Request: fmt.Sprintf("create events %d-%d", start, end),
		}, alice.Client, []jmap.URI{calendar.URI}, false, &calendarevent.Set{
			Account: aliceID,
			Create:  jmap.Some(create),
		}, func(resp *calendarevent.SetResponse) (string, error) {
			made, err := createdEventIDs(resp, keys)
			if err != nil {
				return "", err
			}
			for n := start; n <= end; n++ {
				eventIDs[n-1] = made[jmap.ID(fmt.Sprintf("e%d", n))]
			}
			return "created=" + strconv.Itoa(len(made)), nil
		})
	}
	bobEventID := eventIDs[0]
	dropID := eventIDs[4]

	call[*calendarevent.QueryResponse](t, sc, step{
		RFC: calendarRFC, Method: "CalendarEvent/query", Account: alice.Name, Request: "inCalendar window",
	}, alice.Client, []jmap.URI{calendar.URI}, false, &calendarevent.Query{
		Account: aliceID,
		Filter: &calendarevent.FilterCondition{
			InCalendar: calID,
			After:      "2026-09-22T08:00:00",
			Before:     "2026-09-22T12:00:00",
		},
	}, func(resp *calendarevent.QueryResponse) (string, error) {
		if resp == nil || len(resp.IDs) != 5 || !containsAll(resp.IDs, eventIDs) {
			return "", errString("expected 5 events in the calendar")
		}
		return "ids=5", nil
	})

	call[*calendarevent.GetResponse](t, sc, step{
		RFC: calendarRFC, Method: "CalendarEvent/get", Account: alice.Name, Request: "id=" + string(bobEventID),
	}, alice.Client, []jmap.URI{calendar.URI}, false, &calendarevent.Get{
		Account: aliceID,
		IDs:     jmap.Some([]jmap.ID{bobEventID}),
	}, func(resp *calendarevent.GetResponse) (string, error) {
		if resp == nil || len(resp.List) == 0 || !hasCalendarAddress(resp.List[0].Participants, "mailto:bob@example.org") {
			return "", errString("participant mailto:bob@example.org missing")
		}
		return "calendarAddress=mailto:bob@example.org", nil
	})

	call[*calendarevent.ChangesResponse](t, sc, step{
		RFC: calendarRFC, Method: "CalendarEvent/changes", Account: alice.Name, Request: "sinceState=" + stateBefore,
	}, alice.Client, []jmap.URI{calendar.URI}, false, &calendarevent.Changes{
		Account:    aliceID,
		SinceState: stateBefore,
	}, func(resp *calendarevent.ChangesResponse) (string, error) {
		if resp == nil || !containsAll(resp.Created, eventIDs) {
			more := false
			if resp != nil {
				more = resp.HasMoreChanges
			}
			return "", errString(fmt.Sprintf("created missing events hasMore=%v", more))
		}
		return "created=5", nil
	})

	call[*calendarevent.SetResponse](t, sc, step{
		RFC: calendarRFC, Method: "CalendarEvent/set", Account: alice.Name, Request: "update title",
	}, alice.Client, []jmap.URI{calendar.URI}, false, &calendarevent.Set{
		Account: aliceID,
		Update:  jmap.Some(map[jmap.ID]jmap.Patch{dropID: {"title": "E2E Renamed"}}),
	}, func(resp *calendarevent.SetResponse) (string, error) {
		if err := rejectEventSet(resp, dropID); err != nil {
			return "", err
		}
		return "updated=" + string(dropID), nil
	})

	call[*calendarevent.GetResponse](t, sc, step{
		RFC: calendarRFC, Method: "CalendarEvent/get", Account: alice.Name, Request: "id=" + string(dropID),
	}, alice.Client, []jmap.URI{calendar.URI}, false, &calendarevent.Get{
		Account: aliceID,
		IDs:     jmap.Some([]jmap.ID{dropID}),
	}, func(resp *calendarevent.GetResponse) (string, error) {
		if resp == nil || len(resp.List) == 0 || resp.List[0].Title != "E2E Renamed" {
			return "", errString("title mismatch")
		}
		return "title=" + resp.List[0].Title, nil
	})

	call[*calendarevent.SetResponse](t, sc, step{
		RFC: calendarRFC, Method: "CalendarEvent/set", Account: alice.Name, Request: "destroy id=" + string(dropID),
	}, alice.Client, []jmap.URI{calendar.URI}, false, &calendarevent.Set{
		Account: aliceID,
		Destroy: jmap.Some([]jmap.ID{dropID}),
	}, func(resp *calendarevent.SetResponse) (string, error) {
		if err := rejectEventSet(resp, dropID); err != nil {
			return "", err
		}
		return "destroyed=" + string(dropID), nil
	})

	kept := eventIDs[:4]
	call[*calendarevent.QueryResponse](t, sc, step{
		RFC: calendarRFC, Method: "CalendarEvent/query", Account: alice.Name, Request: "inCalendar window",
	}, alice.Client, []jmap.URI{calendar.URI}, false, &calendarevent.Query{
		Account: aliceID,
		Filter: &calendarevent.FilterCondition{
			InCalendar: calID,
			After:      "2026-09-22T08:00:00",
			Before:     "2026-09-22T12:00:00",
		},
	}, func(resp *calendarevent.QueryResponse) (string, error) {
		if resp == nil || len(resp.IDs) != 4 || containsID(resp.IDs, dropID) || !containsAll(resp.IDs, kept) {
			return "", errString("expected 4 events in the calendar")
		}
		return "ids=4", nil
	})

	call[*participantidentity.GetResponse](t, sc, step{
		RFC: calendarRFC, Method: "ParticipantIdentity/get", Account: alice.Name, Request: "ids omitted",
	}, alice.Client, []jmap.URI{calendar.URI}, false, &participantidentity.Get{Account: aliceID}, func(resp *participantidentity.GetResponse) (string, error) {
		n := 0
		if resp != nil {
			n = len(resp.List)
		}
		return "len=" + strconv.Itoa(n), nil
	})

	call[*calendareventnotification.QueryResponse](t, sc, step{
		RFC: calendarRFC, Method: "CalendarEventNotification/query", Account: alice.Name, Request: "empty filter",
	}, alice.Client, []jmap.URI{calendar.URI}, true, &calendareventnotification.Query{
		Account: aliceID,
		Filter:  &calendareventnotification.FilterCondition{},
	}, func(resp *calendareventnotification.QueryResponse) (string, error) {
		n := 0
		if resp != nil {
			n = len(resp.IDs)
		}
		return "len=" + strconv.Itoa(n), nil
	})

	if sc.failed {
		t.Fail()
	}
}

var calendarSteps = []step{
	{RFC: calendarRFC, Method: "Calendar/get", Account: "alice", Request: "list"},
	{RFC: calendarRFC, Method: "Calendar/set", Account: "alice", Request: "create name=E2E"},
	{RFC: calendarRFC, Method: "CalendarEvent/get", Account: "alice", Request: "ids omitted"},
	{RFC: calendarRFC, Method: "CalendarEvent/set", Account: "alice", Request: "create 5 events"},
	{RFC: calendarRFC, Method: "CalendarEvent/query", Account: "alice", Request: "inCalendar window"},
	{RFC: calendarRFC, Method: "CalendarEvent/get", Account: "alice", Request: "participant"},
	{RFC: calendarRFC, Method: "CalendarEvent/changes", Account: "alice", Request: "sinceState"},
	{RFC: calendarRFC, Method: "CalendarEvent/set", Account: "alice", Request: "update title"},
	{RFC: calendarRFC, Method: "CalendarEvent/set", Account: "alice", Request: "destroy"},
	{RFC: calendarRFC, Method: "CalendarEvent/query", Account: "alice", Request: "inCalendar window"},
	{RFC: calendarRFC, Method: "ParticipantIdentity/get", Account: "alice", Request: "ids omitted"},
	{RFC: calendarRFC, Method: "CalendarEventNotification/query", Account: "alice", Request: "empty filter"},
}

func createdCalendarID(resp *calendars.SetResponse, key jmap.ID) (jmap.ID, error) {
	if resp == nil {
		return "", errString("empty set response")
	}
	if nc, ok := resp.NotCreated.Value(); ok {
		if se := nc[key]; se != nil {
			return "", errString(setErrorText(se))
		}
	}
	created, ok := resp.Created[key]
	if !ok || created.ID == "" {
		return "", errString("calendar was not created")
	}
	return created.ID, nil
}

func createdEventIDs(resp *calendarevent.SetResponse, keys []jmap.ID) (map[jmap.ID]jmap.ID, error) {
	if resp == nil {
		return nil, errString("empty set response")
	}
	if nc, ok := resp.NotCreated.Value(); ok && len(nc) > 0 {
		return nil, errString(setErrorSummary(nc))
	}
	out := make(map[jmap.ID]jmap.ID, len(keys))
	for _, key := range keys {
		created, ok := resp.Created[key]
		id, hasID := created.ID.Value()
		if !ok || !hasID || id == "" {
			return nil, errString("event " + string(key) + " was not created")
		}
		out[key] = id
	}
	return out, nil
}

func rejectEventSet(resp *calendarevent.SetResponse, id jmap.ID) error {
	if resp == nil {
		return errString("empty set response")
	}
	if nu, ok := resp.NotUpdated.Value(); ok {
		if se := nu[id]; se != nil {
			return errString(setErrorText(se))
		}
	}
	if nd, ok := resp.NotDestroyed.Value(); ok {
		if se := nd[id]; se != nil {
			return errString(setErrorText(se))
		}
	}
	return nil
}

func hasCalendarAddress(participants map[string]*jscalendar.Participant, address string) bool {
	for _, participant := range participants {
		if participant != nil && participant.CalendarAddress == address {
			return true
		}
	}
	return false
}
