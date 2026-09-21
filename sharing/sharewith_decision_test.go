package sharing_test

import "testing"

// TestShareWithNormativeDecision records the RFC 9670 §4 spike for client
// Mailbox.shareWith (Calendar/AddressBook already define shareWith in their
// own specs; this decision is only about adding it to Mailbox).
//
// RFC 9670 §4 is a framework for shareable data types: types that opt in MUST
// define isSubscribed, myRights, and shareWith. It does not amend RFC 8621
// Mailbox. Mailbox already has myRights/isSubscribed from RFC 8621 and has no
// shareWith property there, so shareWith is not normative for client Mailbox
// in v1-alpha.
func TestShareWithNormativeDecision(t *testing.T) {
	// If RFC 9670 §4 requires client Mailbox.shareWith, this test unmarshals a fixture with shareWith.
	// If not normative, t.Skip("RFC 9670 §4 shareWith not required on client Mailbox for v1-alpha")
	t.Skip("RFC 9670 §4 shareWith not required on client Mailbox for v1-alpha")
}
