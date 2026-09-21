package subscription

import (
	"github.com/Janso123/go-jmap"
)

func init() {
	jmap.RegisterObject[PushSubscription](jmap.MethodGet | jmap.MethodSet)
}

// Server side push notification
// https://www.rfc-editor.org/rfc/rfc8620.html#section-7.2
type PushSubscription struct {
	ID jmap.ID `json:"id,omitzero"`

	DeviceClientID string `json:"deviceClientId,omitzero"`

	URL string `json:"url,omitzero"`

	Keys *Key `json:"keys,omitzero"`

	VerificationCode string `json:"verificationCode,omitzero"`

	Expires *jmap.UTCDate `json:"expires,omitzero"`

	Types []string `json:"types,omitzero"`
}

func (PushSubscription) JMAPType() string { return "PushSubscription" }

func (PushSubscription) Requires() []jmap.URI { return nil }

// A Push Subscription Encryption key. This key must be a P-256 ECDH key
type Key struct {
	// The public key, base64 encoded
	Public string `json:"p256dh"`
	// The authentication secret, base64 encoded
	Auth string `json:"auth"`
}
