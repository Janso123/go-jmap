package email

import (
	"github.com/Janso123/go-jmap"
)

// Parse binary blobs as RFC5322 messages
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.9
type Parse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	BlobIDs []jmap.ID `json:"blobIds,omitzero"`

	Properties []string `json:"properties,omitzero"`

	BodyProperties []string `json:"bodyProperties,omitzero"`

	FetchTextBodyValues bool `json:"fetchTextBodyValues,omitzero"`

	FetchHTMLBodyValues bool `json:"fetchHTMLBodyValues,omitzero"`

	FetchAllBodyValues bool `json:"fetchAllBodyValues,omitzero"`

	// MaxBodyValueBytes is optional. Nil omits the field (server default).
	// A pointer to 0 is sent so the server fetches zero bytes (RFC 8621 §4.9).
	MaxBodyValueBytes *jmap.UnsignedInt `json:"maxBodyValueBytes,omitzero"`
}

func (m *Parse) Name() string { return "Email/parse" }

func (m *Parse) Requires() []jmap.URI {
	return mailRequires(propertiesNeedSMIME(m.Properties))
}

type ParseResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	// Parsed, NotParsable and NotFound are JSON null when empty
	// (RFC 8621 §4.9).
	Parsed jmap.Optional[map[jmap.ID]*Email] `json:"parsed,omitzero"`

	NotParsable jmap.Optional[[]jmap.ID] `json:"notParsable,omitzero"`

	NotFound jmap.Optional[[]jmap.ID] `json:"notFound,omitzero"`
}

func newParseResponse() jmap.MethodResponse { return &ParseResponse{} }
