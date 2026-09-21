package email

// Email property name constants (RFC 8621 §4.1).
const (
	PropID                    = "id"
	PropBlobID                = "blobId"
	PropThreadID              = "threadId"
	PropMailboxIDs            = "mailboxIds"
	PropKeywords              = "keywords"
	PropSize                  = "size"
	PropReceivedAt            = "receivedAt"
	PropHeaders               = "headers"
	PropMessageID             = "messageId"
	PropInReplyTo             = "inReplyTo"
	PropReferences            = "references"
	PropSender                = "sender"
	PropFrom                  = "from"
	PropTo                    = "to"
	PropCC                    = "cc"
	PropBCC                   = "bcc"
	PropReplyTo               = "replyTo"
	PropSubject               = "subject"
	PropSentAt                = "sentAt"
	PropBodyStructure         = "bodyStructure"
	PropBodyValues            = "bodyValues"
	PropTextBody              = "textBody"
	PropHTMLBody              = "htmlBody"
	PropAttachments           = "attachments"
	PropHasAttachment         = "hasAttachment"
	PropPreview               = "preview"
	PropSMIMEStatus           = "smimeStatus"
	PropSMIMEStatusAtDelivery = "smimeStatusAtDelivery"
	PropSMIMEErrors           = "smimeErrors"
	PropSMIMEVerifiedAt       = "smimeVerifiedAt"
)

// DefaultListProps is a sensible default properties set for Email/get list views.
var DefaultListProps = []string{
	PropID,
	PropThreadID,
	PropMailboxIDs,
	PropKeywords,
	PropReceivedAt,
	PropFrom,
	PropSubject,
	PropPreview,
	PropHasAttachment,
}
