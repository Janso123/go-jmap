package jmap

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// originAuthTransport adds Authorization only for requests whose origin
// (scheme+host+port) is in the allow-set. Cross-origin redirects therefore do
// not receive credentials (RFC 9110 §15.4 / RFC 8620 §8.3).
type originAuthTransport struct {
	base   http.RoundTripper
	header string

	mu      sync.RWMutex
	origins map[string]struct{}
}

func (t *originAuthTransport) allow(origin string) {
	if origin == "" {
		return
	}
	t.mu.Lock()
	if t.origins == nil {
		t.origins = make(map[string]struct{})
	}
	t.origins[origin] = struct{}{}
	t.mu.Unlock()
}

func (t *originAuthTransport) allowed(origin string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.origins[origin]
	return ok
}

func (t *originAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	if t.allowed(originOfURL(req.URL)) {
		req.Header.Set("Authorization", t.header)
	} else {
		req.Header.Del("Authorization")
	}
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

func unwrapOriginAuth(rt http.RoundTripper) http.RoundTripper {
	if t, ok := rt.(*originAuthTransport); ok {
		return t.base
	}
	return rt
}

func originOf(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return originOfURL(u)
}

func originOfURL(u *url.URL) string {
	if u == nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	scheme := strings.ToLower(u.Scheme)
	host := strings.ToLower(u.Hostname())
	if host == "" {
		return ""
	}
	port := u.Port()
	switch {
	case (scheme == "https" || scheme == "wss") && (port == "" || port == "443"):
		return scheme + "://" + host
	case (scheme == "http" || scheme == "ws") && (port == "" || port == "80"):
		return scheme + "://" + host
	case port != "":
		return scheme + "://" + host + ":" + port
	default:
		return scheme + "://" + host
	}
}

// CheckWebSocketURL reports whether wsURL may be dialed for this client.
// The URL must be wss, except ws on a loopback host. Unless allowForeign is
// set, its HTTP-equivalent origin must match the session endpoint origin
// when that origin is known (RFC 8887 §4 / RFC 8620 §2).
func (c *Client) CheckWebSocketURL(wsURL string, allowForeign bool) error {
	u, err := url.Parse(wsURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("jmap: invalid websocket url")
	}
	switch strings.ToLower(u.Scheme) {
	case "wss":
	case "ws":
		if !isLoopbackHost(u.Hostname()) {
			return fmt.Errorf("jmap: websocket url must use wss")
		}
	default:
		return fmt.Errorf("jmap: websocket url must use wss")
	}
	if allowForeign || c == nil {
		return nil
	}
	c.Lock()
	endpoint := c.SessionEndpoint
	c.Unlock()
	want := originOf(endpoint)
	if want == "" {
		return nil
	}
	got := originOf(httpURLForWebSocket(wsURL))
	if got == "" || got != want {
		return fmt.Errorf("jmap: refusing websocket origin %s (session origin %s)", got, want)
	}
	return nil
}

// WebSocketResponseOriginOK reports whether the handshake response stayed on
// the HTTP-equivalent origin of wsURL. WebSocket handshakes are HTTP requests,
// so ws/wss is compared as http/https.
func WebSocketResponseOriginOK(wsURL string, resp *http.Response) error {
	return ResponseOriginOK(httpURLForWebSocket(wsURL), resp)
}

func httpURLForWebSocket(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u == nil {
		return raw
	}
	switch strings.ToLower(u.Scheme) {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	}
	return u.String()
}

// AllowAuthOrigin adds the origin (scheme+host+port) of rawURL to the
// credential allow-set used by WithBearer/WithBasic. Call this before
// connecting to a WebSocket URL whose origin is not the session origin.
func (c *Client) AllowAuthOrigin(rawURL string) {
	if c == nil || c.HttpClient == nil {
		return
	}
	t, ok := c.HttpClient.Transport.(*originAuthTransport)
	if !ok {
		return
	}
	o := originOf(rawURL)
	t.allow(o)
	// WebSocket handshakes are HTTP(S); ws/wss origins must also allow the
	// corresponding http/https origin so Authorization is sent (RFC 8887).
	if alt := httpOriginForWS(o); alt != "" {
		t.allow(alt)
	}
}

func httpOriginForWS(origin string) string {
	switch {
	case strings.HasPrefix(origin, "ws://"):
		return "http://" + strings.TrimPrefix(origin, "ws://")
	case strings.HasPrefix(origin, "wss://"):
		return "https://" + strings.TrimPrefix(origin, "wss://")
	default:
		return ""
	}
}
