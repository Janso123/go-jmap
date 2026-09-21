package email

import "github.com/Janso123/go-jmap"

type Filter interface {
	implementsFilter()
}

type FilterOperator struct {
	Operator jmap.Operator `json:"operator,omitzero"`

	Conditions []Filter `json:"conditions,omitzero"`
}

func (fo *FilterOperator) implementsFilter() {}

// Email query condition that can be compounded with FilterOperator
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.4.1
type FilterCondition struct {
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

	HasAttachment bool `json:"hasAttachment,omitzero"`

	Text string `json:"text,omitzero"`

	From string `json:"from,omitzero"`

	To string `json:"to,omitzero"`

	Cc string `json:"cc,omitzero"`

	Bcc string `json:"bcc,omitzero"`

	Subject string `json:"subject,omitzero"`

	Body string `json:"body,omitzero"`

	Header []string `json:"header,omitzero"`

	HasSMIME bool `json:"hasSmime,omitzero"`

	HasVerifiedSMIME bool `json:"hasVerifiedSmime,omitzero"`

	HasVerifiedSMIMEAtDelivery bool `json:"hasVerifiedSmimeAtDelivery,omitzero"`
}

func (fc *FilterCondition) implementsFilter() {}

// And returns a FilterOperator matching when all conditions match.
func And(conds ...Filter) *FilterOperator {
	return &FilterOperator{Operator: jmap.OperatorAND, Conditions: conds}
}

// Or returns a FilterOperator matching when any condition matches.
func Or(conds ...Filter) *FilterOperator {
	return &FilterOperator{Operator: jmap.OperatorOR, Conditions: conds}
}

// Not returns a FilterOperator matching when none of the conditions match.
func Not(conds ...Filter) *FilterOperator {
	return &FilterOperator{Operator: jmap.OperatorNOT, Conditions: conds}
}

// InMailbox filters emails in the given mailbox.
func InMailbox(id jmap.ID) *FilterCondition {
	return &FilterCondition{InMailbox: id}
}

// InMailboxesOtherThan filters emails in any mailbox other than the given ones.
func InMailboxesOtherThan(ids ...jmap.ID) *FilterCondition {
	return &FilterCondition{InMailboxOtherThan: ids}
}

// HasKeyword returns a filter for emails that have the given keyword.
// Note: Email.HasKeyword is the object accessor with the same name.
func HasKeyword(kw string) *FilterCondition {
	return &FilterCondition{HasKeyword: kw}
}

// NotKeyword returns a filter for emails that do not have the given keyword.
func NotKeyword(kw string) *FilterCondition {
	return &FilterCondition{NotKeyword: kw}
}

// Subject returns a filter matching subject substring.
func Subject(s string) *FilterCondition {
	return &FilterCondition{Subject: s}
}

// From returns a filter matching From substring.
func From(s string) *FilterCondition {
	return &FilterCondition{From: s}
}

// To returns a filter matching To substring.
func To(s string) *FilterCondition {
	return &FilterCondition{To: s}
}

// Text returns a full-text search filter.
func Text(s string) *FilterCondition {
	return &FilterCondition{Text: s}
}

// Body returns a filter matching body substring.
func Body(s string) *FilterCondition {
	return &FilterCondition{Body: s}
}

// Before returns a filter for emails received before d.
func Before(d jmap.UTCDate) *FilterCondition {
	return &FilterCondition{Before: &d}
}

// After returns a filter for emails received after d.
func After(d jmap.UTCDate) *FilterCondition {
	return &FilterCondition{After: &d}
}

// WithAttachment returns a filter for emails that have attachments.
func WithAttachment() *FilterCondition {
	return &FilterCondition{HasAttachment: true}
}
