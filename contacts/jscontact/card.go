package jscontact

import "encoding/json"

type UTCDateTime string

type PatchObject map[string]any

type Card struct {
	Type    string `json:"@type,omitempty"`
	Version string `json:"version,omitempty"`
	UID     string `json:"uid,omitempty"`

	Created  UTCDateTime     `json:"created,omitempty"`
	Kind     string          `json:"kind,omitempty"`
	Language string          `json:"language,omitempty"`
	Members  map[string]bool `json:"members,omitempty"`
	ProdID   string          `json:"prodId,omitempty"`

	RelatedTo map[string]*Relation `json:"relatedTo,omitempty"`
	Updated   UTCDateTime          `json:"updated,omitempty"`

	Name          *Name                    `json:"name,omitempty"`
	Nicknames     map[string]*Nickname     `json:"nicknames,omitempty"`
	Organizations map[string]*Organization `json:"organizations,omitempty"`
	SpeakToAs     *SpeakToAs               `json:"speakToAs,omitempty"`
	Titles        map[string]*Title        `json:"titles,omitempty"`

	Emails              map[string]*EmailAddress      `json:"emails,omitempty"`
	OnlineServices      map[string]*OnlineService     `json:"onlineServices,omitempty"`
	Phones              map[string]*Phone             `json:"phones,omitempty"`
	PreferredLanguages  map[string]*LanguagePref      `json:"preferredLanguages,omitempty"`
	Calendars           map[string]*Calendar          `json:"calendars,omitempty"`
	SchedulingAddresses map[string]*SchedulingAddress `json:"schedulingAddresses,omitempty"`
	Addresses           map[string]*Address           `json:"addresses,omitempty"`
	CryptoKeys          map[string]*CryptoKey         `json:"cryptoKeys,omitempty"`
	Directories         map[string]*Directory         `json:"directories,omitempty"`
	Links               map[string]*Link              `json:"links,omitempty"`
	Media               map[string]*Media             `json:"media,omitempty"`

	Localizations map[string]PatchObject `json:"localizations,omitempty"`

	Anniversaries map[string]*Anniversary  `json:"anniversaries,omitempty"`
	Keywords      map[string]bool          `json:"keywords,omitempty"`
	Notes         map[string]*Note         `json:"notes,omitempty"`
	PersonalInfo  map[string]*PersonalInfo `json:"personalInfo,omitempty"`

	extensions map[string]json.RawMessage
}

type Relation struct {
	Type     string          `json:"@type,omitempty"`
	Relation map[string]bool `json:"relation,omitempty"`

	extensions map[string]json.RawMessage
}

type Name struct {
	Type             string            `json:"@type,omitempty"`
	Components       []*NameComponent  `json:"components,omitempty"`
	IsOrdered        bool              `json:"isOrdered,omitempty"`
	DefaultSeparator string            `json:"defaultSeparator,omitempty"`
	Full             string            `json:"full,omitempty"`
	SortAs           map[string]string `json:"sortAs,omitempty"`
	PhoneticScript   string            `json:"phoneticScript,omitempty"`
	PhoneticSystem   string            `json:"phoneticSystem,omitempty"`

	extensions map[string]json.RawMessage
}

