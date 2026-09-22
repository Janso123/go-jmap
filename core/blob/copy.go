package blob

import (
	"github.com/Janso123/go-jmap"
)

// Copy a binary blob from one account to another
// https://www.rfc-editor.org/rfc/rfc8620.html#section-6.3
type Copy struct {
	FromAccount jmap.ID `json:"fromAccountId,omitzero"`

	Account jmap.ID `json:"accountId,omitzero"`

	IDs []jmap.ID `json:"blobIds,omitzero"`
}

func (m *Copy) Name() string { return "Blob/copy" }

func (m *Copy) Requires() []jmap.URI { return []jmap.URI{jmap.CoreURI} }

type CopyResponse struct {
	FromAccount jmap.ID `json:"fromAccountId,omitzero"`

	Account jmap.ID `json:"accountId,omitzero"`

	Copied jmap.Optional[map[jmap.ID]jmap.ID] `json:"copied,omitzero"`

	NotCopied jmap.Optional[map[jmap.ID]*jmap.SetError] `json:"notCopied,omitzero"`
}

func newCopyResponse() jmap.MethodResponse { return &CopyResponse{} }
