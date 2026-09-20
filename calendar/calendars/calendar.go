package calendars

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/calendar/jscalendar"
)

func init() {
	jmap.RegisterMethod("Calendar/get", newGetResponse)
	jmap.RegisterMethod("Calendar/changes", newChangesResponse)
	jmap.RegisterMethod("Calendar/set", newSetResponse)
}

// Calendar is a named collection of events.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-4
type Calendar struct {
	ID jmap.ID `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Description *string `json:"description,omitempty"`

	Color *string `json:"color,omitempty"`

	SortOrder uint64 `json:"sortOrder,omitempty"`

	IsSubscribed bool `json:"isSubscribed,omitempty"`

	IsVisible bool `json:"isVisible,omitempty"`

	IsDefault bool `json:"isDefault,omitempty"`

	IncludeInAvailability IncludeInAvailability `json:"includeInAvailability,omitempty"`

	DefaultAlertsWithTime map[jmap.ID]*jscalendar.Alert `json:"defaultAlertsWithTime,omitempty"`

	DefaultAlertsWithoutTime map[jmap.ID]*jscalendar.Alert `json:"defaultAlertsWithoutTime,omitempty"`

	TimeZone jscalendar.TimeZoneID `json:"timeZone,omitempty"`

	ShareWith map[jmap.ID]*Rights `json:"shareWith,omitempty"`

	MyRights *Rights `json:"myRights,omitempty"`
}

type IncludeInAvailability string

const (
	IncludeInAvailabilityAll       IncludeInAvailability = "all"
	IncludeInAvailabilityAttending IncludeInAvailability = "attending"
	IncludeInAvailabilityNone      IncludeInAvailability = "none"
)

// Rights is the set of permissions the user has in relation to a Calendar.
type Rights struct {
	MayReadFreeBusy bool `json:"mayReadFreeBusy,omitempty"`

	MayReadItems bool `json:"mayReadItems,omitempty"`

	MayWriteAll bool `json:"mayWriteAll,omitempty"`

	MayWriteOwn bool `json:"mayWriteOwn,omitempty"`

	MayUpdatePrivate bool `json:"mayUpdatePrivate,omitempty"`

	MayRSVP bool `json:"mayRSVP,omitempty"`

	MayShare bool `json:"mayShare,omitempty"`

	MayDelete bool `json:"mayDelete,omitempty"`
}
