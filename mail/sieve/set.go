package sieve

import "github.com/Janso123/go-jmap"

// Set creates, updates, destroys, activates, and deactivates Sieve scripts.
// https://www.rfc-editor.org/rfc/rfc9661.html#section-2.4
type Set struct {
	jmap.Set[SieveScript]

	OnSuccessActivateScript jmap.ID `json:"onSuccessActivateScript,omitzero"`

	OnSuccessDeactivateScript bool `json:"onSuccessDeactivateScript,omitzero"`
}

// SetResponse is the result of SieveScript/set.
type SetResponse = jmap.SetResponse[SieveScript]
