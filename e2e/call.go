//go:build e2e

package e2e

import (
	"context"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core"
)

type account struct {
	Name     string
	Email    string
	Password string
	Client   *jmap.Client
}

var (
	alice *account
	bob   *account
)

type scenario struct {
	name   string
	failed bool
}

type step struct {
	RFC     string
	Method  string
	Account string
	Request string
}

func hasCap(c *jmap.Client, uri jmap.URI) bool {
	_, ok := c.Session.Capabilities[uri]
	return ok
}

func maxInSet(c *jmap.Client) int {
	coreCap, ok := c.Session.Capabilities[jmap.CoreURI].(*core.Core)
	if !ok || coreCap.MaxObjectsInSet == 0 {
		return 30
	}
	return int(coreCap.MaxObjectsInSet)
}

func call[T jmap.MethodResponse](t *testing.T, sc *scenario, st step, client *jmap.Client, caps []jmap.URI, allowUnknown bool, method jmap.Method, check func(T) (string, error)) (T, bool) {
	t.Helper()
	var zero T
	if sc.failed {
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: "an earlier step failed", Result: kindNotRun})
		return zero, false
	}
	for _, uri := range caps {
		if !hasCap(client, uri) {
			rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: "missing " + string(uri), Result: kindSkip})
			return zero, false
		}
	}
	got, err := jmap.Call[T](context.Background(), client, method)
	if err != nil {
		var me *jmap.MethodError
		if allowUnknown && errors.As(err, &me) && me.Type == string(jmap.MethodErrUnknownMethod) {
			rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: "unknownMethod", Result: kindSkip})
			return zero, false
		}
		sc.failed = true
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: err.Error(), Result: kindFail})
		t.Errorf("%s %s: %v", sc.name, st.Method, err)
		return zero, false
	}
	summary, err := check(got)
	if err != nil {
		sc.failed = true
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: err.Error(), Result: kindFail})
		t.Errorf("%s %s: %v", sc.name, st.Method, err)
		return got, false
	}
	rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: summary, Result: kindPass})
	return got, true
}

func uploadBlob(t *testing.T, sc *scenario, st step, client *jmap.Client, account jmap.ID, body, contentType string) (jmap.ID, bool) {
	t.Helper()
	if sc.failed {
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: "an earlier step failed", Result: kindNotRun})
		return "", false
	}
	resp, err := client.Upload(context.Background(), account, strings.NewReader(body), contentType)
	if err != nil {
		sc.failed = true
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: err.Error(), Result: kindFail})
		t.Errorf("%s %s: %v", sc.name, st.Method, err)
		return "", false
	}
	if resp == nil || resp.ID == "" {
		sc.failed = true
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: "blob id missing", Result: kindFail})
		t.Errorf("%s %s: blob id missing", sc.name, st.Method)
		return "", false
	}
	summary := "blobId=" + string(resp.ID)
	rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: summary, Result: kindPass})
	return resp.ID, true
}

func downloadBlob(t *testing.T, sc *scenario, st step, client *jmap.Client, account, blobID jmap.ID, want string) ([]byte, bool) {
	t.Helper()
	if sc.failed {
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: "an earlier step failed", Result: kindNotRun})
		return nil, false
	}
	rc, err := client.Download(context.Background(), account, blobID, jmap.DownloadOptions{Type: "text/plain"})
	if err != nil {
		sc.failed = true
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: err.Error(), Result: kindFail})
		t.Errorf("%s %s: %v", sc.name, st.Method, err)
		return nil, false
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		sc.failed = true
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: err.Error(), Result: kindFail})
		t.Errorf("%s %s: %v", sc.name, st.Method, err)
		return nil, false
	}
	if want != "" && string(data) != want {
		sc.failed = true
		msg := "got " + strconv.Quote(string(data))
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: msg, Result: kindFail})
		t.Errorf("%s %s: %s", sc.name, st.Method, msg)
		return data, false
	}
	summary := "bytes=" + strconv.Itoa(len(data))
	rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: summary, Result: kindPass})
	return data, true
}

func skipRest(t *testing.T, sc *scenario, steps []step, reason string) {
	t.Helper()
	for _, st := range steps {
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: reason, Result: kindSkip})
	}
}

func joinComma(parts []string) string { return strings.Join(parts, ", ") }

func errString(msg string) error { return errors.New(msg) }
