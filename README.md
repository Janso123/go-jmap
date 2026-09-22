# go-jmap

A JMAP **client** library for Go: typed methods, session/HTTPS transport, EventSource, and RFC 8887 WebSocket.

**Module:** [`github.com/Janso123/go-jmap`](https://github.com/Janso123/go-jmap) · **Tag:** `v1.0.0-rc.1` · **Go:** 1.27.1

Release candidate: APIs may still change before a stable `v1.0.0`. Not a drop-in for older rockorager import paths.

## Install

```bash
go get github.com/Janso123/go-jmap@v1.0.0-rc.1
```

Requires Go **1.27.1** (`go 1.27.1` in `go.mod`; `github.com/coder/websocket` v1.8.x).

## Quick start

### HTTPS (session + method calls)

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

	// One-shot typed call
	mboxResp, err := jmap.Call[*mailbox.GetResponse](ctx, client, &mailbox.Get{Account: id})
	if err != nil {
		panic(err)
	}
	for _, mbox := range mboxResp.List {
		fmt.Println("Mailbox:", mbox.Name)
	}

	// Batch with result references
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

### WebSocket (RFC 8887)

Requires a session that advertises `urn:ietf:params:jmap:websocket`.

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

Blobs stay on HTTPS (RFC 8887 §4). `maxConcurrentRequests` is enforced per WebSocket `Conn`; budget HTTPS + WS together if you multiplex.

HTTP/2 Extended CONNECT is **not** available (`DialH2Connect` returns `ErrH2ConnectUnsupported`).

## Capability registration (blank imports)

Packages register capabilities and method response types in `init()`. Blank-import what you use **before** `Client.Do` / `websocket.Conn.Do`:

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

Import `contacts/jscontact` and `calendar/jscalendar` when building typed Card/Event values; they do not register JMAP methods.

Mail packages (`mail/email`, `mail/mailbox`, …) register via their normal imports — no extra blank import needed for core Mail.

## API highlights (v1 release candidate)

| Area | What you get |
|------|----------------|
| **Generic method kit** | `Get` / `Changes` / `Query` / `QueryChanges` / `Set` / `Copy` parameterized on object types; thin packages embed and add RFC fields |
| **One-shot** | `jmap.Call[T](ctx, client, method)` — typed response or `*MethodError` |
| **Responses** | `Response.ByCallID`, `jmap.As[T]`; unknown methods → `*UnknownResponse` (does not fail the whole response) |
| **Filters / sort** | `jmap.And(email.InMailbox(...))`; domain `FilterCondition` implements `jmap.Filter`; sort is `[]*jmap.Comparator` |
| **JSContact / JSCalendar** | Unknown properties in public `Extra map[string]jsontext.Value` (`json:",embed"`); json/v2 |
| **UTCDate / Date** | RFC 8620 §1.4: `*jmap.UTCDate` always `Z`; `jmap.Date` is RFC 3339 date-time with offset preserved |
| **Transport** | Context on all I/O; origin-scoped `Authorization`; RFC 6570 L1 URI templates; `problem+json` → `RequestError`; non-JSON HTTP → `HTTPError`; blob upload/download; Discover (SRV + `.well-known`); `WithTrustedHosts`, `WithTimeout` |
| **Session** | `Do` marks stale on `sessionState` mismatch — call `SessionStale()` / `RefreshSession(ctx)` yourself |
| **Push** | EventSource (`core/push`) and WebSocket (`core/push/websocket`) with push enable/disable and `pushState` |

### Release candidate caveats

| Caveat | Detail |
|--------|--------|
| Auth option order | `WithBearer` / `WithBasic` wrap the current `Transport` and keep Timeout, CheckRedirect, and Jar. Credentials are sent only to the session origin plus `apiUrl`/`uploadUrl`/`downloadUrl`/`eventSourceUrl` (and WebSocket URL after `AllowAuthOrigin`). `WithHTTPClient` after auth replaces the whole client (including auth). Timeout / TrustedHosts compose with Bearer/Basic. |
| Null, zero, and false | `T\|null` is `jmap.Optional[T]` with `omitzero`: an unset value omits the key, `Null()` is JSON null, and `Some(0)` / `Some(false)` / `Some("")` stay on the wire. A nil `*bool` or `*UnsignedInt` omits the key; `new(false)` and `new(jmap.UnsignedInt(0))` send zero. `Bool`, `UintPtr`, and `IDPtr` are `//go:fix inline` wrappers around `new`. |
| No H2 CONNECT | Use HTTP/1.1 WebSocket upgrade |
| Calendars | Draft-pinned; stays **Partial** until RFCs ship (see below) |
| No fluent one-shots | Typed structs + `Do` / `Call` / `Conn.Do` — not a high-level mail sync layer |

## Spec coverage

Landscape: [jmap.io Specifications](https://jmap.io/spec.html).

Two different measures:

- **Spec status** — whether client types/methods for that RFC (or draft) are implemented and usable.
- **Test coverage (approx)** — Go statement coverage from `go test ./… -cover` for the main package(s). High package % does not mean every RFC semantic is asserted; larger packages (especially `mail/email`) still leave implementation ahead of tests.

**Spec status legend:** **Done** = usable client surface · **Partial** = draft-pinned or intentionally incomplete · **Blocked** = known gap (transport).

Published RFCs on jmap.io are **Done** for a client library except calendar drafts (**Partial**) and H2 CONNECT (**Blocked**). Roughly the full published matrix is implemented; remaining gaps are drafts + one transport stub.

### Core

| Spec | Spec status | Test coverage (approx) | Notes |
|------|-------------|------------------------|-------|
| [RFC 8620](https://www.rfc-editor.org/rfc/rfc8620) Core | Done | root ~69%; `core` ~83%; EventSource ~80%; subscription ~67% | Session, Request/Response, errors, HTTPS blobs, `Blob/copy`, PushSubscription, EventSource, Discover |
| [RFC 8887](https://www.rfc-editor.org/rfc/rfc8887) WebSocket | Done (H1) | `websocket` ~81% | H1 Upgrade via `coder/websocket`; H2 Extended CONNECT **Blocked** (stub) |
| [RFC 9749](https://www.rfc-editor.org/rfc/rfc9749) VAPID | Done | `vapid` ~100% | Capability + `applicationServerKey` |
| [RFC 9670](https://www.rfc-editor.org/rfc/rfc9670) Sharing | Done | sharing 62–100% | Principal + ShareNotification (set is destroy-only) |
| [RFC 9425](https://www.rfc-editor.org/rfc/rfc9425) Quotas | Done | `quota` ~64% | get/changes/query/queryChanges |
| [RFC 9404](https://www.rfc-editor.org/rfc/rfc9404) Blob Management | Done | `blob` ~70% | upload/get/lookup (+ Core `Blob/copy`) |

### Mail

| Spec | Spec status | Test coverage (approx) | Notes |
|------|-------------|------------------------|-------|
| [RFC 8621](https://www.rfc-editor.org/rfc/rfc8621) Mail | Done | `email` ~57%; mailbox ~58%; others ~67–100% | Types/methods present; filters/keywords/`header:*`/S/MIME/Import/Parse wire fixtures; large Email surface still incomplete |
| [RFC 9007](https://www.rfc-editor.org/rfc/rfc9007) MDN | Done | `mdn` ~82% | send/parse; `$mdnsent` keyword |
| [RFC 9219](https://www.rfc-editor.org/rfc/rfc9219) S/MIME verify | Done | (in `email`, ~57%) | Capability + Email SMIME fields/filters; `using` includes `smimeverify` when properties or filters need it |
| [RFC 9661](https://www.rfc-editor.org/rfc/rfc9661) Sieve | Done | `sieve` ~88% | get/set/query/validate |

### Contacts and calendars

| Spec | Spec status | Test coverage (approx) | Notes |
|------|-------------|------------------------|-------|
| [RFC 9610](https://www.rfc-editor.org/rfc/rfc9610) Contacts | Done | addressbook ~86%; contactcard ~72% | AddressBook + ContactCard |
| [RFC 9553](https://www.rfc-editor.org/rfc/rfc9553) JSContact | Done | `jscontact` ~55% | Typed Card model |
| JMAP Calendars (draft) | Partial | calendar pkgs ~70–100% | Pinned [draft-ietf-jmap-calendars-29](https://datatracker.ietf.org/doc/draft-ietf-jmap-calendars/29/); methods complete; stays Partial until RFC |
| JSCalendar 2.0 (draft) | Partial | `jscalendar` ~47% | Pinned [draft-ietf-calext-jscalendarbis-20](https://datatracker.ietf.org/doc/html/draft-ietf-calext-jscalendarbis-20); typed Event core including `priority`/`privacy`/`freeBusyStatus`; unknown keys in `Extra`; Task/Group not typed |

### Blocked

| Item | Spec status | Notes |
|------|-------------|-------|
| HTTP/2 Extended CONNECT (RFC 8887 §4.2) | Blocked | `DialH2Connect` → `ErrH2ConnectUnsupported`; needs [coder/websocket#4](https://github.com/coder/websocket/issues/4) |

## Tests

Unit tests do not need a server:

```bash
go test -race -count=1 ./...
go test ./... -cover   # package statement coverage
```

CI runs those race tests on pushes and pull requests to `main`. It does not start Docker.

### Stalwart end-to-end

`e2e/` is built only with `-tags e2e`. The test process starts `stalwartlabs/stalwart:v0.16` through Docker Compose, creates `alice@example.org` and `bob@example.org`, and exercises session, mail, contacts, calendar, blob, quota, and push. It then writes `e2e/report.md` and removes the container. `go test ./...` without the tag does not start Docker.

Docker is required. The suite binds `127.0.0.1:18080`. The `e2e` job in [`.github/workflows/go.yml`](.github/workflows/go.yml) runs it on pushes to `main` and on pull requests. The unit-test job does not pass `-tags e2e` and does not start Docker.

```bash
go test -tags e2e -count=1 -timeout 10m -v ./e2e/
```

`-v` prints `e2e report: e2e/report.md` after the journal is written. The report is gitignored.

## Releasing

Push a version tag (`vX.Y.Z` or prerelease like `v1.0.0-rc.1`). The [release workflow](.github/workflows/release.yml) runs tests, then creates a GitHub Release with auto-generated notes (tags containing `alpha` / `beta` / `rc` are marked prerelease).

```bash
git tag v1.0.0-rc.1
git push origin v1.0.0-rc.1
```

## License

[MIT](LICENSE) — SPDX: `MIT`.

Copyright © 2019 Max Mazurov; Copyright © 2022 Tim Culverhouse; and subsequent contributors.

See [LICENSE](LICENSE) for terms.
