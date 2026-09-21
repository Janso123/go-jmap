package jscontact

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

type UTCDateTime string

type PatchObject map[string]any

type Card struct {
	Type    string `json:"@type,omitzero"`
	Version string `json:"version,omitzero"`
	UID     string `json:"uid,omitzero"`

	Created  UTCDateTime     `json:"created,omitzero"`
	Kind     string          `json:"kind,omitzero"`
	Language string          `json:"language,omitzero"`
	Members  map[string]bool `json:"members,omitzero"`
	ProdID   string          `json:"prodId,omitzero"`

	RelatedTo map[string]*Relation `json:"relatedTo,omitzero"`
	Updated   UTCDateTime          `json:"updated,omitzero"`

	Name          *Name                    `json:"name,omitzero"`
	Nicknames     map[string]*Nickname     `json:"nicknames,omitzero"`
	Organizations map[string]*Organization `json:"organizations,omitzero"`
	SpeakToAs     *SpeakToAs               `json:"speakToAs,omitzero"`
	Titles        map[string]*Title        `json:"titles,omitzero"`

	Emails              map[string]*EmailAddress      `json:"emails,omitzero"`
	OnlineServices      map[string]*OnlineService     `json:"onlineServices,omitzero"`
	Phones              map[string]*Phone             `json:"phones,omitzero"`
	PreferredLanguages  map[string]*LanguagePref      `json:"preferredLanguages,omitzero"`
	Calendars           map[string]*Calendar          `json:"calendars,omitzero"`
	SchedulingAddresses map[string]*SchedulingAddress `json:"schedulingAddresses,omitzero"`
	Addresses           map[string]*Address           `json:"addresses,omitzero"`
	CryptoKeys          map[string]*CryptoKey         `json:"cryptoKeys,omitzero"`
	Directories         map[string]*Directory         `json:"directories,omitzero"`
	Links               map[string]*Link              `json:"links,omitzero"`
	Media               map[string]*Media             `json:"media,omitzero"`

	Localizations map[string]PatchObject `json:"localizations,omitzero"`

	Anniversaries map[string]*Anniversary  `json:"anniversaries,omitzero"`
	Keywords      map[string]bool          `json:"keywords,omitzero"`
	Notes         map[string]*Note         `json:"notes,omitzero"`
	PersonalInfo  map[string]*PersonalInfo `json:"personalInfo,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Relation struct {
	Type     string          `json:"@type,omitzero"`
	Relation map[string]bool `json:"relation,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Name struct {
	Type             string            `json:"@type,omitzero"`
	Components       []*NameComponent  `json:"components,omitzero"`
	IsOrdered        bool              `json:"isOrdered,omitzero"`
	DefaultSeparator string            `json:"defaultSeparator,omitzero"`
	Full             string            `json:"full,omitzero"`
	SortAs           map[string]string `json:"sortAs,omitzero"`
	PhoneticScript   string            `json:"phoneticScript,omitzero"`
	PhoneticSystem   string            `json:"phoneticSystem,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type NameComponent struct {
	Type     string `json:"@type,omitzero"`
	Value    string `json:"value,omitzero"`
	Kind     string `json:"kind,omitzero"`
	Phonetic string `json:"phonetic,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Nickname struct {
	Type     string          `json:"@type,omitzero"`
	Name     string          `json:"name,omitzero"`
	Contexts map[string]bool `json:"contexts,omitzero"`
	Pref     uint64          `json:"pref,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Organization struct {
	Type     string          `json:"@type,omitzero"`
	Name     string          `json:"name,omitzero"`
	Units    []*OrgUnit      `json:"units,omitzero"`
	SortAs   string          `json:"sortAs,omitzero"`
	Contexts map[string]bool `json:"contexts,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type OrgUnit struct {
	Type   string `json:"@type,omitzero"`
	Name   string `json:"name,omitzero"`
	SortAs string `json:"sortAs,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type SpeakToAs struct {
	Type              string               `json:"@type,omitzero"`
	GrammaticalGender string               `json:"grammaticalGender,omitzero"`
	Pronouns          map[string]*Pronouns `json:"pronouns,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Pronouns struct {
	Type     string          `json:"@type,omitzero"`
	Pronouns string          `json:"pronouns,omitzero"`
	Contexts map[string]bool `json:"contexts,omitzero"`
	Pref     uint64          `json:"pref,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Title struct {
	Type           string `json:"@type,omitzero"`
	Name           string `json:"name,omitzero"`
	Kind           string `json:"kind,omitzero"`
	OrganizationID string `json:"organizationId,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type EmailAddress struct {
	Type     string          `json:"@type,omitzero"`
	Address  string          `json:"address,omitzero"`
	Contexts map[string]bool `json:"contexts,omitzero"`
	Pref     uint64          `json:"pref,omitzero"`
	Label    string          `json:"label,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type OnlineService struct {
	Type     string          `json:"@type,omitzero"`
	Service  string          `json:"service,omitzero"`
	URI      string          `json:"uri,omitzero"`
	User     string          `json:"user,omitzero"`
	Contexts map[string]bool `json:"contexts,omitzero"`
	Pref     uint64          `json:"pref,omitzero"`
	Label    string          `json:"label,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Phone struct {
	Type     string          `json:"@type,omitzero"`
	Number   string          `json:"number,omitzero"`
	Features map[string]bool `json:"features,omitzero"`
	Contexts map[string]bool `json:"contexts,omitzero"`
	Pref     uint64          `json:"pref,omitzero"`
	Label    string          `json:"label,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type LanguagePref struct {
	Type     string          `json:"@type,omitzero"`
	Language string          `json:"language,omitzero"`
	Contexts map[string]bool `json:"contexts,omitzero"`
	Pref     uint64          `json:"pref,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Resource struct {
	Type      string          `json:"@type,omitzero"`
	Kind      string          `json:"kind,omitzero"`
	URI       string          `json:"uri,omitzero"`
	MediaType string          `json:"mediaType,omitzero"`
	Contexts  map[string]bool `json:"contexts,omitzero"`
	Pref      uint64          `json:"pref,omitzero"`
	Label     string          `json:"label,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Calendar struct {
	Resource

	Extra map[string]jsontext.Value `json:",embed"`
}

type SchedulingAddress struct {
	Type     string          `json:"@type,omitzero"`
	URI      string          `json:"uri,omitzero"`
	Contexts map[string]bool `json:"contexts,omitzero"`
	Pref     uint64          `json:"pref,omitzero"`
	Label    string          `json:"label,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Address struct {
	Type             string              `json:"@type,omitzero"`
	Components       []*AddressComponent `json:"components,omitzero"`
	IsOrdered        bool                `json:"isOrdered,omitzero"`
	CountryCode      string              `json:"countryCode,omitzero"`
	Coordinates      string              `json:"coordinates,omitzero"`
	TimeZone         string              `json:"timeZone,omitzero"`
	Contexts         map[string]bool     `json:"contexts,omitzero"`
	Full             string              `json:"full,omitzero"`
	DefaultSeparator string              `json:"defaultSeparator,omitzero"`
	Pref             uint64              `json:"pref,omitzero"`
	PhoneticScript   string              `json:"phoneticScript,omitzero"`
	PhoneticSystem   string              `json:"phoneticSystem,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type AddressComponent struct {
	Type     string `json:"@type,omitzero"`
	Value    string `json:"value,omitzero"`
	Kind     string `json:"kind,omitzero"`
	Phonetic string `json:"phonetic,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type CryptoKey struct {
	Resource

	Extra map[string]jsontext.Value `json:",embed"`
}

