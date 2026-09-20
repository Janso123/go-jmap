package quota

import "git.sr.ht/~rockorager/go-jmap"

func init() {
	jmap.RegisterMethod("Quota/changes", newChangesResponse)
}

// Changes gets quota changes for the whole account.
type Changes struct {
	Account jmap.ID `json:"accountId,omitempty"`

	SinceState string `json:"sinceState,omitempty"`

	MaxChanges uint64 `json:"maxChanges,omitempty"`
}

func (m *Changes) Name() string { return "Quota/changes" }

func (m *Changes) Requires() []jmap.URI { return []jmap.URI{URI} }

type ChangesResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	OldState string `json:"oldState,omitempty"`

	NewState string `json:"newState,omitempty"`

	HasMoreChanges bool `json:"hasMoreChanges,omitempty"`

	Created []jmap.ID `json:"created,omitempty"`

	Updated []jmap.ID `json:"updated,omitempty"`

	Destroyed []jmap.ID `json:"destroyed,omitempty"`

	UpdatedProperties []string `json:"updatedProperties,omitempty"`
}

func newChangesResponse() jmap.MethodResponse { return &ChangesResponse{} }
