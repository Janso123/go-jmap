package blob

import "github.com/Janso123/go-jmap"

// Lookup reverse references to blobs across JMAP data types.
// https://www.rfc-editor.org/rfc/rfc9404.html#section-4.3
type Lookup struct {
	Account jmap.ID `json:"accountId,omitzero"`

	TypeNames []string `json:"typeNames,omitzero"`

	IDs []jmap.ID `json:"ids,omitzero"`
}

func (m *Lookup) Name() string { return "Blob/lookup" }

func (m *Lookup) Requires() []jmap.URI { return []jmap.URI{URI} }

type LookupResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	List []*BlobInfo `json:"list,omitzero"`

	NotFound []jmap.ID `json:"notFound,omitzero"`
}

type BlobInfo struct {
	ID jmap.ID `json:"id,omitzero"`

	MatchedIDs map[string][]jmap.ID `json:"matchedIds,omitzero"`
}

func newLookupResponse() jmap.MethodResponse { return &LookupResponse{} }
