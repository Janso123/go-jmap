package searchsnippet

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/email"
)

// Get search snippet details
// https://www.rfc-editor.org/rfc/rfc8621.html#section-5.1
type Get struct {
	Account jmap.ID `json:"accountId,omitzero"`

	Filter jmap.Filter `json:"filter,omitzero"`

	EmailIDs []jmap.ID `json:"emailIds,omitzero"`

	ReferenceIDs *jmap.ResultReference `json:"#emailIds,omitzero"`
}

type getAlias Get

func (m *Get) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var s struct {
		getAlias
		Filter jsontext.Value `json:"filter"`
	}
	if err := jsonv2.UnmarshalDecode(dec, &s); err != nil {
		return err
	}
	f, err := jmap.UnmarshalFilter[email.FilterCondition](s.Filter)
	if err != nil {
		return err
	}
	*m = Get(s.getAlias)
	m.Filter = f
	return nil
}

func (m *Get) Name() string { return "SearchSnippet/get" }

func (m *Get) Requires() []jmap.URI {
	uris := []jmap.URI{mail.URI}
	if email.FilterNeedsSMIME(m.Filter) {
		uris = append(uris, email.SMIMEVerify)
	}
	return uris
}

type GetResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	List []*SearchSnippet `json:"list,omitzero"`

	// NotFound is JSON null when every requested email was found
	// (RFC 8621 §5.1).
	NotFound jmap.Optional[[]jmap.ID] `json:"notFound,omitzero"`
}

func newGetResponse() jmap.MethodResponse { return &GetResponse{} }
