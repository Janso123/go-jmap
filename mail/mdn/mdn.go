package mdn

import "github.com/Janso123/go-jmap"

const URI jmap.URI = "urn:ietf:params:jmap:mdn"

func init() {
	jmap.RegisterCapability(&Capability{})
	jmap.RegisterMethod("MDN/send", newSendResponse)
	jmap.RegisterMethod("MDN/parse", newParseResponse)
}

// The MDN Capability
type Capability struct{}

func (m *Capability) URI() jmap.URI { return URI }

func (m *Capability) New() jmap.Capability { return &Capability{} }

// A Message Delivery Notification (MDN) object
// https://www.rfc-editor.org/rfc/rfc9007.html#section-2
type MDN struct {
	ForEmailID jmap.Optional[jmap.ID] `json:"forEmailId,omitzero"`

	Subject jmap.Optional[string] `json:"subject,omitzero"`

	TextBody jmap.Optional[string] `json:"textBody,omitzero"`

	IncludeOriginalMessage *bool `json:"includeOriginalMessage,omitzero"`

	ReportingUA jmap.Optional[string] `json:"reportingUA,omitzero"`

	Disposition *Disposition `json:"disposition,omitzero"`

	MDNGateway jmap.Optional[string] `json:"mdnGateway,omitzero"`

	OriginalRecipient jmap.Optional[string] `json:"originalRecipient,omitzero"`

	FinalRecipient jmap.Optional[string] `json:"finalRecipient,omitzero"`

	OriginalMessageID jmap.Optional[string] `json:"originalMessageId,omitzero"`

	Error jmap.Optional[[]string] `json:"error,omitzero"`

	ExtensionFields jmap.Optional[map[string]string] `json:"extensionFields,omitzero"`
}

type ActionMode string

const (
	ActionManual    ActionMode = "manual-action"
	ActionAutomatic ActionMode = "automatic-action"
)

type SendingMode string

const (
	SendingManual    SendingMode = "mdn-sent-manually"
	SendingAutomatic SendingMode = "mdn-sent-automatically"
)

type DispositionType string

const (
	DispositionDeleted    DispositionType = "deleted"
	DispositionDispatched DispositionType = "dispatched"
	DispositionDisplayed  DispositionType = "displayed"
	DispositionProcessed  DispositionType = "processed"
)

type Disposition struct {
	ActionMode ActionMode `json:"actionMode,omitzero"`

	SendingMode SendingMode `json:"sendingMode,omitzero"`

	Type DispositionType `json:"type,omitzero"`
}
