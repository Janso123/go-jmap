package websocket

import (
	"time"

	"github.com/Janso123/go-jmap/calendar"
	cws "github.com/coder/websocket"
)

// Compression modes for Options.CompressionMode (RFC 7692 permessage-deflate,
// advertised in RFC 8887 §4 examples). Aliased from coder/websocket.
const (
	CompressionDisabled          = cws.CompressionDisabled
	CompressionContextTakeover   = cws.CompressionContextTakeover
	CompressionNoContextTakeover = cws.CompressionNoContextTakeover
)

// Options configures Dial / DialURL behavior.
type Options struct {
	// CompressionMode selects permessage-deflate negotiation.
	// The zero value is CompressionDisabled (coder/websocket default).
	CompressionMode cws.CompressionMode

	// MaxConcurrentRequests limits in-flight Conn.Do calls on this connection
	// (RFC 8887 §4.3.2 / RFC 8620 Core maxConcurrentRequests).
	// If 0, the session Core capability value is used when > 0; otherwise unlimited.
	MaxConcurrentRequests uint64

	// Reconnect enables automatic reconnection after an unexpected disconnect.
	// Close() disables reconnect permanently for that Conn.
	Reconnect *ReconnectOptions

	// PingInterval, when > 0, starts a keepalive goroutine that calls Ping
	// on this interval. Zero disables automatic pings (default).
	PingInterval time.Duration

	// OnFrameError is called when a text frame cannot be decoded or has an
	// unknown @type. The connection stays open; see Conn readLoop for how
	// pending Do calls are unblocked when requestId is present.
	OnFrameError func(err error, raw []byte)

	// OnCalendarAlert is called when the server sends a CalendarAlert frame
	// (draft-ietf-jmap-calendars §6.4). Unknown other @types still go to
	// OnFrameError.
	OnCalendarAlert func(*calendar.CalendarAlert)
}

// ReconnectOptions controls auto-reconnect backoff and callbacks.
type ReconnectOptions struct {
	// MinBackoff is the initial wait after disconnect. Default: 1s.
	MinBackoff time.Duration
	// MaxBackoff caps exponential backoff. Default: 30s.
	MaxBackoff time.Duration
	// OnDisconnect is called when the underlying socket drops (before backoff).
	OnDisconnect func(error)
	// OnReconnect is called after a successful re-dial and optional re-EnablePush.
	OnReconnect func()
}

func mergeOptions(opts []Options) Options {
	var o Options
	if len(opts) > 0 {
		o = opts[0]
	}
	return o
}

func (r *ReconnectOptions) minBackoff() time.Duration {
	if r == nil || r.MinBackoff <= 0 {
		return time.Second
	}
	return r.MinBackoff
}

func (r *ReconnectOptions) maxBackoff() time.Duration {
	if r == nil || r.MaxBackoff <= 0 {
		return 30 * time.Second
	}
	return r.MaxBackoff
}
