package email

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
)

// Import email from binary blobs
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.8
type Import struct {
	Account jmap.ID `json:"accountId,omitzero"`

	IfInState string `json:"ifInState,omitzero"`

	Emails map[string]*EmailImport `json:"emails,omitzero"`
}

func (m *Import) Name() string { return "Email/import" }

func (m *Import) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

type EmailImport struct {
	BlobID jmap.ID `json:"blobId,omitzero"`

	MailboxIDs map[jmap.ID]bool `json:"mailboxIds,omitzero"`

	Keywords map[string]bool `json:"keywords,omitzero"`

	ReceivedAt *jmap.UTCDate `json:"receivedAt,omitzero"`
}

type ImportResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	OldState string `json:"oldState,omitzero"`

	NewState string `json:"newState,omitzero"`

	Created map[jmap.ID]*Email `json:"created,omitzero"`

	NotCreated map[jmap.ID]*jmap.SetError `json:"notCreated,omitzero"`
}

func newImportResponse() jmap.MethodResponse { return &ImportResponse{} }
