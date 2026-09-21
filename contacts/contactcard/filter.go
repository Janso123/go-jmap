package contactcard

import "github.com/Janso123/go-jmap"

// FilterCondition defines ContactCard query filters.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.3.1
type FilterCondition struct {
	jmap.FilterBase `json:"-"`

	InAddressBook jmap.ID `json:"inAddressBook,omitzero"`

	UID string `json:"uid,omitzero"`

	HasMember string `json:"hasMember,omitzero"`

	Kind string `json:"kind,omitzero"`

	CreatedBefore *jmap.UTCDate `json:"createdBefore,omitzero"`

	CreatedAfter *jmap.UTCDate `json:"createdAfter,omitzero"`

	UpdatedBefore *jmap.UTCDate `json:"updatedBefore,omitzero"`

	UpdatedAfter *jmap.UTCDate `json:"updatedAfter,omitzero"`

	Text string `json:"text,omitzero"`

	Name string `json:"name,omitzero"`

	NameGiven string `json:"name/given,omitzero"`

	NameSurname string `json:"name/surname,omitzero"`

	NameSurname2 string `json:"name/surname2,omitzero"`

	Nickname string `json:"nickname,omitzero"`

	Organization string `json:"organization,omitzero"`

	Email string `json:"email,omitzero"`

	Phone string `json:"phone,omitzero"`

	OnlineService string `json:"onlineService,omitzero"`

	Address string `json:"address,omitzero"`

	Note string `json:"note,omitzero"`
}

func And(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.And(conds...) }
func Or(conds ...jmap.Filter) *jmap.FilterOperator  { return jmap.Or(conds...) }
func Not(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.Not(conds...) }
