package websocket

import "github.com/Janso123/go-jmap"

// URI is the JMAP capability for WebSocket transport (RFC 8887).
const URI jmap.URI = "urn:ietf:params:jmap:websocket"

// Subprotocol is the value negotiated in Sec-WebSocket-Protocol.
const Subprotocol = "jmap"

func init() {
	jmap.RegisterCapability(&WebSocket{})
}

// WebSocket is the session capability object for urn:ietf:params:jmap:websocket.
type WebSocket struct {
	// The wss-URI to use for initiating a JMAP-over-WebSocket handshake.
	URL string `json:"url"`

	// True if the server supports push notifications over the WebSocket.
	SupportsPush bool `json:"supportsPush"`
}

func (w *WebSocket) URI() jmap.URI { return URI }

func (w *WebSocket) New() jmap.Capability { return &WebSocket{} }
