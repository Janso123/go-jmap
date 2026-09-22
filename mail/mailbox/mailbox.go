package mailbox

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
)

func init() {
	jmap.RegisterObject[Mailbox](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodQuery |
			jmap.MethodQueryChanges |
			jmap.MethodSet,
	)
	// Mailbox/changes includes updatedProperties (RFC 8621 §2.2); override kit factory.
	jmap.RegisterMethod("Mailbox/changes", func() jmap.MethodResponse { return &ChangesResponse{} })
}

// Named set of Email objects. Can be viewed as a folder or a label.
// An email must be part of at least one Mailbox.
// https://www.rfc-editor.org/rfc/rfc8621.html#section-2
type Mailbox struct {
	ID jmap.ID `json:"id,omitzero"`

	Name string `json:"name,omitzero"`

	// ParentID is JSON null for a top-level mailbox and omitted when unset.
	ParentID jmap.Optional[jmap.ID] `json:"parentId,omitzero"`

	// Role is JSON null for a mailbox with no special-use role.
	Role jmap.Optional[Role] `json:"role,omitzero"`

	// SortOrder is server-set: nil omits it so a create does not send 0.
	SortOrder *jmap.UnsignedInt `json:"sortOrder,omitzero"`

	TotalEmails *jmap.UnsignedInt `json:"totalEmails,omitzero"`

	UnreadEmails *jmap.UnsignedInt `json:"unreadEmails,omitzero"`

	TotalThreads *jmap.UnsignedInt `json:"totalThreads,omitzero"`

	UnreadThreads *jmap.UnsignedInt `json:"unreadThreads,omitzero"`

	Rights *Rights `json:"myRights,omitzero"`

	IsSubscribed *bool `json:"isSubscribed,omitzero"`
}

func (Mailbox) JMAPType() string { return "Mailbox" }

func (Mailbox) JMAPCreatable() {}

func (Mailbox) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

// Access Control Lists (ACLs)
//
// Every right is a mandatory Boolean; false must stay on the wire. The whole
// object sits behind *Rights,omitzero, so a create omits it entirely.
type Rights struct {
	MayReadItems bool `json:"mayReadItems"`

	MayAddItems bool `json:"mayAddItems"`

	MayRemoveItems bool `json:"mayRemoveItems"`

	MaySetSeen bool `json:"maySetSeen"`

	MaySetKeywords bool `json:"maySetKeywords"`

	MayCreateChild bool `json:"mayCreateChild"`

	MayRename bool `json:"mayRename"`

	MayDelete bool `json:"mayDelete"`

	MaySubmit bool `json:"maySubmit"`
}

// Identifies Mailboxes that have a particular common purpose (e.g., the
// “inbox”), regardless of the name property (which may be localised).
// Values are from the IANA IMAP Mailbox Name Attributes registry (SPECIAL-USE).
type Role string

const (
	RoleAll Role = "all"

	RoleArchive Role = "archive"

	RoleDrafts Role = "drafts"

	RoleFlagged Role = "flagged"

	RoleImportant Role = "important"

	RoleInbox Role = "inbox"

	RoleJunk Role = "junk"

	RoleSent Role = "sent"

	RoleTrash Role = "trash"
)
