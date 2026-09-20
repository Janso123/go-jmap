package sharenotification

import (
	"time"

	"git.sr.ht/~rockorager/go-jmap"
)

func init() {
	jmap.RegisterMethod("ShareNotification/get", newGetResponse)
	jmap.RegisterMethod("ShareNotification/changes", newChangesResponse)
	jmap.RegisterMethod("ShareNotification/set", newSetResponse)
	jmap.RegisterMethod("ShareNotification/query", newQueryResponse)
	jmap.RegisterMethod("ShareNotification/queryChanges", newQueryChangesResponse)
}

const (
	// ShareNotificationEvent is the ShareNotification event type.
	ShareNotificationEvent jmap.EventType = "ShareNotification"
)

// Entity describes who made a sharing change.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3
type Entity struct {
	Name string `json:"name,omitempty"`

	Email *string `json:"email,omitempty"`

	PrincipalID *jmap.ID `json:"principalId,omitempty"`
}

// ShareNotification records a change to a user's access rights on a shared object.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3
type ShareNotification struct {
	ID jmap.ID `json:"id,omitempty"`

	Created *time.Time `json:"created,omitempty"`

	ChangedBy *Entity `json:"changedBy,omitempty"`

	ObjectAccountID jmap.ID `json:"objectAccountId,omitempty"`

	ObjectType string `json:"objectType,omitempty"`

	ObjectID jmap.ID `json:"objectId,omitempty"`

	OldRights map[string]bool `json:"oldRights,omitempty"`

	NewRights map[string]bool `json:"newRights,omitempty"`

	Name string `json:"name,omitempty"`
}
