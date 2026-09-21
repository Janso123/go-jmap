package email

import (
	"encoding/json/jsontext"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
)

func init() {
	jmap.RegisterCapability(&smimeVerify{})
	jmap.RegisterObject[Email](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodQuery |
			jmap.MethodQueryChanges |
			jmap.MethodSet |
			jmap.MethodCopy,
	)
	jmap.RegisterMethod("Email/import", newImportResponse)
	jmap.RegisterMethod("Email/parse", newParseResponse)
}

// Representation of an RFC5322 message
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4
type Email struct {
	ID jmap.ID `json:"id,omitzero"`

	BlobID jmap.ID `json:"blobId,omitzero"`

	ThreadID jmap.ID `json:"threadId,omitzero"`

	MailboxIDs map[jmap.ID]bool `json:"mailboxIds,omitzero"`

	Keywords map[string]bool `json:"keywords,omitzero"`

	Size uint64 `json:"size,omitzero"`

	ReceivedAt *jmap.UTCDate `json:"receivedAt,omitzero"`

	Headers []*Header `json:"headers,omitzero"`

	MessageID []string `json:"messageId,omitzero"`

	InReplyTo []string `json:"inReplyTo,omitzero"`

	References []string `json:"references,omitzero"`

	Sender []*mail.Address `json:"sender,omitzero"`

	From []*mail.Address `json:"from,omitzero"`

	To []*mail.Address `json:"to,omitzero"`

	CC []*mail.Address `json:"cc,omitzero"`

	BCC []*mail.Address `json:"bcc,omitzero"`

	ReplyTo []*mail.Address `json:"replyTo,omitzero"`

	Subject string `json:"subject,omitzero"`

	SentAt *jmap.Date `json:"sentAt,omitzero"`

	BodyStructure *BodyPart `json:"bodyStructure,omitzero"`

	BodyValues map[string]*BodyValue `json:"bodyValues,omitzero"`

	TextBody []*BodyPart `json:"textBody,omitzero"`

	HTMLBody []*BodyPart `json:"htmlBody,omitzero"`

	Attachments []*BodyPart `json:"attachments,omitzero"`

	HasAttachment bool `json:"hasAttachment,omitzero"`

	Preview string `json:"preview,omitzero"`

	SMIMEStatus string `json:"smimeStatus,omitzero"`

	SMIMEStatusAtDelivery string `json:"smimeStatusAtDelivery,omitzero"`

	SMIMEErrors []string `json:"smimeErrors,omitzero"`

	SMIMEVerifiedAt *jmap.UTCDate `json:"smimeVerifiedAt,omitzero"`

	// Extra holds header:* and other unrecognized properties (Approach A embed).
	Extra map[string]jsontext.Value `json:",embed"`
}

func (Email) JMAPType() string { return "Email" }

func (Email) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

type AddressGroup struct {
	Name string `json:"name,omitzero"`

	Addresses []*mail.Address `json:"addresses,omitzero"`
}

type Header struct {
	Name string `json:"name,omitzero"`

	Value string `json:"value,omitzero"`
}

type BodyPart struct {
	PartID string `json:"partId,omitzero"`

	BlobID jmap.ID `json:"blobId,omitzero"`

	Size uint64 `json:"size,omitzero"`

	Headers []*Header `json:"headers,omitzero"`

	Name string `json:"name,omitzero"`

	Type string `json:"type,omitzero"`

	Charset string `json:"charset,omitzero"`

	Disposition string `json:"disposition,omitzero"`

	CID string `json:"cid,omitzero"`

	Language []string `json:"language,omitzero"`

	Location string `json:"location,omitzero"`

	SubParts []*BodyPart `json:"subParts,omitzero"`

	// Extra holds header:* properties on body parts.
	Extra map[string]jsontext.Value `json:",embed"`
}

type BodyValue struct {
	Value string `json:"value"`

	IsEncodingProblem bool `json:"isEncodingProblem,omitzero"`

	// IsTruncated has no omitzero: servers/clients often need explicit false.
	IsTruncated bool `json:"isTruncated"`
}
