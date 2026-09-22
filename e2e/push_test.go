//go:build e2e

package e2e

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/push"
	"github.com/Janso123/go-jmap/core/push/websocket"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/Janso123/go-jmap/mail/mailbox"
)

func TestPush(t *testing.T) {
	sc := &scenario{name: "Push"}
	if alice == nil || alice.Client == nil || alice.Client.Session == nil {
		t.Fatal("alice session missing")
	}
	if bob == nil {
		t.Fatal("bob session missing")
	}

	pushSubject := "e2e-push-" + started.Format("20060102T150405Z")
	emailStep := step{
		RFC: "RFC 8621", Method: "Email/set", Account: alice.Name,
		Request: "create subject=" + pushSubject,
	}

	changes := make(chan *jmap.StateChange, 32)
	var listenErr <-chan error
	var transportRFC string

	wsCap, hasWS := alice.Client.Session.Capabilities[websocket.URI].(*websocket.WebSocket)
	switch {
	case hasWS && wsCap.SupportsPush:
		transportRFC = "RFC 8887"
		enableStep := step{
			RFC: transportRFC, Method: "pushEnable", Account: alice.Name,
			Request: "dataTypes=Email",
		}
		ctx := context.Background()
		conn, err := websocket.Dial(ctx, alice.Client)
		if err != nil {
			failPush(t, sc, enableStep, err.Error())
			notRunPush(sc, emailStep, pushStateStep(transportRFC))
			return
		}
		defer conn.Close()
		conn.SetHandler(func(change *jmap.StateChange) {
			changes <- change
		})
		if err := conn.EnablePush(ctx, []jmap.EventType{mail.EmailEvent}, conn.PushState()); err != nil {
			failPush(t, sc, enableStep, err.Error())
			notRunPush(sc, emailStep, pushStateStep(transportRFC))
			return
		}
		recordPush(sc, enableStep, "enabled", kindPass)
	case alice.Client.Session.EventSourceURL != "":
		transportRFC = "RFC 8620"
		enableStep := step{
			RFC: transportRFC, Method: "EventSource", Account: alice.Name,
			Request: "types=Email",
		}
		listenCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errCh := make(chan error, 1)
		listenErr = errCh
		go func() {
			errCh <- (&push.EventSource{
				Client: alice.Client,
				Events: []jmap.EventType{mail.EmailEvent},
				Handler: func(change *jmap.StateChange) {
					changes <- change
				},
			}).Listen(listenCtx)
		}()
		recordPush(sc, enableStep, "listening", kindPass)
	default:
		skipRest(t, sc, []step{
			{RFC: "RFC 8887", Method: "pushEnable", Account: alice.Name, Request: "dataTypes=Email"},
			emailStep,
			pushStateStep("RFC 8887"),
		}, "missing push")
		return
	}

	stateStep := pushStateStep(transportRFC)
	aliceID, ok := createPushDraft(t, sc, emailStep, pushSubject)
	if !ok {
		if sc.failed {
			notRunPush(sc, stateStep)
			t.Fail()
		} else {
			skipRest(t, sc, []step{stateStep}, "missing "+string(mail.URI))
		}
		return
	}
	waitPushEmail(t, sc, stateStep, aliceID, changes, listenErr)
	if sc.failed {
		t.Fail()
	}
}

func createPushDraft(t *testing.T, sc *scenario, st step, subject string) (jmap.ID, bool) {
	t.Helper()
	aliceID, err := alice.Client.PrimaryAccount(mail.URI)
	if err != nil {
		failPush(t, sc, st, err.Error())
		return "", false
	}
	inboxID, err := aliceInbox(alice.Client, aliceID)
	if err != nil {
		failPush(t, sc, st, err.Error())
		return "", false
	}
	_, ok := call[*email.SetResponse](t, sc, st, alice.Client, []jmap.URI{mail.URI}, false, &email.Set{
		Account: aliceID,
		Create: jmap.Some(map[jmap.ID]*email.Email{
			"m1": newE2EEmail(inboxID, subject, true),
		}),
	}, func(resp *email.SetResponse) (string, error) {
		created, err := createdEmail(resp, "m1")
		if err != nil {
			return "", err
		}
		return "id=" + string(created.ID), nil
	})
	return aliceID, ok
}

func aliceInbox(client *jmap.Client, account jmap.ID) (jmap.ID, error) {
	resp, err := jmap.Call[*mailbox.GetResponse](context.Background(), client, &mailbox.Get{Account: account})
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errString("empty mailbox list")
	}
	for _, mb := range resp.List {
		role, ok := mb.Role.Value()
		if ok && role == mailbox.RoleInbox {
			return mb.ID, nil
		}
	}
	return "", errString("inbox role missing")
}

func waitPushEmail(t *testing.T, sc *scenario, st step, account jmap.ID, changes <-chan *jmap.StateChange, listenErr <-chan error) {
	t.Helper()
	if sc.failed {
		notRunPush(sc, st)
		return
	}
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()
	seen := 0
	last := ""
	for {
		select {
		case change := <-changes:
			seen++
			last = pushChangeSummary(change)
			if token, ok := emailState(change, account); ok {
				recordPush(sc, st, "Email="+token, kindPass)
				return
			}
		default:
			select {
			case change := <-changes:
				seen++
				last = pushChangeSummary(change)
				if token, ok := emailState(change, account); ok {
					recordPush(sc, st, "Email="+token, kindPass)
					return
				}
			case err := <-listenErr:
				if err == nil || errors.Is(err, context.Canceled) {
					continue
				}
				failPush(t, sc, st, err.Error())
				return
			case <-timer.C:
				sc.failed = true
				recordPush(sc, st, "timeout", kindFail)
				if seen == 0 {
					t.Errorf("%s %s: timeout", sc.name, st.Method)
				} else {
					t.Errorf("%s %s: timeout (%d events, last %s)", sc.name, st.Method, seen, last)
				}
				return
			}
		}
	}
}

func emailState(change *jmap.StateChange, account jmap.ID) (string, bool) {
	if change == nil {
		return "", false
	}
	types, ok := change.Changed[account]
	if !ok {
		return "", false
	}
	token, ok := types["Email"]
	return token, ok
}

func pushChangeSummary(change *jmap.StateChange) string {
	if change == nil || len(change.Changed) == 0 {
		return "empty"
	}
	parts := make([]string, 0, len(change.Changed))
	for account, types := range change.Changed {
		names := make([]string, 0, len(types))
		for name := range types {
			names = append(names, name)
		}
		parts = append(parts, string(account)+"="+joinComma(names))
	}
	return joinComma(parts)
}

func pushStateStep(rfc string) step {
	return step{
		RFC: rfc, Method: "StateChange", Account: alice.Name, Request: "type=Email",
	}
}

func recordPush(sc *scenario, st step, response string, result kind) {
	rep.add(invocation{
		Scenario: sc.name,
		RFC:      st.RFC,
		Method:   st.Method,
		Account:  st.Account,
		Request:  st.Request,
		Response: response,
		Result:   result,
	})
}

func failPush(t *testing.T, sc *scenario, st step, response string) {
	t.Helper()
	sc.failed = true
	recordPush(sc, st, response, kindFail)
	t.Errorf("%s %s: %s", sc.name, st.Method, response)
}

func notRunPush(sc *scenario, steps ...step) {
	for _, st := range steps {
		recordPush(sc, st, "an earlier step failed", kindNotRun)
	}
}
