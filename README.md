# go-jmap (fork)

A JMAP client library for Go, brought up to current
[JMAP RFC / jmap.io](https://jmap.io/spec.html) coverage for mail-client use
(WebSocket, published extensions, contacts, calendars drafts, and related fixes).

## Lineage

This is a **fork of a fork**:

1. Originally [github.com/foxcpp/go-jmap](https://github.com/foxcpp/go-jmap)
2. Continued as [`git.sr.ht/~rockorager/go-jmap`](https://git.sr.ht/~rockorager/go-jmap)
   (upstream here) — Core (RFC 8620), Mail (RFC 8621), MDN, S/MIME verify
3. **This repository** extends that tree to fuller current RFC / jmap.io coverage

We keep the upstream module path (`git.sr.ht/~rockorager/go-jmap`) so consumers
can switch with a `replace` or by changing the remote once upstream catches up.

## Maintenance policy

We **actively maintain this fork** (bugs, features, RFC catch-up) **until**
upstream merges equivalent work or otherwise provides the same coverage.

**When upstream catches up / absorbs this work, maintenance of this fork stops.**
Consumers should then switch back to
[`git.sr.ht/~rockorager/go-jmap`](https://git.sr.ht/~rockorager/go-jmap).
Until that happens, treat this fork as the source of truth for the additions
below.

## What this fork adds

| Change | Detail |
|--------|--------|
| **RFC 8887 WebSocket** | `core/push/websocket` — Dial/Do, push frames, `pushState`, permessage-deflate, per-conn concurrency gate, reconnect (`github.com/coder/websocket`) |
| **jmap.io extensions** | VAPID, Sharing, Quotas, Blob Management, Sieve, Contacts/JSContact; Calendars/JSCalendar on pinned IETF drafts |
| **Blob upload 2xx** | Treat any HTTP 2xx as a successful upload (RFC 8620 §6.1); e.g. Cyrus `201 Created` |
| **Stable `using` URIs** | `mergeURIs` preserves capability order (target then extras) instead of map iteration |

**Not in this fork yet:** HTTP/2 Extended CONNECT (RFC 8887 §4.2) — API stub returns
`ErrH2ConnectUnsupported` until [`coder/websocket#4`](https://github.com/coder/websocket/issues/4).

## Module path / import

Keep the upstream module path (do not rewrite):

```text
git.sr.ht/~rockorager/go-jmap
```

Local `replace` (example):

```go
replace git.sr.ht/~rockorager/go-jmap => ../../forks/go-jmap
```

Go **1.23+** (required by `github.com/coder/websocket` v1.8.x).

Blank-import packages whose `init()` registers capabilities/methods before
`Client.Do` / `websocket.Conn.Do`:

```go
import (
	_ "git.sr.ht/~rockorager/go-jmap/core"
	_ "git.sr.ht/~rockorager/go-jmap/core/blob"
	_ "git.sr.ht/~rockorager/go-jmap/core/push/vapid"
	_ "git.sr.ht/~rockorager/go-jmap/quota"
	_ "git.sr.ht/~rockorager/go-jmap/sharing"
	_ "git.sr.ht/~rockorager/go-jmap/sharing/principal"
	_ "git.sr.ht/~rockorager/go-jmap/sharing/sharenotification"
	_ "git.sr.ht/~rockorager/go-jmap/mail/sieve"
	_ "git.sr.ht/~rockorager/go-jmap/contacts"
	_ "git.sr.ht/~rockorager/go-jmap/contacts/addressbook"
	_ "git.sr.ht/~rockorager/go-jmap/contacts/contactcard"
	_ "git.sr.ht/~rockorager/go-jmap/calendar"
	_ "git.sr.ht/~rockorager/go-jmap/calendar/calendars"
	_ "git.sr.ht/~rockorager/go-jmap/calendar/calendarevent"
	_ "git.sr.ht/~rockorager/go-jmap/calendar/participantidentity"
	_ "git.sr.ht/~rockorager/go-jmap/calendar/calendareventnotification"
)
```

Import `contacts/jscontact` and `calendar/jscalendar` when constructing typed
Card/Event values; they do not register JMAP methods on their own.

## Quick start

### HTTPS (session + method calls)

```go
package main

import (
	"fmt"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
	"git.sr.ht/~rockorager/go-jmap/mail/mailbox"
)

func main() {
	client := &jmap.Client{
		SessionEndpoint: "https://api.fastmail.com/jmap/session",
	}
	client.WithAccessToken("my-access-token")

	if err := client.Authenticate(); err != nil {
		// handle
	}

	id := client.Session.PrimaryAccounts[mail.URI]
	req := &jmap.Request{}
	req.Invoke(&mailbox.Get{Account: id})
	callID := req.Invoke(&email.Changes{
		Account:    id,
		SinceState: "some-known-state",
	})
	req.Invoke(&email.Get{
		Account: id,
		ReferenceIDs: &jmap.ResultReference{
			ResultOf: callID,
			Name:     "Email/changes",
			Path:     "/created",
		},
	})

	resp, err := client.Do(req)
	if err != nil {
		// handle
	}
	for _, inv := range resp.Responses {
		switch r := inv.Args.(type) {
		case *mailbox.GetResponse:
			for _, mbox := range r.List {
				fmt.Printf("Mailbox: %s\n", mbox.Name)
			}
		case *email.GetResponse:
			for _, eml := range r.List {
				fmt.Printf("Subject: %s\n", eml.Subject)
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

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/core/push/websocket"
	_ "git.sr.ht/~rockorager/go-jmap/core"
)

conn, err := websocket.Dial(ctx, client, websocket.Options{
	CompressionMode: websocket.CompressionNoContextTakeover,
	Reconnect: &websocket.ReconnectOptions{
		MinBackoff: time.Second,
		MaxBackoff: 30 * time.Second,
	},
})
if err != nil {
	// handle
}
defer conn.Close()

conn.SetHandler(func(sc *jmap.StateChange) {
	// coalesce → Email/changes, etc.
})
_ = conn.EnablePush(nil, conn.PushState()) // nil dataTypes = all types

resp, err := conn.Do(ctx, req) // same *jmap.Request as Client.Do
```

Blobs stay on HTTPS (RFC 8887 §4). `maxConcurrentRequests` is enforced per
WebSocket `Conn`; budget HTTPS + WS together if you multiplex.

## Spec coverage

Landscape: [jmap.io Specifications](https://jmap.io/spec.html).

Status: **Done** = usable client types/methods · **Partial** = draft-pinned or
intentionally incomplete · **Blocked** = known transport gap.

### Core

| Spec | Status | Notes |
|------|--------|-------|
| [RFC 8620](https://www.rfc-editor.org/rfc/rfc8620) Core | Done | Session, Request/Response, errors, HTTPS blobs, `Blob/copy`, PushSubscription, EventSource, `core.Discover` |
| [RFC 8887](https://www.rfc-editor.org/rfc/rfc8887) WebSocket | Done | H1 Upgrade via `coder/websocket`; H2 Extended CONNECT **Blocked** (stub) |
| [RFC 9749](https://www.rfc-editor.org/rfc/rfc9749) VAPID | Done | `core/push/vapid` |
| [RFC 9670](https://www.rfc-editor.org/rfc/rfc9670) Sharing | Done | Principal + ShareNotification |
| [RFC 9425](https://www.rfc-editor.org/rfc/rfc9425) Quotas | Done | `quota` |
| [RFC 9404](https://www.rfc-editor.org/rfc/rfc9404) Blob Management | Done | `core/blob` — upload/get/lookup (+ Core `Blob/copy`) |

### Mail

| Spec | Status | Notes |
|------|--------|-------|
| [RFC 8621](https://www.rfc-editor.org/rfc/rfc8621) Mail | Done | Mailbox/Thread/Email/Identity/Submission/Vacation/SearchSnippet |
| [RFC 9007](https://www.rfc-editor.org/rfc/rfc9007) MDN | Done | |
| [RFC 9219](https://www.rfc-editor.org/rfc/rfc9219) S/MIME verify | Done | Capability + Email SMIME fields |
| [RFC 9661](https://www.rfc-editor.org/rfc/rfc9661) Sieve | Done | `mail/sieve` |

### Contacts and calendars

| Spec | Status | Notes |
|------|--------|-------|
| [RFC 9610](https://www.rfc-editor.org/rfc/rfc9610) Contacts | Done | AddressBook + ContactCard |
| [RFC 9553](https://www.rfc-editor.org/rfc/rfc9553) JSContact | Done | `contacts/jscontact` Card model |
| JMAP Calendars (draft) | Partial | Pinned [draft-ietf-jmap-calendars-29](https://datatracker.ietf.org/doc/draft-ietf-jmap-calendars/29/); methods complete, stays Partial until RFC |
| JSCalendar 2.0 (draft) | Partial | Pinned [draft-ietf-calext-jscalendarbis-20](https://datatracker.ietf.org/doc/html/draft-ietf-calext-jscalendarbis-20); typed Event core; Task/Group and some Event fields still via extensions |

### Blocked

| Item | Status | Notes |
|------|--------|-------|
| HTTP/2 Extended CONNECT | Blocked | `DialH2Connect` → `ErrH2ConnectUnsupported`; needs `coder/websocket` client support |

## Remotes

| Remote | URL |
|--------|-----|
| `upstream` | `https://git.sr.ht/~rockorager/go-jmap` |
| `origin` (this fork) | `https://github.com/Janso123/go-jmap.git` |

Do not rewrite the module path. See **Maintenance policy** above.

### Contributing upstream (sourcehut)

Upstream is on [sourcehut](https://sr.ht/~rockorager/go-jmap/): patches via email,
not GitHub PRs against the canonical repo.

1. Open or reference a [ticket](https://todo.sr.ht/~rockorager/go-jmap) if useful.
2. Send a patch series to **[~rockorager/go-jmap-devel](https://lists.sr.ht/~rockorager/go-jmap-devel)**
   (`~rockorager/go-jmap-devel@lists.sr.ht`) using
   [`git send-email`](https://git-send-email.io) or the
   [git.sr.ht prepare-patchset UI](https://man.sr.ht/git.sr.ht/#sending-patches-upstream).
3. Prefer small, reviewable series (see below). Announce list:
   [go-jmap-announce](https://lists.sr.ht/~rockorager/go-jmap-announce).

### Recommended upstream series (phased)

Do **not** dump the entire fork (~30+ commits spanning WS + every extension) in
one patchset — unlikely to review or land.

| Phase | Scope | Why first |
|-------|--------|-----------|
| 1 | Blob upload accept any 2xx | One-line RFC correctness fix; easy review |
| 2 | RFC 8887 WebSocket (`core/push/websocket`) | Highest unique value vs upstream; keep optional extras (reconnect/deflate) clearly separated if maintainer prefers a minimal Dial/Do first |
| 3 | Published RFCs one package at a time | VAPID → Quotas → Blob mgmt → Sharing → Sieve → Contacts/JSContact |
| 4 | Calendars / JSCalendar | After drafts stabilize or with explicit draft pins in docs |

License allows redistribution and contribution under MIT (keep copyright notices).

## Tests

```bash
go test ./...
```

## License

[MIT](LICENSE) — SPDX: `MIT`.

Copyright © 2019 Max Mazurov (fox.cpp); Copyright © 2022 Tim Culverhouse; and
subsequent contributors.

Obligations: include the copyright and permission notice in all copies or
substantial portions. No copyleft / same-license requirement beyond MIT terms.
See [LICENSE](LICENSE).
