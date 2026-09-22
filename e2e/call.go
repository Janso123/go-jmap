//go:build e2e

package e2e

import (
	"context"
	"errors"
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

func skipRest(t *testing.T, sc *scenario, steps []step, reason string) {
	t.Helper()
	for _, st := range steps {
		rep.add(invocation{Scenario: sc.name, RFC: st.RFC, Method: st.Method, Account: st.Account, Request: st.Request, Response: reason, Result: kindSkip})
	}
}

func joinComma(parts []string) string { return strings.Join(parts, ", ") }

func errString(msg string) error { return errors.New(msg) }