type Directory struct {
	Resource
	ListAs uint64 `json:"listAs,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Link struct {
	Resource

	Extra map[string]jsontext.Value `json:",embed"`
}

type Media struct {
	Resource

	Extra map[string]jsontext.Value `json:",embed"`
}

type Anniversary struct {
	Type  string           `json:"@type,omitzero"`
	Kind  string           `json:"kind,omitzero"`
	Date  *AnniversaryDate `json:"date,omitzero"`
	Place *Address         `json:"place,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type AnniversaryDate struct {
	PartialDate *PartialDate
	Timestamp   *Timestamp
}

func (d AnniversaryDate) MarshalJSON() ([]byte, error) {
	switch {
	case d.Timestamp != nil:
		return jsonv2.Marshal(d.Timestamp)
	case d.PartialDate != nil:
		return jsonv2.Marshal(d.PartialDate)
	default:
		return []byte("null"), nil
	}
}

func (d *AnniversaryDate) UnmarshalJSON(data []byte) error {
	var probe struct {
		Type string `json:"@type"`
	}
	if err := jsonv2.Unmarshal(data, &probe); err != nil {
		return err
	}
	if probe.Type == "Timestamp" {
		var ts Timestamp
		if err := jsonv2.Unmarshal(data, &ts); err != nil {
			return err
		}
		d.Timestamp = &ts
		d.PartialDate = nil
		return nil
	}
	var pd PartialDate
	if err := jsonv2.Unmarshal(data, &pd); err != nil {
		return err
	}
	d.PartialDate = &pd
	d.Timestamp = nil
	return nil
}

type PartialDate struct {
	Type          string `json:"@type,omitzero"`
	Year          uint64 `json:"year,omitzero"`
	Month         uint64 `json:"month,omitzero"`
	Day           uint64 `json:"day,omitzero"`
	CalendarScale string `json:"calendarScale,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Timestamp struct {
	Type string      `json:"@type,omitzero"`
	UTC  UTCDateTime `json:"utc,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Note struct {
	Type    string      `json:"@type,omitzero"`
	Note    string      `json:"note,omitzero"`
	Created UTCDateTime `json:"created,omitzero"`
	Author  *Author     `json:"author,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type Author struct {
	Type string `json:"@type,omitzero"`
	Name string `json:"name,omitzero"`
	URI  string `json:"uri,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}

type PersonalInfo struct {
	Type   string `json:"@type,omitzero"`
	Kind   string `json:"kind,omitzero"`
	Value  string `json:"value,omitzero"`
	Level  string `json:"level,omitzero"`
	ListAs uint64 `json:"listAs,omitzero"`
	Label  string `json:"label,omitzero"`

	Extra map[string]jsontext.Value `json:",embed"`
}
