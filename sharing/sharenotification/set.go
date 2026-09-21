package sharenotification

import "github.com/Janso123/go-jmap"

// Set destroys share notifications (RFC 9670 §3.3: create/update are not permitted).
type Set struct {
	Account          jmap.ID               `json:"accountId,omitzero"`
	IfInState        string                `json:"ifInState,omitzero"`
	Destroy          []jmap.ID             `json:"destroy,omitzero"`
	ReferenceDestroy *jmap.ResultReference `json:"#destroy,omitzero"`
}

func (m *Set) Name() string { return "ShareNotification/set" }

func (m *Set) Requires() []jmap.URI {
	var zero ShareNotification
	return zero.Requires()
}

// SetResponse is the result of ShareNotification/set.
type SetResponse = jmap.SetResponse[ShareNotification]
