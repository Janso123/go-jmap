package vacationresponse

import (
	"github.com/Janso123/go-jmap"
)

const URI jmap.URI = "urn:ietf:params:jmap:vacationresponse"

func init() {
	jmap.RegisterCapability(&Capability{})
	jmap.RegisterObject[VacationResponse](jmap.MethodGet | jmap.MethodSet)
}

// The VacationResponse capability is an empty object
type Capability struct{}

func (m *Capability) URI() jmap.URI { return URI }

func (m *Capability) New() jmap.Capability { return &Capability{} }

// Automatic reply when a message is delivered to the mail store
// https://www.rfc-editor.org/rfc/rfc8621.html#section-8
type VacationResponse struct {
	ID string `json:"id,omitzero"`

	IsEnabled *bool `json:"isEnabled,omitzero"`

	FromDate *jmap.UTCDate `json:"fromDate,omitzero"`

	ToDate *jmap.UTCDate `json:"toDate,omitzero"`

	Subject *string `json:"subject,omitzero"`

	TextBody *string `json:"textBody,omitzero"`

	HTMLBody *string `json:"htmlBody,omitzero"`
}

func (VacationResponse) JMAPType() string { return "VacationResponse" }

func (VacationResponse) Requires() []jmap.URI {
	return []jmap.URI{URI}
}
