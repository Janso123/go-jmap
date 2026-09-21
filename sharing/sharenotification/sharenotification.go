package sharenotification

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/sharing"
)

func init() {
	jmap.RegisterObject[ShareNotification](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodQuery |
			jmap.MethodQueryChanges |
			jmap.MethodSet,
	)
	// ShareNotification/changes includes updatedProperties; override kit factory.
	jmap.RegisterMethod("ShareNotification/changes", func() jmap.MethodResponse { return &ChangesResponse{} })
}

const (
	// ShareNotificationEvent is the ShareNotification event type.
	ShareNotificationEvent jmap.EventType = "ShareNotification"
)

// Entity describes who made a sharing change.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3
type Entity struct {
	Name string `json:"name,omitzero"`

	Email *string `json:"email,omitzero"`

	PrincipalID *jmap.ID `json:"principalId,omitzero"`
}

// ShareNotification records a change to a user's access rights on a shared object.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3
type ShareNotification struct {
	ID jmap.ID `json:"id,omitzero"`

	Created *jmap.UTCDate `json:"created,omitzero"`

	ChangedBy *Entity `json:"changedBy,omitzero"`

	ObjectAccountID jmap.ID `json:"objectAccountId,omitzero"`

	ObjectType string `json:"objectType,omitzero"`

	ObjectID jmap.ID `json:"objectId,omitzero"`

	OldRights map[string]bool `json:"oldRights,omitzero"`

	NewRights map[string]bool `json:"newRights,omitzero"`

	Name string `json:"name,omitzero"`
}

func (ShareNotification) JMAPType() string { return "ShareNotification" }

func (ShareNotification) Requires() []jmap.URI { return []jmap.URI{sharing.URI} }
