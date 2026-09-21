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

	MaxBodyValueBytes uint64 `json:"maxBodyValueBytes,omitzero"`
}

// GetResponse is the result of Email/get.
type GetResponse = jmap.GetResponse[Email]
