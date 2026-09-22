package jmap

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

// An account is a collection of data the authenticated user has access to
//
// See RFC 8620 section 1.6.2 for details.
type Account struct {
	// The ID of the account
	ID string `json:"-"`

	// A user-friendly string to show when presenting content from this
	// account, e.g. the email address representing the owner of the account.
	Name string `json:"name"`

	// True if this account belongs to the authenticated user
	IsPersonal bool `json:"isPersonal"`

	IsReadOnly bool `json:"isReadOnly"`

	// The set of capability URIs for the methods supported in this account.
	Capabilities map[URI]Capability `json:"-"`

	// The raw JSON of accountCapabilities
	RawCapabilities map[URI]jsontext.Value `json:"accountCapabilities"`
}

type account Account

func (a *Account) UnmarshalJSON(data []byte) error {
	return jsonv2.Unmarshal(data, a)
}

func (a *Account) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	raw := (*account)(a)
	if err := jsonv2.UnmarshalDecode(dec, &raw); err != nil {
		return err
	}

	a.Capabilities = make(map[URI]Capability)
	var decodeErr error
	rangeCapabilities(func(key URI, cap Capability) {
		if decodeErr != nil {
			return
		}
		rawCap, ok := raw.RawCapabilities[key]
		if !ok {
			return
		}
		newCap := cap.New()
		if err := jsonv2.Unmarshal(rawCap, newCap); err != nil {
			decodeErr = err
			return
		}
		a.Capabilities[key] = newCap
	})
	if decodeErr != nil {
		return decodeErr
	}

	return nil
}