type NameComponent struct {
	Type     string `json:"@type,omitempty"`
	Value    string `json:"value,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Phonetic string `json:"phonetic,omitempty"`

	extensions map[string]json.RawMessage
}

type Nickname struct {
	Type     string          `json:"@type,omitempty"`
	Name     string          `json:"name,omitempty"`
	Contexts map[string]bool `json:"contexts,omitempty"`
	Pref     uint64          `json:"pref,omitempty"`

	extensions map[string]json.RawMessage
}

type Organization struct {
	Type     string          `json:"@type,omitempty"`
	Name     string          `json:"name,omitempty"`
	Units    []*OrgUnit      `json:"units,omitempty"`
	SortAs   string          `json:"sortAs,omitempty"`
	Contexts map[string]bool `json:"contexts,omitempty"`

	extensions map[string]json.RawMessage
}

type OrgUnit struct {
	Type   string `json:"@type,omitempty"`
	Name   string `json:"name,omitempty"`
	SortAs string `json:"sortAs,omitempty"`

	extensions map[string]json.RawMessage
}

type SpeakToAs struct {
	Type              string               `json:"@type,omitempty"`
	GrammaticalGender string               `json:"grammaticalGender,omitempty"`
	Pronouns          map[string]*Pronouns `json:"pronouns,omitempty"`

	extensions map[string]json.RawMessage
}

type Pronouns struct {
	Type     string          `json:"@type,omitempty"`
	Pronouns string          `json:"pronouns,omitempty"`
	Contexts map[string]bool `json:"contexts,omitempty"`
	Pref     uint64          `json:"pref,omitempty"`

	extensions map[string]json.RawMessage
}

type Title struct {
	Type           string `json:"@type,omitempty"`
	Name           string `json:"name,omitempty"`
	Kind           string `json:"kind,omitempty"`
	OrganizationID string `json:"organizationId,omitempty"`

	extensions map[string]json.RawMessage
}

type EmailAddress struct {
	Type     string          `json:"@type,omitempty"`
	Address  string          `json:"address,omitempty"`
	Contexts map[string]bool `json:"contexts,omitempty"`
	Pref     uint64          `json:"pref,omitempty"`
	Label    string          `json:"label,omitempty"`

	extensions map[string]json.RawMessage
}

type OnlineService struct {
	Type     string          `json:"@type,omitempty"`
	Service  string          `json:"service,omitempty"`
	URI      string          `json:"uri,omitempty"`
	User     string          `json:"user,omitempty"`
	Contexts map[string]bool `json:"contexts,omitempty"`
	Pref     uint64          `json:"pref,omitempty"`
	Label    string          `json:"label,omitempty"`

	extensions map[string]json.RawMessage
}

type Phone struct {
	Type     string          `json:"@type,omitempty"`
	Number   string          `json:"number,omitempty"`
	Features map[string]bool `json:"features,omitempty"`
	Contexts map[string]bool `json:"contexts,omitempty"`
	Pref     uint64          `json:"pref,omitempty"`
	Label    string          `json:"label,omitempty"`

	extensions map[string]json.RawMessage
}

type LanguagePref struct {
	Type     string          `json:"@type,omitempty"`
	Language string          `json:"language,omitempty"`
	Contexts map[string]bool `json:"contexts,omitempty"`
	Pref     uint64          `json:"pref,omitempty"`

	extensions map[string]json.RawMessage
}

type Resource struct {
	Type      string          `json:"@type,omitempty"`
	Kind      string          `json:"kind,omitempty"`
	URI       string          `json:"uri,omitempty"`
	MediaType string          `json:"mediaType,omitempty"`
	Contexts  map[string]bool `json:"contexts,omitempty"`
	Pref      uint64          `json:"pref,omitempty"`
	Label     string          `json:"label,omitempty"`

	extensions map[string]json.RawMessage
}

type Calendar struct {
	Resource

	extensions map[string]json.RawMessage
}

type SchedulingAddress struct {
	Type     string          `json:"@type,omitempty"`
	URI      string          `json:"uri,omitempty"`
	Contexts map[string]bool `json:"contexts,omitempty"`
	Pref     uint64          `json:"pref,omitempty"`
	Label    string          `json:"label,omitempty"`

	extensions map[string]json.RawMessage
}

type Address struct {
	Type             string              `json:"@type,omitempty"`
	Components       []*AddressComponent `json:"components,omitempty"`
	IsOrdered        bool                `json:"isOrdered,omitempty"`
	CountryCode      string              `json:"countryCode,omitempty"`
	Coordinates      string              `json:"coordinates,omitempty"`
	TimeZone         string              `json:"timeZone,omitempty"`
	Contexts         map[string]bool     `json:"contexts,omitempty"`
	Full             string              `json:"full,omitempty"`
	DefaultSeparator string              `json:"defaultSeparator,omitempty"`
	Pref             uint64              `json:"pref,omitempty"`
	PhoneticScript   string              `json:"phoneticScript,omitempty"`
	PhoneticSystem   string              `json:"phoneticSystem,omitempty"`

	extensions map[string]json.RawMessage
}

type AddressComponent struct {
	Type     string `json:"@type,omitempty"`
	Value    string `json:"value,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Phonetic string `json:"phonetic,omitempty"`

	extensions map[string]json.RawMessage
}

type CryptoKey struct {
	Resource

	extensions map[string]json.RawMessage
}

type Directory struct {
	Resource
	ListAs uint64 `json:"listAs,omitempty"`

	extensions map[string]json.RawMessage
}

type Link struct {
	Resource

	extensions map[string]json.RawMessage
}

type Media struct {
	Resource

	extensions map[string]json.RawMessage
}

type Anniversary struct {
	Type  string           `json:"@type,omitempty"`
	Kind  string           `json:"kind,omitempty"`
	Date  *AnniversaryDate `json:"date,omitempty"`
	Place *Address         `json:"place,omitempty"`

	extensions map[string]json.RawMessage
}

type AnniversaryDate struct {
	PartialDate *PartialDate
	Timestamp   *Timestamp
}

func (d AnniversaryDate) MarshalJSON() ([]byte, error) {
	switch {
	case d.Timestamp != nil:
		return json.Marshal(d.Timestamp)
	case d.PartialDate != nil:
		return json.Marshal(d.PartialDate)
	default:
		return []byte("null"), nil
	}
}

func (d *AnniversaryDate) UnmarshalJSON(data []byte) error {
	var probe struct {
		Type string `json:"@type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}
	if probe.Type == "Timestamp" {
		var ts Timestamp
		if err := json.Unmarshal(data, &ts); err != nil {
			return err
		}
		d.Timestamp = &ts
		d.PartialDate = nil
		return nil
	}
	var pd PartialDate
	if err := json.Unmarshal(data, &pd); err != nil {
		return err
	}
	d.PartialDate = &pd
	d.Timestamp = nil
	return nil
}

type PartialDate struct {
	Type          string `json:"@type,omitempty"`
	Year          uint64 `json:"year,omitempty"`
	Month         uint64 `json:"month,omitempty"`
	Day           uint64 `json:"day,omitempty"`
	CalendarScale string `json:"calendarScale,omitempty"`

	extensions map[string]json.RawMessage
}

type Timestamp struct {
	Type string      `json:"@type,omitempty"`
	UTC  UTCDateTime `json:"utc,omitempty"`

	extensions map[string]json.RawMessage
}

type Note struct {
	Type    string      `json:"@type,omitempty"`
	Note    string      `json:"note,omitempty"`
	Created UTCDateTime `json:"created,omitempty"`
	Author  *Author     `json:"author,omitempty"`

	extensions map[string]json.RawMessage
}

type Author struct {
	Type string `json:"@type,omitempty"`
	Name string `json:"name,omitempty"`
	URI  string `json:"uri,omitempty"`

	extensions map[string]json.RawMessage
}

type PersonalInfo struct {
	Type   string `json:"@type,omitempty"`
	Kind   string `json:"kind,omitempty"`
	Value  string `json:"value,omitempty"`
	Level  string `json:"level,omitempty"`
	ListAs uint64 `json:"listAs,omitempty"`
	Label  string `json:"label,omitempty"`

	extensions map[string]json.RawMessage
}
