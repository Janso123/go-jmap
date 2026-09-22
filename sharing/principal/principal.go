package principal

import (
	"encoding/json/jsontext"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/sharing"
)

func init() {
	jmap.RegisterObject[Principal](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodQuery |
			jmap.MethodQueryChanges |
			jmap.MethodSet,
	)
	// Manual method (not part of the standard kit).
	jmap.RegisterMethod("Principal/getAvailability", newGetAvailabilityResponse)
}

const (
	// PrincipalEvent is the Principal event type.
	PrincipalEvent jmap.EventType = "Principal"
)

// PrincipalType identifies the kind of principal the object represents.
type PrincipalType string

const (
	TypeIndividual PrincipalType = "individual"

	TypeGroup PrincipalType = "group"

	TypeResource PrincipalType = "resource"

	TypeLocation PrincipalType = "location"

	TypeOther PrincipalType = "other"
)

// Principal represents an individual, team, location, or other shared entity.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2
type Principal struct {
	ID jmap.ID `json:"id,omitzero"`

	Type PrincipalType `json:"type,omitzero"`

	Name string `json:"name,omitzero"`

	Description jmap.Optional[string] `json:"description,omitzero"`

	Email jmap.Optional[string] `json:"email,omitzero"`

	TimeZone jmap.Optional[string] `json:"timeZone,omitzero"`

	Capabilities map[jmap.URI]jsontext.Value `json:"capabilities,omitzero"`

	Accounts jmap.Optional[map[jmap.ID]jmap.Account] `json:"accounts,omitzero"`
}

func (Principal) JMAPType() string { return "Principal" }

func (Principal) JMAPCreatable() {}

func (Principal) Requires() []jmap.URI { return []jmap.URI{sharing.URI} }
