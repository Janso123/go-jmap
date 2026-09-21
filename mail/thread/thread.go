package thread

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
)

func init() {
	jmap.RegisterObject[Thread](jmap.MethodGet | jmap.MethodChanges)
}

// See RFC8621, Section 3.
type Thread struct {
	ID jmap.ID `json:"id,omitzero"`

	EmailIDs []jmap.ID `json:"emailIds,omitzero"`
}

func (Thread) JMAPType() string { return "Thread" }

func (Thread) Requires() []jmap.URI { return []jmap.URI{mail.URI} }
