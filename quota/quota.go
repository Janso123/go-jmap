package quota

import "git.sr.ht/~rockorager/go-jmap"

// urn:ietf:params:jmap:quota represents support for the Quota data type and
// associated API methods.
const URI jmap.URI = "urn:ietf:params:jmap:quota"

const (
	// The Quota event type.
	QuotaEvent jmap.EventType = "Quota"
)

func init() {
	jmap.RegisterCapability(&Capability{})
}

// Capability broadcasts support for quota methods.
type Capability struct{}

func (c *Capability) URI() jmap.URI { return URI }

func (c *Capability) New() jmap.Capability { return &Capability{} }

// Quota describes quota usage and limits for a resource.
type Quota struct {
	ID string `json:"id,omitempty"`

	ResourceType string `json:"resourceType,omitempty"`

	Used uint64 `json:"used,omitempty"`

	HardLimit uint64 `json:"hardLimit,omitempty"`

	Scope string `json:"scope,omitempty"`

	Name string `json:"name,omitempty"`

	Types []string `json:"types,omitempty"`

	WarnLimit uint64 `json:"warnLimit,omitempty"`

	SoftLimit uint64 `json:"softLimit,omitempty"`

	Description string `json:"description,omitempty"`
}
