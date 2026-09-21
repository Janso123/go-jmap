package searchsnippet

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
)

// Get search snippet details
// https://www.rfc-editor.org/rfc/rfc8621.html#section-5.1
type Get struct {
	Account jmap.ID `json:"accountId,omitempty"`

	Filter any `json:"filter,omitempty"`

	EmailIDs []jmap.ID `json:"emailIds,omitempty"`

	ReferenceIDs *jmap.ResultReference `json:"#emailIds,omitempty"`
}

func (m *Get) Name() string { return "SearchSnippet/get" }

func (m *Get) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

type GetResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	List []*SearchSnippet `json:"list,omitempty"`

	NotFound []jmap.ID `json:"notFound,omitempty"`
}

func newGetResponse() jmap.MethodResponse { return &GetResponse{} }
