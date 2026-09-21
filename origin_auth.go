package jmap

import (
	"net/http"
	"net/url"
	"sync"
)

// originAuthTransport adds Authorization only for requests whose host is in
// the allow-set. Cross-origin redirects therefore do not receive credentials
// (RFC 9110 §15.4 / RFC 8620 §8.3).
type originAuthTransport struct {
	base   http.RoundTripper
	header string

	mu    sync.RWMutex
	hosts map[string]struct{}
}

func (t *originAuthTransport) allow(host string) {
	if host == "" {
		return
	}
	t.mu.Lock()
	if t.hosts == nil {
		t.hosts = make(map[string]struct{})
	}
	t.hosts[host] = struct{}{}
	t.mu.Unlock()
}

func (t *originAuthTransport) allowed(host string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.hosts[host]
	return ok
}

func (t *originAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	if t.allowed(req.URL.Host) {
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

func hostOf(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return u.Host
}

// AllowAuthOrigin adds the host of rawURL to the credential allow-set used by
// WithBearer/WithBasic. Call this before connecting to a WebSocket URL whose
// host is not the session origin.
func (c *Client) AllowAuthOrigin(rawURL string) {
	if c == nil || c.HttpClient == nil {
		return
	}
	t, ok := c.HttpClient.Transport.(*originAuthTransport)
	if !ok {
		return
	}
	t.allow(hostOf(rawURL))
}
