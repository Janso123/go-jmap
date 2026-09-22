//go:build e2e

package e2e

import (
	"context"
	"strconv"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/blob"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/Janso123/go-jmap/mail/mailbox"
)

const blobPayload = "e2e-blob-hello"

func TestBlob(t *testing.T) {
	sc := &scenario{name: "Blob"}
	if alice == nil || alice.Client == nil || bob == nil || bob.Client == nil {
		t.Fatal("account session missing")
	}

	aliceID, err := alice.Client.PrimaryAccount(mail.URI)
	if err != nil {
		t.Errorf("Blob primary account: %v", err)
		sc.failed = true
	}
	bobID, err := bob.Client.PrimaryAccount(mail.URI)
	if err != nil {
		t.Errorf("Blob bob primary account: %v", err)
	}

	subject := "e2e-blob-" + started.Format("20060102T150405Z")

	blobID, _ := uploadBlob(t, sc, step{
		RFC: "RFC 8620", Method: "Blob/upload", Account: alice.Name, Request: "text/plain",
	}, alice.Client, aliceID, blobPayload, "text/plain")

	got, ok := downloadBlob(t, sc, step{
		RFC: "RFC 8620", Method: "Blob/download", Account: alice.Name, Request: "blobId=" + string(blobID),
	}, alice.Client, aliceID, blobID, blobPayload)
	if ok && string(got) != blobPayload {
		sc.failed = true
		t.Errorf("Blob Blob/download: got %q", got)
	}

	call[*blob.GetResponse](t, sc, step{
		RFC: "RFC 9404", Method: "Blob/get", Account: alice.Name, Request: "properties=data:asText",
	}, alice.Client, []jmap.URI{blob.URI}, false, &blob.Get{
		Account:    aliceID,
		IDs:        jmap.Some([]jmap.ID{blobID}),
		Properties: []string{"data:asText"},
	}, func(resp *blob.GetResponse) (string, error) {
		if resp == nil || len(resp.List) == 0 || resp.List[0] == nil {
			return "", errString("blob missing")
		}
		text, hasText := resp.List[0].DataAsText.Value()
		if !hasText || text != blobPayload {
			return "", errString("data:asText mismatch")
		}
		return "data:asText=" + text, nil
	})

	var inboxID jmap.ID
	call[*mailbox.GetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Mailbox/get", Account: alice.Name, Request: "role=inbox",
	}, alice.Client, []jmap.URI{mail.URI}, false, &mailbox.Get{Account: aliceID}, func(resp *mailbox.GetResponse) (string, error) {
		if resp == nil {
			return "", errString("empty mailbox list")
		}
		for _, mb := range resp.List {
			role, ok := mb.Role.Value()
			if ok && role == mailbox.RoleInbox {
				inboxID = mb.ID
				return "inbox=" + string(inboxID), nil
			}
		}
		return "", errString("inbox role missing")
	})

	var emailID jmap.ID
	call[*email.SetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Email/set", Account: alice.Name, Request: "create subject=" + subject,
	}, alice.Client, []jmap.URI{mail.URI}, false, &email.Set{
		Account: aliceID,
		Create: jmap.Some(map[jmap.ID]*email.Email{
			"m1": {
				MailboxIDs: map[jmap.ID]bool{inboxID: true},
				From:       jmap.Some([]*mail.Address{{Email: alice.Email}}),
				To:         jmap.Some([]*mail.Address{{Email: bob.Email}}),
				Subject:    jmap.Some(subject),
				TextBody: []*email.BodyPart{{
					BlobID: jmap.Some(blobID),
					Type:   "text/plain",
				}},
			},
		}),
	}, func(resp *email.SetResponse) (string, error) {
		id, err := createdBlobEmail(resp, "m1")
		if err != nil {
			return "", err
		}
		emailID = id
		return "id=" + string(emailID), nil
	})

	var messageBlob, textBlob jmap.ID
	call[*email.GetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Email/get", Account: alice.Name, Request: "id=" + string(emailID),
	}, alice.Client, []jmap.URI{mail.URI}, false, &email.Get{
		Account:        aliceID,
		IDs:            jmap.Some([]jmap.ID{emailID}),
		Properties:     jmap.Some([]string{"blobId", "textBody"}),
		BodyProperties: []string{"blobId", "type", "partId"},
	}, func(resp *email.GetResponse) (string, error) {
		if resp == nil || len(resp.List) == 0 {
			return "", errString("email missing")
		}
		msg := resp.List[0]
		if msg.BlobID == "" {
			return "", errString("message blob missing")
		}
		if len(msg.TextBody) == 0 || msg.TextBody[0] == nil {
			return "", errString("text body missing")
		}
		partID, hasPart := msg.TextBody[0].BlobID.Value()
		if !hasPart || partID == "" {
			return "", errString("text body blob missing")
		}
		messageBlob = msg.BlobID
		textBlob = partID
		return "message=" + string(messageBlob) + " part=" + string(textBlob), nil
	})

	call[*blob.LookupResponse](t, sc, step{
		RFC: "RFC 9404", Method: "Blob/lookup", Account: alice.Name, Request: "typeNames=Email message",
	}, alice.Client, []jmap.URI{blob.URI}, false, &blob.Lookup{
		Account:   aliceID,
		TypeNames: []string{"Email"},
		IDs:       []jmap.ID{messageBlob},
	}, func(resp *blob.LookupResponse) (string, error) {
		if resp == nil {
			return "", errString("empty lookup")
		}
		for _, info := range resp.List {
			if info != nil && info.ID == messageBlob && containsID(info.MatchedIDs["Email"], emailID) {
				return "emailId=" + string(emailID), nil
			}
		}
		return "", errString("email id missing from matchedIds")
	})

	call[*blob.LookupResponse](t, sc, step{
		RFC: "RFC 9404", Method: "Blob/lookup", Account: alice.Name, Request: "typeNames=Email upload",
	}, alice.Client, []jmap.URI{blob.URI}, false, &blob.Lookup{
		Account:   aliceID,
		TypeNames: []string{"Email"},
		IDs:       []jmap.ID{blobID},
	}, func(resp *blob.LookupResponse) (string, error) {
		if resp == nil {
			return "", errString("empty lookup")
		}
		if containsID(resp.NotFound, blobID) {
			return "", errString("upload id in notFound")
		}
		var info *blob.BlobInfo
		for _, item := range resp.List {
			if item != nil && item.ID == blobID {
				info = item
				break
			}
		}
		if info == nil {
			return "", errString("upload id missing from list")
		}
		n := 0
		for _, ids := range info.MatchedIDs {
			n += len(ids)
		}
		return "matched=" + strconv.Itoa(n), nil
	})

	textGot, textOK := downloadBlob(t, sc, step{
		RFC: "RFC 8620", Method: "Blob/download", Account: alice.Name, Request: "part=" + string(textBlob),
	}, alice.Client, aliceID, textBlob, blobPayload)
	if textOK && string(textGot) != blobPayload {
		sc.failed = true
		t.Errorf("Blob Blob/download: got %q", textGot)
	}

	call[*mailbox.SetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Mailbox/set", Account: alice.Name, Request: "shareWith bob mayReadItems",
	}, alice.Client, []jmap.URI{mail.URI}, false, &mailbox.Set{
		Account: aliceID,
		Update: jmap.Some(map[jmap.ID]jmap.Patch{
			inboxID: {
				"shareWith": map[string]any{
					string(bobID): map[string]bool{"mayReadItems": true},
				},
			},
		}),
	}, func(resp *mailbox.SetResponse) (string, error) {
		if err := rejectMailboxUpdate(resp, inboxID); err != nil {
			return "", err
		}
		return "shared=" + string(bobID), nil
	})

	copiedID := copyBlobFromAlice(t, sc, aliceID, bobID, textBlob)

	bobGot, bobOK := downloadBlob(t, sc, step{
		RFC: "RFC 8620", Method: "Blob/download", Account: bob.Name, Request: "blobId=" + string(copiedID),
	}, bob.Client, bobID, copiedID, blobPayload)
	if bobOK && string(bobGot) != blobPayload {
		sc.failed = true
		t.Errorf("Blob Blob/download: got %q", bobGot)
	}

	if sc.failed {
		t.Fail()
	}
}

