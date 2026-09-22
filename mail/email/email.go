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

	// Size is server-set: nil omits it so a create does not send 0.
	Size *jmap.UnsignedInt `json:"size,omitzero"`

	ReceivedAt *jmap.UTCDate `json:"receivedAt,omitzero"`

	Headers []*Header `json:"headers,omitzero"`

	// The header-derived properties below are JSON null when the message has
	// no such header (RFC 8621 §4.1.2).
	MessageID jmap.Optional[[]string] `json:"messageId,omitzero"`

	InReplyTo jmap.Optional[[]string] `json:"inReplyTo,omitzero"`

	References jmap.Optional[[]string] `json:"references,omitzero"`

	Sender jmap.Optional[[]*mail.Address] `json:"sender,omitzero"`

	From jmap.Optional[[]*mail.Address] `json:"from,omitzero"`

	To jmap.Optional[[]*mail.Address] `json:"to,omitzero"`

	CC jmap.Optional[[]*mail.Address] `json:"cc,omitzero"`

	BCC jmap.Optional[[]*mail.Address] `json:"bcc,omitzero"`

	ReplyTo jmap.Optional[[]*mail.Address] `json:"replyTo,omitzero"`

	Subject jmap.Optional[string] `json:"subject,omitzero"`

	SentAt jmap.Optional[jmap.Date] `json:"sentAt,omitzero"`

	BodyStructure *BodyPart `json:"bodyStructure,omitzero"`

	BodyValues map[string]*BodyValue `json:"bodyValues,omitzero"`

	TextBody []*BodyPart `json:"textBody,omitzero"`

	HTMLBody []*BodyPart `json:"htmlBody,omitzero"`

	Attachments []*BodyPart `json:"attachments,omitzero"`

	// HasAttachment is server-set: nil omits it, *false stays on the wire.
	HasAttachment *bool `json:"hasAttachment,omitzero"`

	Preview string `json:"preview,omitzero"`

	SMIMEStatus jmap.Optional[string] `json:"smimeStatus,omitzero"`

	SMIMEStatusAtDelivery jmap.Optional[string] `json:"smimeStatusAtDelivery,omitzero"`

	SMIMEErrors jmap.Optional[[]string] `json:"smimeErrors,omitzero"`

	SMIMEVerifiedAt jmap.Optional[jmap.UTCDate] `json:"smimeVerifiedAt,omitzero"`

	// Extra holds header:* and other unrecognized properties (Approach A embed).
	Extra map[string]jsontext.Value `json:",embed"`
}

func (Email) JMAPType() string { return "Email" }

func (Email) JMAPCreatable() {}

func (Email) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

type AddressGroup struct {
	// Name is JSON null for an ungrouped list of addresses.
	Name jmap.Optional[string] `json:"name,omitzero"`

	Addresses []*mail.Address `json:"addresses,omitzero"`
}

type Header struct {
	Name string `json:"name,omitzero"`

	// Value has no omitzero: an empty header value must stay on the wire.
	Value string `json:"value"`
}

type BodyPart struct {
	// PartID is JSON null for a part with no body content (e.g. multipart).
	PartID jmap.Optional[string] `json:"partId,omitzero"`

	BlobID jmap.Optional[jmap.ID] `json:"blobId,omitzero"`

	// Size is server-set: nil omits it so a create does not send 0.
	Size *jmap.UnsignedInt `json:"size,omitzero"`

	Headers []*Header `json:"headers,omitzero"`

	Name jmap.Optional[string] `json:"name,omitzero"`

	Type string `json:"type,omitzero"`

	Charset jmap.Optional[string] `json:"charset,omitzero"`

	Disposition jmap.Optional[string] `json:"disposition,omitzero"`

	CID jmap.Optional[string] `json:"cid,omitzero"`

	Language jmap.Optional[[]string] `json:"language,omitzero"`

	Location jmap.Optional[string] `json:"location,omitzero"`

	SubParts jmap.Optional[[]*BodyPart] `json:"subParts,omitzero"`

	// Extra holds header:* properties on body parts.
	Extra map[string]jsontext.Value `json:",embed"`
}

type BodyValue struct {
	Value string `json:"value"`

	// IsEncodingProblem has no omitzero: it is a mandatory Boolean.
	IsEncodingProblem bool `json:"isEncodingProblem"`

	// IsTruncated has no omitzero: servers/clients often need explicit false.
	IsTruncated bool `json:"isTruncated"`
}
