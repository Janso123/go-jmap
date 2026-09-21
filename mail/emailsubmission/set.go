package emailsubmission

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
)

// Create, delete or modify an email submission
// https://www.rfc-editor.org/rfc/rfc8621.html#section-7.5
type Set struct {
	jmap.Set[EmailSubmission]

	OnSuccessUpdateEmail map[jmap.ID]jmap.Patch `json:"onSuccessUpdateEmail,omitzero"`

	OnSuccessDestroyEmail []jmap.ID `json:"onSuccessDestroyEmail,omitzero"`
}

func (s *Set) Requires() []jmap.URI {
	if len(s.OnSuccessUpdateEmail) > 0 || len(s.OnSuccessDestroyEmail) > 0 {
		return []jmap.URI{URI, mail.URI}
	}
	return []jmap.URI{URI}
}

// SetResponse is the result of EmailSubmission/set.
type SetResponse = jmap.SetResponse[EmailSubmission]
