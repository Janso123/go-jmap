package sieve

import "github.com/Janso123/go-jmap"

// Get Sieve script details.
// https://www.rfc-editor.org/rfc/rfc9661.html#section-2.3
type Get struct {
	jmap.Get[SieveScript]
}

// GetResponse is the result of SieveScript/get.
type GetResponse = jmap.GetResponse[SieveScript]
