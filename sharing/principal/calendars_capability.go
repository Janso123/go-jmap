package principal

import (
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
)

// CalendarsCapability is the Principal capabilities value for
// urn:ietf:params:jmap:calendars.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-2.1
//
// Do not pass this type to jmap.RegisterCapability. That URI is already
// registered as calendar.Capability.
type CalendarsCapability struct {
	// AccountID is Id|null.
	AccountID jmap.Optional[jmap.ID] `json:"accountId,omitzero"`

	MayGetAvailability bool `json:"mayGetAvailability"`

	MayShareWith bool `json:"mayShareWith"`

	CalendarAddress string `json:"calendarAddress"`
}

// CalendarsCapability decodes the urn:ietf:params:jmap:calendars entry of
// p.Capabilities. ok is false when the entry is absent.
func (p *Principal) CalendarsCapability() (*CalendarsCapability, bool, error) {
	if p == nil {
		return nil, false, nil
	}
	raw, ok := p.Capabilities[calendar.URI]
	if !ok {
		return nil, false, nil
	}
	var cap CalendarsCapability
	if err := jsonv2.Unmarshal(raw, &cap); err != nil {
		return nil, false, err
	}
	return &cap, true, nil
}
