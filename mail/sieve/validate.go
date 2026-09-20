package sieve

import "git.sr.ht/~rockorager/go-jmap"

// Validate checks Sieve script validity without storing it on the server.
// https://www.rfc-editor.org/rfc/rfc9661.html#section-2.6
type Validate struct {
	Account jmap.ID `json:"accountId,omitempty"`

	BlobID jmap.ID `json:"blobId,omitempty"`
}

func (m *Validate) Name() string { return "SieveScript/validate" }

func (m *Validate) Requires() []jmap.URI { return []jmap.URI{URI} }

type ValidateResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	Error *jmap.SetError `json:"error,omitempty"`
}

func newValidateResponse() jmap.MethodResponse { return &ValidateResponse{} }
