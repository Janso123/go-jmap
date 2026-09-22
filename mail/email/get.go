package email

import "github.com/Janso123/go-jmap"

// Get email details
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.2
type Get struct {
	jmap.Get[Email]

	BodyProperties []string `json:"bodyProperties,omitzero"`

	FetchTextBodyValues bool `json:"fetchTextBodyValues,omitzero"`

	FetchHTMLBodyValues bool `json:"fetchHTMLBodyValues,omitzero"`

	FetchAllBodyValues bool `json:"fetchAllBodyValues,omitzero"`

	// MaxBodyValueBytes is optional. Nil omits the field (server default).
	// A pointer to 0 is sent so the server fetches zero bytes (RFC 8621 §4.2).
	MaxBodyValueBytes *jmap.UnsignedInt `json:"maxBodyValueBytes,omitzero"`
}

func (g *Get) Requires() []jmap.URI {
	props, _ := g.Properties.Value()
	return mailRequires(propertiesNeedSMIME(props))
}

// GetResponse is the result of Email/get.
type GetResponse = jmap.GetResponse[Email]
