package websocket

// H1 spike (2026-09-19): RFC 8887 §4.2 / RFC 8441 HTTP/2 Extended CONNECT with
// github.com/coder/websocket v1.8.15 — outcome **Blocked**.
//
// Goal: dial JMAP WebSocket over H2 CONNECT (:method CONNECT, :protocol websocket)
// using only coder/websocket, or pass an existing net.Conn from a manual CONNECT
// into the library.
//
// Dial API (module cache: dial.go):
//   - websocket.Dial(ctx, url, *DialOptions) is the only client entrypoint.
//   - DialOptions: HTTPClient, HTTPHeader, Host, Subprotocols, CompressionMode,
//     CompressionThreshold, OnPingReceived, OnPongReceived — no net.Conn, no
//     alternate handshake mode, no HTTP/2 CONNECT hook.
//
// Handshake path is fixed to HTTP/1.1 Upgrade:
//   - handshakeRequest sends GET with Connection: Upgrade and Upgrade: websocket.
//   - verifyServerResponse requires http.StatusSwitchingProtocols (101) and
//     Connection/Upgrade/Sec-WebSocket-Accept headers (RFC 6455), not an H2
//     CONNECT 200 with :protocol websocket.
//
// HTTPClient customization does not help:
//   - A custom Transport could perform H2 CONNECT elsewhere, but Dial still
//     issues GET Upgrade via HTTPClient.Do and validates a 101 response; it
//     never accepts a pre-established tunnel Conn from outside the response body.
//   - After handshake, Dial type-asserts resp.Body to io.ReadWriteCloser and
//     passes it to unexported newConn — there is no exported NewClientConn(rwc).
//
// net.Conn direction is server-oriented only:
//   - websocket.NetConn(ctx, *Conn, msgType) wraps an existing *Conn as net.Conn
//     (typically after Accept). There is no inverse for client dial.
//
// Upstream tracking: https://github.com/coder/websocket/issues/4 (HTTP/2 CONNECT).
//
// Next (Task H2): export DialH2Connect returning ErrH2ConnectUnsupported until
// coder/websocket adds client-side H2 CONNECT or a documented way to adopt a
// post-CONNECT io.ReadWriteCloser as *Conn.

import (
	"context"
	"errors"

	"git.sr.ht/~rockorager/go-jmap"
)

// ErrH2ConnectUnsupported is returned by DialH2Connect while github.com/coder/websocket
// lacks RFC 8441 client Extended CONNECT (see h2connect.go spike comments).
var ErrH2ConnectUnsupported = errors.New("jmap websocket: HTTP/2 Extended CONNECT unsupported by github.com/coder/websocket")

// DialH2Connect attempts RFC 8887 §4.2 / RFC 8441. Until coder/websocket
// supports it, this returns ErrH2ConnectUnsupported.
func DialH2Connect(ctx context.Context, client *jmap.Client, opts ...Options) (*Conn, error) {
	_ = ctx
	_ = client
	_ = mergeOptions(opts)
	return nil, ErrH2ConnectUnsupported
}
