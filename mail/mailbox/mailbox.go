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

	// ParentID has no omit: nil marshals as JSON null (top-level mailbox).
	ParentID *jmap.ID `json:"parentId"`

	Role Role `json:"role,omitzero"`

	SortOrder uint64 `json:"sortOrder,omitzero"`

	TotalEmails uint64 `json:"totalEmails,omitzero"`

	UnreadEmails uint64 `json:"unreadEmails,omitzero"`

	TotalThreads uint64 `json:"totalThreads,omitzero"`

	UnreadThreads uint64 `json:"unreadThreads,omitzero"`

	Rights *Rights `json:"myRights,omitzero"`

	IsSubscribed *bool `json:"isSubscribed,omitzero"`
}

func (Mailbox) JMAPType() string { return "Mailbox" }

func (Mailbox) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

// Access Control Lists (ACLs)
type Rights struct {
	MayReadItems bool `json:"mayReadItems,omitzero"`

	MayAddItems bool `json:"mayAddItems,omitzero"`

	MayRemoveItems bool `json:"mayRemoveItems,omitzero"`

	MaySetSeen bool `json:"maySetSeen,omitzero"`

	MaySetKeywords bool `json:"maySetKeywords,omitzero"`

	MayCreateChild bool `json:"mayCreateChild,omitzero"`

	MayRename bool `json:"mayRename,omitzero"`

	MayDelete bool `json:"mayDelete,omitzero"`

	MaySubmit bool `json:"maySubmit,omitzero"`
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
