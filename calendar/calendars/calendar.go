package calendars

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
)

func init() {
	jmap.RegisterObject[Calendar](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodSet,
	)
}

// Calendar is a named collection of events.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-4
type Calendar struct {
	ID jmap.ID `json:"id,omitzero"`

	Name string `json:"name,omitzero"`

	Description jmap.Optional[string] `json:"description,omitzero"`

	Color jmap.Optional[string] `json:"color,omitzero"`

	SortOrder uint64 `json:"sortOrder,omitzero"`

	IsSubscribed *bool `json:"isSubscribed,omitzero"`

	IsVisible *bool `json:"isVisible,omitzero"`

	IsDefault *bool `json:"isDefault,omitzero"`

	IncludeInAvailability IncludeInAvailability `json:"includeInAvailability,omitzero"`

	DefaultAlertsWithTime jmap.Optional[map[jmap.ID]*jscalendar.Alert] `json:"defaultAlertsWithTime,omitzero"`

	DefaultAlertsWithoutTime jmap.Optional[map[jmap.ID]*jscalendar.Alert] `json:"defaultAlertsWithoutTime,omitzero"`

	TimeZone jmap.Optional[jscalendar.TimeZoneID] `json:"timeZone,omitzero"`

	ShareWith jmap.Optional[map[jmap.ID]*Rights] `json:"shareWith,omitzero"`

	MyRights *Rights `json:"myRights,omitzero"`
}

func (Calendar) JMAPType() string { return "Calendar" }

func (Calendar) JMAPCreatable() {}

func (Calendar) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }

type IncludeInAvailability string

const (
	IncludeInAvailabilityAll       IncludeInAvailability = "all"
	IncludeInAvailabilityAttending IncludeInAvailability = "attending"
	IncludeInAvailabilityNone      IncludeInAvailability = "none"
)

// Rights is the set of permissions the user has in relation to a Calendar.
type Rights struct {
	MayReadFreeBusy bool `json:"mayReadFreeBusy"`

	MayReadItems bool `json:"mayReadItems"`

	MayWriteAll bool `json:"mayWriteAll"`

	MayWriteOwn bool `json:"mayWriteOwn"`

	MayUpdatePrivate bool `json:"mayUpdatePrivate"`

	MayRSVP bool `json:"mayRSVP"`

	MayShare bool `json:"mayShare"`

	MayDelete bool `json:"mayDelete"`
}
