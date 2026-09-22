package mdn

import (
	"github.com/Janso123/go-jmap"
)

// Parse blobs as messages in the style of RFC5322 to get MDN objects
// https://www.rfc-editor.org/rfc/rfc9007.html#section-2.2
type Parse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	BlobIDs []jmap.ID `json:"blobIds,omitzero"`
}

func (m *Parse) Name() string { return "MDN/parse" }

func (m *Parse) Requires() []jmap.URI { return []jmap.URI{URI} }

type ParseResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	Parsed jmap.Optional[map[jmap.ID]*MDN] `json:"parsed,omitzero"`

	NotParsable jmap.Optional[[]jmap.ID] `json:"notParsable,omitzero"`

	NotFound jmap.Optional[[]jmap.ID] `json:"notFound,omitzero"`
}

func newParseResponse() jmap.MethodResponse { return &ParseResponse{} }
