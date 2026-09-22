package jscalendar

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

type UTCDateTime string
type LocalDateTime string
type Duration string
type SignedDuration string
type TimeZoneID string

type FreeBusyStatus string

const (
	FreeBusyFree FreeBusyStatus = "free"
	FreeBusyBusy FreeBusyStatus = "busy"
)

type Privacy string

const (
	PrivacyPublic  Privacy = "public"
	PrivacyPrivate Privacy = "private"
	PrivacySecret  Privacy = "secret"
)

type PatchObject map[string]any

type Event struct {
	Type    string `json:"@type,omitzero"`
	Version string `json:"version,omitzero"`
	UID     string `json:"uid,omitzero"`

	RelatedTo map[string]*Relation `json:"relatedTo,omitzero"`
	ProdID    string               `json:"prodId,omitzero"`
	Created   UTCDateTime          `json:"created,omitzero"`
	Updated   UTCDateTime          `json:"updated,omitzero"`
	Sequence  uint64               `json:"sequence,omitzero"`
	Method    string               `json:"method,omitzero"`

	Title                  string                      `json:"title,omitzero"`
	Description            string                      `json:"description,omitzero"`
	DescriptionContentType string                      `json:"descriptionContentType,omitzero"`
	ShowWithoutTime        bool                        `json:"showWithoutTime,omitzero"`
	Locations              map[string]*Location        `json:"locations,omitzero"`
	MainLocationID         string                      `json:"mainLocationId,omitzero"`
	VirtualLocations       map[string]*VirtualLocation `json:"virtualLocations,omitzero"`
	Links                  map[string]*Link            `json:"links,omitzero"`
	Locale                 string                      `json:"locale,omitzero"`
	Keywords               map[string]bool             `json:"keywords,omitzero"`
	Categories             map[string]bool             `json:"categories,omitzero"`
	Color                  string                      `json:"color,omitzero"`

	RecurrenceID         LocalDateTime                         `json:"recurrenceId,omitzero"`
	RecurrenceIDTimeZone TimeZoneID                            `json:"recurrenceIdTimeZone,omitzero"`
	RecurrenceRule       jmap.Optional[RecurrenceRule]         `json:"recurrenceRule,omitzero"`
	RecurrenceOverrides  jmap.Optional[map[string]PatchObject] `json:"recurrenceOverrides,omitzero"`

	OrganizerCalendarAddress string                  `json:"organizerCalendarAddress,omitzero"`
	SentBy                   string                  `json:"sentBy,omitzero"`
	Participants             map[string]*Participant `json:"participants,omitzero"`
	Alerts                   map[string]*Alert       `json:"alerts,omitzero"`

	TimeZone       TimeZoneID     `json:"timeZone,omitzero"`
	Start          LocalDateTime  `json:"start,omitzero"`
	Duration       Duration       `json:"duration,omitzero"`
	EndTimeZone    TimeZoneID     `json:"endTimeZone,omitzero"`
	Status         string         `json:"status,omitzero"`
	Priority       uint64         `json:"priority,omitzero"`
	FreeBusyStatus FreeBusyStatus `json:"freeBusyStatus,omitzero"`
	Privacy        Privacy        `json:"privacy,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Relation struct {
	Type     string          `json:"@type,omitzero"`
	Relation map[string]bool `json:"relation,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Link struct {
	Type        string            `json:"@type,omitzero"`
	Href        string            `json:"href,omitzero"`
	ContentType string            `json:"contentType,omitzero"`
	Size        *jmap.UnsignedInt `json:"size,omitzero"`
	Rel         string            `json:"rel,omitzero"`
	Title       string            `json:"title,omitzero"`
	Display     map[string]bool   `json:"display,omitzero"`
	BlobID      jmap.ID           `json:"blobId,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Location struct {
	Type          string           `json:"@type,omitzero"`
	Name          string           `json:"name,omitzero"`
	LocationTypes map[string]bool  `json:"locationTypes,omitzero"`
	Coordinates   string           `json:"coordinates,omitzero"`
	Links         map[string]*Link `json:"links,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type VirtualLocation struct {
	Type     string          `json:"@type,omitzero"`
	Name     string          `json:"name,omitzero"`
	URI      string          `json:"uri,omitzero"`
	Features map[string]bool `json:"features,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type NDay struct {
	Type        string `json:"@type,omitzero"`
	Day         string `json:"day,omitzero"`
	NthOfPeriod int64  `json:"nthOfPeriod,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type RecurrenceRule struct {
	Type           string        `json:"@type,omitzero"`
	Frequency      string        `json:"frequency,omitzero"`
	Interval       uint64        `json:"interval,omitzero"`
	RScale         string        `json:"rscale,omitzero"`
	Skip           string        `json:"skip,omitzero"`
	FirstDayOfWeek string        `json:"firstDayOfWeek,omitzero"`
	ByDay          []*NDay       `json:"byDay,omitzero"`
	ByMonthDay     []int64       `json:"byMonthDay,omitzero"`
	ByMonth        []string      `json:"byMonth,omitzero"`
	ByYearDay      []int64       `json:"byYearDay,omitzero"`
	ByWeekNo       []int64       `json:"byWeekNo,omitzero"`
	ByHour         []uint64      `json:"byHour,omitzero"`
	ByMinute       []uint64      `json:"byMinute,omitzero"`
	BySecond       []uint64      `json:"bySecond,omitzero"`
	BySetPosition  []int64       `json:"bySetPosition,omitzero"`
	Count          uint64        `json:"count,omitzero"`
	Until          LocalDateTime `json:"until,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Participant struct {
	Type                   string            `json:"@type,omitzero"`
	Name                   string            `json:"name,omitzero"`
	Email                  string            `json:"email,omitzero"`
	Description            string            `json:"description,omitzero"`
	DescriptionContentType string            `json:"descriptionContentType,omitzero"`
	CalendarAddress        string            `json:"calendarAddress,omitzero"`
	Kind                   string            `json:"kind,omitzero"`
	Roles                  map[string]bool   `json:"roles,omitzero"`
	ParticipationStatus    string            `json:"participationStatus,omitzero"`
	ExpectReply            bool              `json:"expectReply,omitzero"`
	SentBy                 string            `json:"sentBy,omitzero"`
	DelegatedTo            map[string]bool   `json:"delegatedTo,omitzero"`
	DelegatedFrom          map[string]bool   `json:"delegatedFrom,omitzero"`
	MemberOf               map[string]bool   `json:"memberOf,omitzero"`
	Links                  map[string]*Link  `json:"links,omitzero"`
	Progress               string            `json:"progress,omitzero"`
	PercentComplete        *jmap.UnsignedInt `json:"percentComplete,omitzero"`
	ScheduleSequence       jmap.UnsignedInt  `json:"scheduleSequence,omitzero"`
	ScheduleUpdated        UTCDateTime       `json:"scheduleUpdated,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Alert struct {
	Type         string               `json:"@type,omitzero"`
	Trigger      *Trigger             `json:"trigger,omitzero"`
	Acknowledged UTCDateTime          `json:"acknowledged,omitzero"`
	RelatedTo    map[string]*Relation `json:"relatedTo,omitzero"`
	Action       string               `json:"action,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Trigger struct {
	OffsetTrigger   *OffsetTrigger
	AbsoluteTrigger *AbsoluteTrigger
	Unknown         jsontext.Value
}

func (t Trigger) MarshalJSON() ([]byte, error) {
	switch {
	case len(t.Unknown) > 0:
		return []byte(t.Unknown), nil
	case t.AbsoluteTrigger != nil:
		return jsonv2.Marshal(t.AbsoluteTrigger)
	case t.OffsetTrigger != nil:
		return jsonv2.Marshal(t.OffsetTrigger)
	default:
		return []byte("null"), nil
	}
}

func (t *Trigger) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = Trigger{}
		return nil
	}

	var probe struct {
		Type string `json:"@type"`
	}
	if err := jsonv2.Unmarshal(data, &probe); err != nil {
		return err
	}

	switch probe.Type {
	case "", "OffsetTrigger":
		var offset OffsetTrigger
		if err := jsonv2.Unmarshal(data, &offset); err != nil {
			return err
		}
		t.OffsetTrigger = &offset
		t.AbsoluteTrigger = nil
		t.Unknown = nil
	case "AbsoluteTrigger":
		var absolute AbsoluteTrigger
		if err := jsonv2.Unmarshal(data, &absolute); err != nil {
			return err
		}
		t.AbsoluteTrigger = &absolute
		t.OffsetTrigger = nil
		t.Unknown = nil
	default:
		t.Unknown = append(t.Unknown[:0], data...)
		t.OffsetTrigger = nil
		t.AbsoluteTrigger = nil
	}

	return nil
}

type OffsetTrigger struct {
	Type       string         `json:"@type,omitzero"`
	Offset     SignedDuration `json:"offset,omitzero"`
	RelativeTo string         `json:"relativeTo,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type AbsoluteTrigger struct {
	Type string      `json:"@type,omitzero"`
	When UTCDateTime `json:"when,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}
