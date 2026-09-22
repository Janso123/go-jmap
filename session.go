package jmap

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

type Session struct {
	// Capabilities specifies the capabililities the server has.
	Capabilities map[URI]Capability `json:"-"`

	RawCapabilities map[URI]jsontext.Value `json:"capabilities"`

	Accounts map[ID]Account `json:"accounts"`

	// PrimaryAccounts maps a Capability to the primary account associated
	// with it
	PrimaryAccounts map[URI]ID `json:"primaryAccounts"`

	// The username associated with the given credentials
	Username string `json:"username"`

	// The URL to use for JMAP API requests.
	APIURL string `json:"apiUrl"`

	// The URL endpoint to use when downloading files
	DownloadURL string `json:"downloadUrl"`

	// The URL endpoint to use when uploading files
	UploadURL string `json:"uploadUrl"`

	// The URL to connect to for push events
	EventSourceURL string `json:"eventSourceUrl"`

	// A string representing the state of this object on the server
	State string `json:"state"`
}

type session Session

func (s *Session) UnmarshalJSON(data []byte) error {
	return jsonv2.Unmarshal(data, s)
}

func (s *Session) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	raw := (*session)(s)
	if err := jsonv2.UnmarshalDecode(dec, &raw); err != nil {
		return err
	}

	s.Capabilities = make(map[URI]Capability)
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
		s.Capabilities[key] = newCap
	})
	if decodeErr != nil {
		return decodeErr
	}

	for id, acc := range s.Accounts {
		acc.ID = string(id)
		s.Accounts[id] = acc
	}

	return nil
}
