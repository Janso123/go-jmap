//go:build e2e

package e2e

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/Janso123/go-jmap/mail/emailsubmission"
	"github.com/Janso123/go-jmap/mail/identity"
	"github.com/Janso123/go-jmap/mail/mailbox"
	"github.com/Janso123/go-jmap/mail/thread"
)

func TestMail(t *testing.T) {
	sc := &scenario{name: "Mail"}
	if alice == nil || alice.Client == nil {
		t.Fatal("alice session missing")
	}
	if !hasCap(alice.Client, mail.URI) {
		skipRest(t, sc, mailSteps, "missing "+string(mail.URI))
		return
	}

	mailSubject := "e2e-mail-" + started.Format("20060102T150405Z")
	sendSubject := "e2e-send-" + started.Format("20060102T150405Z")

	aliceID, err := alice.Client.PrimaryAccount(mail.URI)
	if err != nil {
		t.Errorf("Mail primary account: %v", err)
	}
	var bobID jmap.ID
	if bob != nil && bob.Client != nil {
		bobID, err = bob.Client.PrimaryAccount(mail.URI)
		if err != nil {
			t.Errorf("Mail bob primary account: %v", err)
		}
	}

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

	call[*mailbox.QueryResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Mailbox/query", Account: alice.Name, Request: "filter role=inbox",
	}, alice.Client, []jmap.URI{mail.URI}, false, &mailbox.Query{
		Account: aliceID,
		Filter:  &mailbox.FilterCondition{Role: jmap.Some(mailbox.RoleInbox)},
	}, func(resp *mailbox.QueryResponse) (string, error) {
		if resp == nil || !containsID(resp.IDs, inboxID) {
			return "", errString("inbox id missing from query")
		}
		return "ids=" + strconv.Itoa(len(resp.IDs)), nil
	})

	var identityID jmap.ID
	call[*identity.GetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Identity/get", Account: alice.Name, Request: "list",
	}, alice.Client, []jmap.URI{mail.URI}, false, &identity.Get{Account: aliceID}, func(resp *identity.GetResponse) (string, error) {
		if resp == nil || len(resp.List) < 1 {
			return "", errString("no identity")
		}
		identityID = resp.List[0].ID
		return "identities=" + strconv.Itoa(len(resp.List)), nil
	})

	var stateBefore string
	call[*email.GetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Email/get", Account: alice.Name, Request: "ids omitted",
	}, alice.Client, []jmap.URI{mail.URI}, false, &email.Get{Account: aliceID}, func(resp *email.GetResponse) (string, error) {
		if resp == nil || resp.State == "" {
			return "", errString("email state missing")
		}
		stateBefore = resp.State
		return "state=" + stateBefore, nil
	})

	var queryBefore string
	call[*email.QueryResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Email/query", Account: alice.Name, Request: "inMailbox inbox",
	}, alice.Client, []jmap.URI{mail.URI}, false, &email.Query{
		Account: aliceID,
		Filter:  email.InMailbox(inboxID),
	}, func(resp *email.QueryResponse) (string, error) {
		if resp == nil {
			return "", errString("empty email query")
		}
		queryBefore = resp.QueryState
		return "ids=" + strconv.Itoa(len(resp.IDs)), nil
	})

	var emailID, threadID jmap.ID
	call[*email.SetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Email/set", Account: alice.Name, Request: "create subject=" + mailSubject,
	}, alice.Client, []jmap.URI{mail.URI}, false, &email.Set{
		Account: aliceID,
		Create: jmap.Some(map[jmap.ID]*email.Email{
			"m1": newE2EEmail(inboxID, mailSubject, true),
		}),
	}, func(resp *email.SetResponse) (string, error) {
		created, err := createdEmail(resp, "m1")
		if err != nil {
			return "", err
		}
		emailID = created.ID
		threadID = created.ThreadID
		return "id=" + string(emailID), nil
	})

	call[*email.GetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Email/get", Account: alice.Name, Request: "id=" + string(emailID),
	}, alice.Client, []jmap.URI{mail.URI}, false, &email.Get{
		Account: aliceID,
		IDs:     jmap.Some([]jmap.ID{emailID}),
	}, func(resp *email.GetResponse) (string, error) {
		if resp == nil || len(resp.List) == 0 {
			return "", errString("email missing")
		}
		subject, ok := resp.List[0].Subject.Value()
		if !ok || subject != mailSubject {
			return "", errString("subject mismatch")
		}
		return "subject=" + subject, nil
	})

	call[*email.QueryResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Email/query", Account: alice.Name, Request: "inMailbox inbox",
	}, alice.Client, []jmap.URI{mail.URI}, false, &email.Query{
		Account: aliceID,
		Filter:  email.InMailbox(inboxID),
	}, func(resp *email.QueryResponse) (string, error) {
		if resp == nil || !containsID(resp.IDs, emailID) {
			return "", errString("email id missing from query")
		}
		return "ids=" + strconv.Itoa(len(resp.IDs)), nil
	})

	call[*email.ChangesResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Email/changes", Account: alice.Name, Request: "sinceState=" + stateBefore,
	}, alice.Client, []jmap.URI{mail.URI}, false, &email.Changes{
		Account:    aliceID,
		SinceState: stateBefore,
	}, func(resp *email.ChangesResponse) (string, error) {
		if resp == nil || !containsID(resp.Created, emailID) {
			return "", errString("email id not in created")
		}
		return "created=" + string(emailID), nil
	})

	call[*email.QueryChangesResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Email/queryChanges", Account: alice.Name, Request: "sinceQueryState=" + queryBefore,
	}, alice.Client, []jmap.URI{mail.URI}, false, &email.QueryChanges{
		Account:         aliceID,
		SinceQueryState: queryBefore,
		Filter:          email.InMailbox(inboxID),
	}, func(resp *email.QueryChangesResponse) (string, error) {
		if resp == nil {
			return "", errString("empty queryChanges")
		}
		for _, item := range resp.Added {
			if item.ID == emailID {
				return "added=" + string(emailID), nil
			}
		}
		return "", errString("email id not in added")
	})

	call[*thread.GetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Thread/get", Account: alice.Name, Request: "id=" + string(threadID),
	}, alice.Client, []jmap.URI{mail.URI}, false, &thread.Get{
		Account: aliceID,
		IDs:     jmap.Some([]jmap.ID{threadID}),
	}, func(resp *thread.GetResponse) (string, error) {
		if resp == nil || len(resp.List) == 0 {
			return "", errString("thread list empty")
		}
		return "threads=" + strconv.Itoa(len(resp.List)), nil
	})

	var sendID jmap.ID
	call[*email.SetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "Email/set", Account: alice.Name, Request: "create subject=" + sendSubject,
	}, alice.Client, []jmap.URI{mail.URI}, false, &email.Set{
		Account: aliceID,
		Create: jmap.Some(map[jmap.ID]*email.Email{
			"m1": newE2EEmail(inboxID, sendSubject, false),
		}),
	}, func(resp *email.SetResponse) (string, error) {
		created, err := createdEmail(resp, "m1")
		if err != nil {
			return "", err
		}
		sendID = created.ID
		return "id=" + string(sendID), nil
	})

	if !hasCap(alice.Client, emailsubmission.URI) {
		skipRest(t, sc, mailSteps[len(mailSteps)-2:], "missing "+string(emailsubmission.URI))
		if sc.failed {
			t.Fail()
		}
		return
	}

	call[*emailsubmission.SetResponse](t, sc, step{
		RFC: "RFC 8621", Method: "EmailSubmission/set", Account: alice.Name, Request: "emailId=" + string(sendID),
	}, alice.Client, []jmap.URI{emailsubmission.URI}, false, &emailsubmission.Set{
		Account: aliceID,
		Create: jmap.Some(map[jmap.ID]*emailsubmission.EmailSubmission{
			"s1": {
				IdentityID: identityID,
				EmailID:    sendID,
				Envelope: jmap.Some(emailsubmission.Envelope{
					MailFrom: &emailsubmission.Address{Email: alice.Email},
					RcptTo:   []*emailsubmission.Address{{Email: bob.Email}},
				}),
			},
		}),
	}, func(resp *emailsubmission.SetResponse) (string, error) {
		if resp == nil {
			return "", errString("empty submission response")
		}
		if nc, ok := resp.NotCreated.Value(); ok {
			if se := nc["s1"]; se != nil {
				return "", errString(se.Error())
			}
		}
		created, ok := resp.Created["s1"]
		if !ok || created.ID == "" {
			return "", errString("submission was not created")
		}
		return "id=" + string(created.ID), nil
	})

	bobClient := (*jmap.Client)(nil)
	if bob != nil {
		bobClient = bob.Client
	}
	if bobClient == nil {
		t.Errorf("Mail bob session missing")
		sc.failed = true
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for poll := 1; poll <= 30; poll++ {
		last := poll == 30
		resp, ok := call[*email.QueryResponse](t, sc, step{
			RFC: "RFC 8621", Method: "Email/query", Account: "bob", Request: "subject=" + sendSubject,
		}, bobClient, []jmap.URI{mail.URI}, false, &email.Query{
			Account: bobID,
			Filter:  &email.FilterCondition{Subject: sendSubject},
		}, func(resp *email.QueryResponse) (string, error) {
			if resp != nil && len(resp.IDs) > 0 {
				return "ids=" + strconv.Itoa(len(resp.IDs)), nil
			}
			if last {
				return "", errString("bob did not receive the message")
			}
			return "still empty", nil
		})
		if ok && resp != nil && len(resp.IDs) > 0 {
			break
		}
		if last || !ok {
			continue
		}
		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
		}
	}

	if sc.failed {
		t.Fail()
	}
}

