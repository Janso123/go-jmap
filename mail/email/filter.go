package email

import "github.com/Janso123/go-jmap"

type FilterCondition struct {
	jmap.FilterBase `json:"-"`

	InMailbox jmap.ID `json:"inMailbox,omitzero"`

	InMailboxOtherThan []jmap.ID `json:"inMailboxOtherThan,omitzero"`

	Before *jmap.UTCDate `json:"before,omitzero"`

	After *jmap.UTCDate `json:"after,omitzero"`

	MinSize uint64 `json:"minSize,omitzero"`

	MaxSize uint64 `json:"maxSize,omitzero"`

	AllInThreadHaveKeyword string `json:"allInThreadHaveKeyword,omitzero"`

	SomeInThreadHaveKeyword string `json:"someInThreadHaveKeyword,omitzero"`

	NoneInThreadHaveKeyword string `json:"noneInThreadHaveKeyword,omitzero"`

	HasKeyword string `json:"hasKeyword,omitzero"`

	NotKeyword string `json:"notKeyword,omitzero"`

	HasAttachment *bool `json:"hasAttachment,omitzero"`

	Text string `json:"text,omitzero"`

	From string `json:"from,omitzero"`

	To string `json:"to,omitzero"`

	Cc string `json:"cc,omitzero"`

	Bcc string `json:"bcc,omitzero"`

	Subject string `json:"subject,omitzero"`

	Body string `json:"body,omitzero"`

	Header []string `json:"header,omitzero"`

	HasSMIME *bool `json:"hasSmime,omitzero"`

	HasVerifiedSMIME *bool `json:"hasVerifiedSmime,omitzero"`

	HasVerifiedSMIMEAtDelivery *bool `json:"hasVerifiedSmimeAtDelivery,omitzero"`
}

func And(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.And(conds...) }
func Or(conds ...jmap.Filter) *jmap.FilterOperator  { return jmap.Or(conds...) }
func Not(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.Not(conds...) }

func InMailbox(id jmap.ID) *FilterCondition {
	return &FilterCondition{InMailbox: id}
}

func InMailboxesOtherThan(ids ...jmap.ID) *FilterCondition {
	return &FilterCondition{InMailboxOtherThan: ids}
}

func HasKeyword(kw string) *FilterCondition {
	return &FilterCondition{HasKeyword: kw}
}

func NotKeyword(kw string) *FilterCondition {
	return &FilterCondition{NotKeyword: kw}
}

func Subject(s string) *FilterCondition { return &FilterCondition{Subject: s} }
func From(s string) *FilterCondition    { return &FilterCondition{From: s} }
func To(s string) *FilterCondition      { return &FilterCondition{To: s} }
func Text(s string) *FilterCondition    { return &FilterCondition{Text: s} }
func Body(s string) *FilterCondition    { return &FilterCondition{Body: s} }

func Before(d jmap.UTCDate) *FilterCondition { return &FilterCondition{Before: &d} }
func After(d jmap.UTCDate) *FilterCondition  { return &FilterCondition{After: &d} }

func WithAttachment() *FilterCondition {
	return &FilterCondition{HasAttachment: jmap.Bool(true)}
}

func WithoutAttachment() *FilterCondition {
	return &FilterCondition{HasAttachment: jmap.Bool(false)}
}
