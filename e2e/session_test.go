//go:build e2e

package e2e

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core"
)

func TestSessionEcho(t *testing.T) {
	sc := &scenario{name: "Session"}
	for _, user := range []*account{alice, bob} {
		if user == nil || user.Client == nil || user.Client.Session == nil {
			t.Fatal("account session missing")
		}
		caps := make([]string, 0, len(user.Client.Session.RawCapabilities))
		for uri := range user.Client.Session.RawCapabilities {
			caps = append(caps, string(uri))
		}
		call[*core.Echo](t, sc, step{
			RFC: "RFC 8620", Method: "Core/echo", Account: user.Name,
			Request: "ping=e2e; capabilities=" + joinComma(caps),
		}, user.Client, []jmap.URI{jmap.CoreURI}, false, core.Echo{"ping": "e2e"}, func(echo *core.Echo) (string, error) {
			if echo == nil || (*echo)["ping"] != "e2e" {
				return "", errString("echo mismatch")
			}
			return "ping=e2e", nil
		})
	}
	if sc.failed {
		t.Fatal("session")
	}
}