var mailSteps = []step{
	{RFC: "RFC 8621", Method: "Mailbox/get", Account: "alice", Request: "role=inbox"},
	{RFC: "RFC 8621", Method: "Mailbox/query", Account: "alice", Request: "filter role=inbox"},
	{RFC: "RFC 8621", Method: "Identity/get", Account: "alice", Request: "list"},
	{RFC: "RFC 8621", Method: "Email/set", Account: "alice", Request: "create"},
	{RFC: "RFC 8621", Method: "Email/get", Account: "alice", Request: "subject"},
	{RFC: "RFC 8621", Method: "Email/query", Account: "alice", Request: "inMailbox inbox"},
	{RFC: "RFC 8621", Method: "Email/changes", Account: "alice", Request: "sinceState"},
	{RFC: "RFC 8621", Method: "Email/queryChanges", Account: "alice", Request: "sinceQueryState"},
	{RFC: "RFC 8621", Method: "Thread/get", Account: "alice", Request: "threadId"},
	{RFC: "RFC 8621", Method: "EmailSubmission/set", Account: "alice", Request: "submit"},
	{RFC: "RFC 8621", Method: "Email/query", Account: "bob", Request: "subject"},
}

func newE2EEmail(inboxID jmap.ID, subject string, draft bool) *email.Email {
	msg := &email.Email{
		MailboxIDs: map[jmap.ID]bool{inboxID: true},
		From:       jmap.Some([]*mail.Address{{Email: alice.Email}}),
		To:         jmap.Some([]*mail.Address{{Email: bob.Email}}),
		Subject:    jmap.Some(subject),
		TextBody:   []*email.BodyPart{{PartID: jmap.Some("1"), Type: "text/plain"}},
		BodyValues: map[string]*email.BodyValue{"1": {Value: "e2e"}},
	}
	if draft {
		msg.Keywords = map[string]bool{"$draft": true}
	}
	return msg
}

func createdEmail(resp *email.SetResponse, key jmap.ID) (email.Email, error) {
	if resp == nil {
		return email.Email{}, errString("empty set response")
	}
	if nc, ok := resp.NotCreated.Value(); ok {
		if se := nc[key]; se != nil {
			return email.Email{}, errString(se.Error())
		}
	}
	created, ok := resp.Created[key]
	if !ok || created.ID == "" {
		return email.Email{}, errString("email was not created")
	}
	return created, nil
}

func containsID(ids []jmap.ID, id jmap.ID) bool {
	for _, candidate := range ids {
		if candidate == id {
			return true
		}
	}
	return false
}