func copyBlobFromAlice(t *testing.T, sc *scenario, from, to, source jmap.ID) jmap.ID {
	t.Helper()
	st := step{
		RFC: "RFC 8620", Method: "Blob/copy", Account: bob.Name,
		Request: "from=" + string(from) + " blobId=" + string(source),
	}
	if sc.failed {
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: "an earlier step failed", Result: kindNotRun})
		return ""
	}
	copyReq := blob.Copy{FromAccount: from, Account: to, IDs: []jmap.ID{source}}
	var method jmap.Method = &copyReq
	if bob != nil && bob.Client != nil && hasCap(bob.Client, blob.URI) {
		method = blobCopyMethod{Copy: copyReq}
	}
	resp, err := jmap.Call[*blob.CopyResponse](context.Background(), bob.Client, method)
	if err != nil {
		sc.failed = true
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: err.Error(), Result: kindFail})
		t.Errorf("%s %s: %v", sc.name, st.Method, err)
		return ""
	}
	if resp == nil {
		sc.failed = true
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: "empty copy response", Result: kindFail})
		t.Errorf("%s %s: empty copy response", sc.name, st.Method)
		return ""
	}
	if nc, hasNC := resp.NotCopied.Value(); hasNC && len(nc) > 0 {
		msg := setErrorSummary(nc)
		sc.failed = true
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: msg, Result: kindFail})
		t.Errorf("%s %s: %s", sc.name, st.Method, msg)
		return ""
	}
	copied, hasCopied := resp.Copied.Value()
	newID := copied[source]
	if !hasCopied || newID == "" {
		sc.failed = true
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: "copied id missing", Result: kindFail})
		t.Errorf("%s %s: copied id missing", sc.name, st.Method)
		return ""
	}
	rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: "id=" + string(newID), Result: kindPass})
	return newID
}

func rejectMailboxUpdate(resp *mailbox.SetResponse, id jmap.ID) error {
	if resp == nil {
		return errString("empty set response")
	}
	if nu, ok := resp.NotUpdated.Value(); ok {
		if se := nu[id]; se != nil {
			return errString(setErrorText(se))
		}
	}
	return nil
}

func createdBlobEmail(resp *email.SetResponse, key jmap.ID) (jmap.ID, error) {
	if resp == nil {
		return "", errString("empty set response")
	}
	if nc, ok := resp.NotCreated.Value(); ok {
		if se := nc[key]; se != nil {
			return "", errString(setErrorText(se))
		}
	}
	created, ok := resp.Created[key]
	if !ok || created.ID == "" {
		return "", errString("email was not created")
	}
	return created.ID, nil
}

// blobCopyMethod advertises the blob capability in "using". Stalwart rejects
// Blob/copy when that URI is absent, and blob.Copy.Requires reports only core.
type blobCopyMethod struct {
	blob.Copy
}

func (m blobCopyMethod) Name() string { return "Blob/copy" }

func (m blobCopyMethod) Requires() []jmap.URI {
	return []jmap.URI{jmap.CoreURI, blob.URI}
}
