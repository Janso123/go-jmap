package vapid

import "github.com/Janso123/go-jmap"

const URI jmap.URI = "urn:ietf:params:jmap:webpush-vapid"

func init() { jmap.RegisterCapability(&Capability{}) }

type Capability struct {
	ApplicationServerKey string `json:"applicationServerKey"`
}

func (c *Capability) URI() jmap.URI        { return URI }
func (c *Capability) New() jmap.Capability { return &Capability{} }
