package searchsnippet

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
)

// Get search snippet details
// https://www.rfc-editor.org/rfc/rfc8621.html#section-5.1
type Get struct {
	Account jmap.ID `json:"accountId,omitzero"`

	Filter jmap.Filter `json:"filter,omitzero"`

	EmailIDs []jmap.ID `json:"emailIds,omitzero"`

	ReferenceIDs *jmap.ResultReference `json:"#emailIds,omitzero"`
}

func (m *Get) Name() string { return "SearchSnippet/get" }

func (m *Get) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

type GetResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	List []*SearchSnippet `json:"list,omitzero"`

	NotFound []jmap.ID `json:"notFound,omitzero"`
}

func newGetResponse() jmap.MethodResponse { return &GetResponse{} }
