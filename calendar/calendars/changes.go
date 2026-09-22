package calendars

import "github.com/Janso123/go-jmap"

// Changes gets calendar changes for the whole account.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-4.2
type Changes struct {
	jmap.Changes[Calendar]
}

// ChangesResponse is the result of Calendar/changes.
type ChangesResponse = jmap.ChangesResponse
