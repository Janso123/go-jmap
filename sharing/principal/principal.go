package principal

import "git.sr.ht/~rockorager/go-jmap"

func init() {
	jmap.RegisterMethod("Principal/get", newGetResponse)
	jmap.RegisterMethod("Principal/getAvailability", newGetAvailabilityResponse)
	jmap.RegisterMethod("Principal/changes", newChangesResponse)
	jmap.RegisterMethod("Principal/set", newSetResponse)
	jmap.RegisterMethod("Principal/query", newQueryResponse)
	jmap.RegisterMethod("Principal/queryChanges", newQueryChangesResponse)
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
	ID jmap.ID `json:"id,omitempty"`

	Type PrincipalType `json:"type,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	Email string `json:"email,omitempty"`

	TimeZone string `json:"timeZone,omitempty"`

	Capabilities map[jmap.URI]jmap.Patch `json:"capabilities,omitempty"`

	Accounts map[jmap.ID]jmap.Account `json:"accounts,omitempty"`
}
