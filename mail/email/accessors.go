package email

import "github.com/Janso123/go-jmap"

// MailboxIDList returns mailbox IDs where the membership value is true.
func (e *Email) MailboxIDList() []jmap.ID {
	if len(e.MailboxIDs) == 0 {
		return nil
	}
	ids := make([]jmap.ID, 0, len(e.MailboxIDs))
	for id, ok := range e.MailboxIDs {
		if ok {
			ids = append(ids, id)
		}
	}
	return ids
}

// HasKeyword reports whether keywords[kw] is true.
func (e *Email) HasKeyword(kw string) bool {
	return e.Keywords[kw]
}

// IsSeen reports whether the $seen keyword is set.
func (e *Email) IsSeen() bool {
	return e.HasKeyword(KeywordSeen)
}

// BodyValue returns the bodyValues entry for partID.
func (e *Email) BodyValue(partID string) (*BodyValue, bool) {
	if e.BodyValues == nil {
		return nil, false
	}
	v, ok := e.BodyValues[partID]
	return v, ok
}
