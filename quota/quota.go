package quota

import "github.com/Janso123/go-jmap"

// urn:ietf:params:jmap:quota represents support for the Quota data type and
// associated API methods.
const URI jmap.URI = "urn:ietf:params:jmap:quota"

const (
	// The Quota event type.
	QuotaEvent jmap.EventType = "Quota"
)

func init() {
	jmap.RegisterCapability(&Capability{})
	jmap.RegisterObject[Quota](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodQuery |
			jmap.MethodQueryChanges,
	)
	// Quota/changes includes updatedProperties; override kit factory.
	jmap.RegisterMethod("Quota/changes", func() jmap.MethodResponse { return &ChangesResponse{} })
}

// Capability broadcasts support for quota methods.
type Capability struct{}

func (c *Capability) URI() jmap.URI { return URI }

func (c *Capability) New() jmap.Capability { return &Capability{} }

// Quota describes quota usage and limits for a resource.
type Quota struct {
	ID jmap.ID `json:"id,omitzero"`

	ResourceType string `json:"resourceType,omitzero"`

	Used uint64 `json:"used,omitzero"`

	HardLimit uint64 `json:"hardLimit,omitzero"`

	Scope string `json:"scope,omitzero"`

	Name string `json:"name,omitzero"`

	Types []string `json:"types,omitzero"`

	WarnLimit uint64 `json:"warnLimit,omitzero"`

	SoftLimit uint64 `json:"softLimit,omitzero"`

	Description string `json:"description,omitzero"`
}

func (Quota) JMAPType() string { return "Quota" }

func (Quota) Requires() []jmap.URI { return []jmap.URI{URI} }
