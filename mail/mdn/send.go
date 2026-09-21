package mdn

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
)

// Sends an RFC5322 message from an MDN object
// https://www.rfc-editor.org/rfc/rfc9007.html#section-2.1
type Send struct {
	Account jmap.ID `json:"accountId,omitzero"`

	IdentityID jmap.ID `json:"identityId,omitzero"`

	Send map[jmap.ID]*MDN `json:"send,omitzero"`

	OnSuccessUpdateEmail map[jmap.ID]*jmap.Patch `json:"onSuccessUpdateEmail,omitzero"`
}

func (m *Send) Name() string { return "MDN/send" }

func (m *Send) Requires() []jmap.URI { return []jmap.URI{mail.URI, URI} }

type SendResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	Sent map[jmap.ID]*MDN `json:"sent,omitzero"`

	NotSent map[jmap.ID]*jmap.SetError `json:"notSent,omitzero"`
}

func newSendResponse() jmap.MethodResponse { return &SendResponse{} }
