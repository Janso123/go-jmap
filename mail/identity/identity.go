package identity

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/emailsubmission"
)

func init() {
	jmap.RegisterObject[Identity](jmap.MethodGet | jmap.MethodChanges | jmap.MethodSet)
}

// Information about an email address or domain the user may send from
// https://www.rfc-editor.org/rfc/rfc8621.html#section-6
type Identity struct {
	ID jmap.ID `json:"id,omitzero"`

	Name string `json:"name,omitzero"`

	Email string `json:"email,omitzero"`

	ReplyTo jmap.Optional[[]*mail.Address] `json:"replyTo,omitzero"`

	Bcc jmap.Optional[[]*mail.Address] `json:"bcc,omitzero"`

	TextSignature string `json:"textSignature,omitzero"`

	HTMLSignature string `json:"htmlSignature,omitzero"`

	MayDelete *bool `json:"mayDelete,omitzero"`
}

func (Identity) JMAPType() string { return "Identity" }

func (Identity) JMAPCreatable() {}

func (Identity) Requires() []jmap.URI {
	return []jmap.URI{emailsubmission.URI}
}
