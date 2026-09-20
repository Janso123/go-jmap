package jscalendar

import "encoding/json"

type UTCDateTime string
type LocalDateTime string
type Duration string
type SignedDuration string
type TimeZoneID string

type PatchObject map[string]any

type Event struct {
	Type    string `json:"@type,omitempty"`
	Version string `json:"version,omitempty"`
	UID     string `json:"uid,omitempty"`

	RelatedTo map[string]*Relation `json:"relatedTo,omitempty"`
	ProdID    string               `json:"prodId,omitempty"`
	Created   UTCDateTime          `json:"created,omitempty"`
	Updated   UTCDateTime          `json:"updated,omitempty"`
	Sequence  uint64               `json:"sequence,omitempty"`
	Method    string               `json:"method,omitempty"`

	Title                  string                      `json:"title,omitempty"`
	Description            string                      `json:"description,omitempty"`
	DescriptionContentType string                      `json:"descriptionContentType,omitempty"`
	ShowWithoutTime        bool                        `json:"showWithoutTime,omitempty"`
	Locations              map[string]*Location        `json:"locations,omitempty"`
	MainLocationID         string                      `json:"mainLocationId,omitempty"`
	VirtualLocations       map[string]*VirtualLocation `json:"virtualLocations,omitempty"`
	Links                  map[string]*Link            `json:"links,omitempty"`
	Locale                 string                      `json:"locale,omitempty"`
	Keywords               map[string]bool             `json:"keywords,omitempty"`
	Categories             map[string]bool             `json:"categories,omitempty"`
	Color                  string                      `json:"color,omitempty"`

	RecurrenceID         LocalDateTime          `json:"recurrenceId,omitempty"`
	RecurrenceIDTimeZone TimeZoneID             `json:"recurrenceIdTimeZone,omitempty"`
	RecurrenceRule       *RecurrenceRule        `json:"recurrenceRule,omitempty"`
	RecurrenceOverrides  map[string]PatchObject `json:"recurrenceOverrides,omitempty"`

	OrganizerCalendarAddress string                  `json:"organizerCalendarAddress,omitempty"`
	SentBy                   string                  `json:"sentBy,omitempty"`
	Participants             map[string]*Participant `json:"participants,omitempty"`
	Alerts                   map[string]*Alert       `json:"alerts,omitempty"`

	TimeZone    TimeZoneID    `json:"timeZone,omitempty"`
	Start       LocalDateTime `json:"start,omitempty"`
	Duration    Duration      `json:"duration,omitempty"`
	EndTimeZone TimeZoneID    `json:"endTimeZone,omitempty"`
	Status      string        `json:"status,omitempty"`

	extensions map[string]json.RawMessage
}

type Relation struct {
	Type     string          `json:"@type,omitempty"`
	Relation map[string]bool `json:"relation,omitempty"`

	extensions map[string]json.RawMessage
}

type Link struct {
	Type        string `json:"@type,omitempty"`
	Href        string `json:"href,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Size        uint64 `json:"size,omitempty"`
	Rel         string `json:"rel,omitempty"`
	Title       string `json:"title,omitempty"`
	CID         string `json:"cid,omitempty"`
	Display     string `json:"display,omitempty"`

	extensions map[string]json.RawMessage
}

type Location struct {
	Type          string           `json:"@type,omitempty"`
	Name          string           `json:"name,omitempty"`
	LocationTypes map[string]bool  `json:"locationTypes,omitempty"`
	Coordinates   string           `json:"coordinates,omitempty"`
	Links         map[string]*Link `json:"links,omitempty"`

	extensions map[string]json.RawMessage
}

type VirtualLocation struct {
	Type     string          `json:"@type,omitempty"`
	Name     string          `json:"name,omitempty"`
	URI      string          `json:"uri,omitempty"`
	Features map[string]bool `json:"features,omitempty"`

	extensions map[string]json.RawMessage
}

type NDay struct {
	Type        string `json:"@type,omitempty"`
	Day         string `json:"day,omitempty"`
	NthOfPeriod int64  `json:"nthOfPeriod,omitempty"`

	extensions map[string]json.RawMessage
}

type RecurrenceRule struct {
	Type           string        `json:"@type,omitempty"`
	Frequency      string        `json:"frequency,omitempty"`
	Interval       uint64        `json:"interval,omitempty"`
	RScale         string        `json:"rscale,omitempty"`
	Skip           string        `json:"skip,omitempty"`
	FirstDayOfWeek string        `json:"firstDayOfWeek,omitempty"`
	ByDay          []*NDay       `json:"byDay,omitempty"`
	ByMonthDay     []int64       `json:"byMonthDay,omitempty"`
	ByMonth        []string      `json:"byMonth,omitempty"`
	ByYearDay      []int64       `json:"byYearDay,omitempty"`
	ByWeekNo       []int64       `json:"byWeekNo,omitempty"`
	ByHour         []uint64      `json:"byHour,omitempty"`
	ByMinute       []uint64      `json:"byMinute,omitempty"`
	BySecond       []uint64      `json:"bySecond,omitempty"`
	BySetPosition  []int64       `json:"bySetPosition,omitempty"`
	Count          uint64        `json:"count,omitempty"`
	Until          LocalDateTime `json:"until,omitempty"`

	extensions map[string]json.RawMessage
}

type Participant struct {
	Type                   string           `json:"@type,omitempty"`
	Name                   string           `json:"name,omitempty"`
	Email                  string           `json:"email,omitempty"`
	Description            string           `json:"description,omitempty"`
	DescriptionContentType string           `json:"descriptionContentType,omitempty"`
	CalendarAddress        string           `json:"calendarAddress,omitempty"`
	Kind                   string           `json:"kind,omitempty"`
	Roles                  map[string]bool  `json:"roles,omitempty"`
	ParticipationStatus    string           `json:"participationStatus,omitempty"`
	ExpectReply            bool             `json:"expectReply,omitempty"`
	SentBy                 string           `json:"sentBy,omitempty"`
	DelegatedTo            map[string]bool  `json:"delegatedTo,omitempty"`
	DelegatedFrom          map[string]bool  `json:"delegatedFrom,omitempty"`
	MemberOf               map[string]bool  `json:"memberOf,omitempty"`
	Links                  map[string]*Link `json:"links,omitempty"`
	Progress               string           `json:"progress,omitempty"`
	PercentComplete        uint64           `json:"percentComplete,omitempty"`

	extensions map[string]json.RawMessage
}

type Alert struct {
	Type         string               `json:"@type,omitempty"`
	Trigger      *Trigger             `json:"trigger,omitempty"`
	Acknowledged UTCDateTime          `json:"acknowledged,omitempty"`
	RelatedTo    map[string]*Relation `json:"relatedTo,omitempty"`
	Action       string               `json:"action,omitempty"`

	extensions map[string]json.RawMessage
}

type Trigger struct {
	OffsetTrigger   *OffsetTrigger
	AbsoluteTrigger *AbsoluteTrigger
	Unknown         json.RawMessage
}

type OffsetTrigger struct {
	Type       string         `json:"@type,omitempty"`
	Offset     SignedDuration `json:"offset,omitempty"`
	RelativeTo string         `json:"relativeTo,omitempty"`

	extensions map[string]json.RawMessage
}

type AbsoluteTrigger struct {
	Type string      `json:"@type,omitempty"`
	When UTCDateTime `json:"when,omitempty"`

	extensions map[string]json.RawMessage
}
