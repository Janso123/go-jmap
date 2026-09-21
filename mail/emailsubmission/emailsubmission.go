package emailsubmission

import (
	"github.com/Janso123/go-jmap"
)

const URI jmap.URI = "urn:ietf:params:jmap:submission"

func init() {
	jmap.RegisterCapability(&Capability{})
	jmap.RegisterObject[EmailSubmission](
		jmap.MethodGet | jmap.MethodChanges | jmap.MethodQuery |
			jmap.MethodQueryChanges | jmap.MethodSet,
	)
}

// The EmailSubmission Capability
type Capability struct {
	// The maximum number of seconds the server supports for delayed
	// sending. A value of 0 indicates delayed sending is not supported
	MaxDelayedSend uint64 `json:"maxDelayedSend,omitzero"`

	// The set of SMTP submission extensions supported by the server, which
	// the client may use when creating an EmailSubmission object (see
	// Section 7). Each key in the object is the ehlo-name, and the value is
	// a list of ehlo-args.
	SubmissionExtensions map[string][]string `json:"submissionExtensions,omitzero"`
}

func (m *Capability) URI() jmap.URI { return URI }

func (m *Capability) New() jmap.Capability { return &Capability{} }

// Submission of an Email for delivery to one or more recipients.
// https://www.rfc-editor.org/rfc/rfc8621.html#section-7
type EmailSubmission struct {
	ID jmap.ID `json:"id,omitzero"`

	IdentityID jmap.ID `json:"identityId,omitzero"`

	EmailID jmap.ID `json:"emailId,omitzero"`

	ThreadID jmap.ID `json:"threadId,omitzero"`

	Envelope *Envelope `json:"envelope,omitzero"`

	SendAt *jmap.UTCDate `json:"sendAt,omitzero"`

	UndoStatus UndoStatus `json:"undoStatus,omitzero"`

	DeliveryStatus map[string]*DeliveryStatus `json:"deliveryStatus,omitzero"`

	DSNBlobIDs []jmap.ID `json:"dsnBlobIds,omitzero"`

	MDNBlobIDs []jmap.ID `json:"mdnBlobIds,omitzero"`
}

func (EmailSubmission) JMAPType() string { return "EmailSubmission" }

func (EmailSubmission) Requires() []jmap.URI {
	return []jmap.URI{URI}
}

// UndoStatus is the undoability of an EmailSubmission.
type UndoStatus string

const (
	UndoPending  UndoStatus = "pending"
	UndoFinal    UndoStatus = "final"
	UndoCanceled UndoStatus = "canceled"
)

type Envelope struct {
	// The email address to use as the return address in the SMTP submission
	MailFrom *Address `json:"mailFrom,omitzero"`

	// The email address to send the message to
	RcptTo []*Address `json:"rcptTo,omitzero"`
}

type Address struct {
	// The email address
	Email string `json:"email"`

	// Parameters to send with the email submission, if any SMTP extensions
	// are used. A nil value means the parameter is present with no value.
	Parameters map[string]*string `json:"parameters,omitzero"`
}

// Delivered is the delivery outcome for a recipient.
type Delivered string

const (
	DeliveredQueued  Delivered = "queued"
	DeliveredYes     Delivered = "yes"
	DeliveredNo      Delivered = "no"
	DeliveredUnknown Delivered = "unknown"
)

// Displayed is whether the message has been displayed by the recipient.
type Displayed string

const (
	DisplayedUnknown Displayed = "unknown"
	DisplayedYes     Displayed = "yes"
)

type DeliveryStatus struct {
	// The SMTP reply returned for the recipient
	SMTPReply string `json:"smtpReply,omitzero"`

	// Represents whether the message has been successfully delivered to the
	// recipient.
	Delivered Delivered `json:"delivered,omitzero"`

	// Whether the message has been displayed by the recipient.
	Displayed Displayed `json:"displayed,omitzero"`
}
