# go-jmap

A JMAP client library for Go: typed methods, HTTPS, EventSource, and RFC 8887 WebSocket.

Module: [`github.com/Janso123/go-jmap`](https://github.com/Janso123/go-jmap). Go 1.27.1.

## Status: semi-stable release candidate

`v1.0.0-rc.1` is semi-stable. It is not `v1.0.0`.

The client API for published JMAP RFCs is frozen for this candidate. A break before `v1.0.0` will be called out in the release notes, not slipped in.

Still open before a stable release:

- Calendars and JSCalendar follow pinned drafts, not RFCs.
- HTTP/2 Extended CONNECT is unavailable. WebSocket uses the HTTP/1.1 upgrade.
- There is no high-level mail sync layer. You send typed structs through `Do` or `Call`.
- The Email type is large. Tests cover the wire shapes we rely on, not every field.

This module path is `github.com/Janso123/go-jmap`. It is not a drop-in for `rockorager/go-jmap` or earlier import paths.

```bash
go get github.com/Janso123/go-jmap@v1.0.0-rc.1
```

`github.com/coder/websocket` v1.8.x.

## Quick start

### HTTPS

```go
package main

import (
	"context"
	"fmt"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/Janso123/go-jmap/mail/mailbox"
)

func main() {
	client := jmap.NewClient(
		"https://api.fastmail.com/jmap/session",
		jmap.WithBearer("my-access-token"),
	)

	ctx := context.Background()
	if err := client.Authenticate(ctx); err != nil {
		panic(err)
	}

	id, err := client.PrimaryAccount(mail.URI)
	if err != nil {
		panic(err)
	}

	mboxResp, err := jmap.Call[*mailbox.GetResponse](ctx, client, &mailbox.Get{Account: id})
	if err != nil {
		panic(err)
	}
	for _, mbox := range mboxResp.List {
		fmt.Println("Mailbox:", mbox.Name)
	}

	req := &jmap.Request{}
	callID := req.Invoke(&email.Changes{
		Account:    id,
		SinceState: "some-known-state",
	})
	req.Invoke(&email.Get{
		Account: id,
		ReferenceIDs: &jmap.ResultReference{
			ResultOf: callID,
			Name:     "Email/changes",
			Path:     jmap.PathCreated, // "/created"
		},
	})

	resp, err := client.Do(ctx, req)
	if err != nil {
		panic(err)
	}
	for _, inv := range resp.Responses {
		if r, ok := inv.Args.(*email.GetResponse); ok {
			for _, eml := range r.List {
				if subject, ok := eml.Subject.Value(); ok {
					fmt.Println("Subject:", subject)
				}
			}
		}
	}
}
```

Call `RefreshSession` yourself when `SessionStale` is set. `Do` only marks the session stale.

### WebSocket

The session must advertise `urn:ietf:params:jmap:websocket`.

```go
import (
	"context"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/push/websocket"
	_ "github.com/Janso123/go-jmap/core"
)

conn, err := websocket.Dial(ctx, client, websocket.Options{
	CompressionMode: websocket.CompressionNoContextTakeover,
	Reconnect: &websocket.ReconnectOptions{
		MinBackoff: time.Second,
		MaxBackoff: 30 * time.Second,
	},
})
if err != nil {
	panic(err)
}
defer conn.Close()

conn.SetHandler(func(sc *jmap.StateChange) {
	// coalesce → Email/changes, etc.
})
_ = conn.EnablePush(nil, conn.PushState()) // nil dataTypes = all types

resp, err := conn.Do(ctx, req) // same *jmap.Request as Client.Do
```

Blobs stay on HTTPS (RFC 8887 §4). `maxConcurrentRequests` applies per WebSocket `Conn`. `DialH2Connect` returns `ErrH2ConnectUnsupported`.

## Imports

Blank-import a capability package before `Client.Do` or `Conn.Do`. Mail packages (`mail/email`, `mail/mailbox`, and the rest of core Mail) register through a normal import.

```go
import (
	_ "github.com/Janso123/go-jmap/core"
	_ "github.com/Janso123/go-jmap/core/blob"
	_ "github.com/Janso123/go-jmap/core/push/vapid"
	_ "github.com/Janso123/go-jmap/quota"
	_ "github.com/Janso123/go-jmap/sharing"
	_ "github.com/Janso123/go-jmap/sharing/principal"
	_ "github.com/Janso123/go-jmap/sharing/sharenotification"
	_ "github.com/Janso123/go-jmap/mail/sieve"
	_ "github.com/Janso123/go-jmap/contacts"
	_ "github.com/Janso123/go-jmap/contacts/addressbook"
	_ "github.com/Janso123/go-jmap/contacts/contactcard"
	_ "github.com/Janso123/go-jmap/calendar"
	_ "github.com/Janso123/go-jmap/calendar/calendars"
	_ "github.com/Janso123/go-jmap/calendar/calendarevent"
	_ "github.com/Janso123/go-jmap/calendar/participantidentity"
	_ "github.com/Janso123/go-jmap/calendar/calendareventnotification"
)
```

`contacts/jscontact` and `calendar/jscalendar` are typed Card and Event models. They do not register JMAP methods.

## API you can rely on

- `jmap.Call[T](ctx, client, method)` returns a typed response or `*MethodError`.
- `Get`, `Changes`, `Query`, `QueryChanges`, `Set`, and `Copy` are generic. Packages add the RFC fields.
- `Response.ByCallID` and `jmap.As[T]` read results. An unknown method becomes `*UnknownResponse` and does not fail the whole response.
- Filters compose with `jmap.And(...)`. Sort is `[]*jmap.Comparator`.
- `jmap.Optional[T]` keeps null, zero, and false distinct. An unset value omits the key. `Null()` sends JSON null. `Some(0)`, `Some(false)`, and `Some("")` stay on the wire.
- A nil `*bool` or `*UnsignedInt` omits the key. `new(false)` and `new(jmap.UnsignedInt(0))` send zero.
- `*jmap.UTCDate` is always `Z`. `jmap.Date` keeps the RFC 3339 offset.
- Unknown JSContact and JSCalendar properties land in `Extra`.
- Authorization is sent only to the session origin and to `apiUrl`, `uploadUrl`, `downloadUrl`, `eventSourceUrl`, and the WebSocket URL after `AllowAuthOrigin`.
- `WithBearer` and `WithBasic` wrap the current transport. `WithHTTPClient` after those options replaces the whole client, including auth.

## Specifications

Landscape: [jmap.io](https://jmap.io/spec.html).

**Done** means a usable client surface. **Partial** means a pinned draft. **Blocked** means a known transport gap.

| Spec | Status | Note |
|------|--------|------|
| [RFC 8620](https://www.rfc-editor.org/rfc/rfc8620) Core | Done | Session, requests, errors, blobs, push, Discover |
| [RFC 8887](https://www.rfc-editor.org/rfc/rfc8887) WebSocket | Done on HTTP/1.1 | H2 Extended CONNECT is blocked; see below |
| [RFC 9749](https://www.rfc-editor.org/rfc/rfc9749) VAPID | Done | Capability and `applicationServerKey` |
| [RFC 9670](https://www.rfc-editor.org/rfc/rfc9670) Sharing | Done | Principal and ShareNotification. Set is destroy-only |
| [RFC 9425](https://www.rfc-editor.org/rfc/rfc9425) Quotas | Done | get, changes, query, queryChanges |
| [RFC 9404](https://www.rfc-editor.org/rfc/rfc9404) Blob Management | Done | upload, get, lookup, plus Core `Blob/copy` |
| [RFC 8621](https://www.rfc-editor.org/rfc/rfc8621) Mail | Done | Types and methods. The Email surface is larger than the tests |
| [RFC 9007](https://www.rfc-editor.org/rfc/rfc9007) MDN | Done | send, parse, `$mdnsent` |
| [RFC 9219](https://www.rfc-editor.org/rfc/rfc9219) S/MIME verify | Done | Capability, fields, and filters |
| [RFC 9661](https://www.rfc-editor.org/rfc/rfc9661) Sieve | Done | get, set, query, validate |
| [RFC 9610](https://www.rfc-editor.org/rfc/rfc9610) Contacts | Done | AddressBook and ContactCard |
| [RFC 9553](https://www.rfc-editor.org/rfc/rfc9553) JSContact | Done | Typed Card |
| [Calendars draft 29](https://datatracker.ietf.org/doc/draft-ietf-jmap-calendars/29/) | Partial | Methods are implemented. Stays partial until the RFC |
| [JSCalendar bis 20](https://datatracker.ietf.org/doc/html/draft-ietf-calext-jscalendarbis-20) | Partial | Typed Event core. Task and Group are not typed |

HTTP/2 Extended CONNECT (RFC 8887 §4.2) is blocked. `DialH2Connect` returns `ErrH2ConnectUnsupported` until [coder/websocket#4](https://github.com/coder/websocket/issues/4) lands.

## Tests

Unit tests need no server. CI runs them with `-race` on pushes to `main` and on pull requests.

```bash
go test -race -count=1 ./...
```

The Stalwart suite is behind `-tags e2e`. It starts `stalwartlabs/stalwart:v0.16` on `127.0.0.1:18080`, runs session, mail, contacts, calendar, blob, quota, and push, writes `e2e/report.md`, and removes the container. `go test ./...` does not start Docker. The report is gitignored.

```bash
go test -tags e2e -count=1 -timeout 10m -v ./e2e/
```

The same suite runs in CI. A version tag also runs it and puts the report in the GitHub Release, collapsed under "Stalwart e2e (pass N, fail N)". The release is not published when that suite fails.

## Releasing

Push a `v*` tag. Tags containing `alpha`, `beta`, or `rc` are prereleases. See [`.github/workflows/release.yml`](.github/workflows/release.yml).

## License

[MIT](LICENSE).

Copyright © 2019 Max Mazurov; Copyright © 2022 Tim Culverhouse; and subsequent contributors.
