package blob

import "github.com/Janso123/go-jmap"

// Lookup reverse references to blobs across JMAP data types.
// https://www.rfc-editor.org/rfc/rfc9404.html#section-4.3
type Lookup struct {
	Account jmap.ID `json:"accountId,omitempty"`

	TypeNames []string `json:"typeNames,omitempty"`

	IDs []jmap.ID `json:"ids,omitempty"`
}

func (m *Lookup) Name() string { return "Blob/lookup" }

func (m *Lookup) Requires() []jmap.URI { return []jmap.URI{URI} }

type LookupResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	List []*BlobInfo `json:"list,omitempty"`

	NotFound []jmap.ID `json:"notFound,omitempty"`
}

type BlobInfo struct {
	ID jmap.ID `json:"id,omitempty"`

	MatchedIDs map[string][]jmap.ID `json:"matchedIds,omitempty"`
}

func newLookupResponse() jmap.MethodResponse { return &LookupResponse{} }
