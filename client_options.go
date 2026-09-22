package jmap

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Option configures a Client created by NewClient.
//
// WithBearer and WithBasic preserve Timeout, CheckRedirect, and Jar on the
// current HttpClient. WithHTTPClient after auth replaces the entire client
// (including auth transport). WithBearer/WithBasic after WithHTTPClient wrap
// the custom client's transport. Option order no longer drops Timeout.
type Option func(*Client)

// NewClient creates a Client for the given session URL with optional
// configuration. The default User-Agent is "go-jmap/" + Version.
//
// Options are applied in order. Auth options preserve existing HttpClient
// settings; see Option for WithHTTPClient interaction.
func NewClient(sessionURL string, opts ...Option) *Client {
	c := &Client{
		SessionEndpoint: sessionURL,
		UserAgent:       "go-jmap/" + Version,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithHTTPClient sets the HTTP client used for requests.
// WithBearer/WithBasic wrap the current Transport for auth; Timeout,
// CheckRedirect, and Jar on that client are preserved. WithHTTPClient after
// auth replaces the entire HttpClient (including the auth transport).
// WithBearer/WithBasic after WithHTTPClient wrap the custom client's transport.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.HttpClient = h }
}

// WithBearer configures bearer token authentication via WithAccessToken.
func WithBearer(token string) Option {
	return func(c *Client) { c.WithAccessToken(token) }
}

// WithBasic configures HTTP basic authentication via WithBasicAuth.
func WithBasic(user, pass string) Option {
	return func(c *Client) { c.WithBasicAuth(user, pass) }
}

// WithTimeout sets the Timeout on the client's HTTP client. If HttpClient is
// nil or DefaultClient, a new http.Client is allocated first.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		hc := c.httpClient()
		if hc == http.DefaultClient {
			hc = &http.Client{}
			c.HttpClient = hc
		}
		hc.Timeout = d
	}
}

// WithTrustedHosts restricts HTTP redirects to the given hosts or origins.
// A host (or host:port) matches only when the redirect keeps the original
// request scheme and that host:port. A full origin (scheme://host[:port])
// matches that origin. Scheme downgrades (https→http) are rejected.
// Redirect chains longer than 5 are rejected. With no list, redirects must
// stay on the original origin (scheme+host+port).
func WithTrustedHosts(hosts ...string) Option {
	return func(c *Client) {
		c.trustedHosts = append([]string(nil), hosts...)
		c.policyClient = nil
	}
}

func enforceRedirectOrigin(req *http.Request, via []*http.Request, allowed []string) error {
	host := ""
	if req != nil && req.URL != nil {
		host = req.URL.Host
	}
	if len(via) >= 5 {
		return fmt.Errorf("stopped after 5 redirects")
	}
	if req == nil || req.URL == nil || len(via) == 0 || via[0] == nil || via[0].URL == nil {
		return fmt.Errorf("redirect to untrusted host %q", host)
	}
	got := originOfURL(req.URL)
	orig := originOfURL(via[0].URL)
	if got == "" || orig == "" || schemeDowngrade(via[0].URL.Scheme, req.URL.Scheme) {
		return fmt.Errorf("redirect to untrusted host %q", host)
	}
	if len(allowed) == 0 {
		if got != orig {
			return fmt.Errorf("redirect to untrusted host %q", host)
		}
		return nil
	}
	for _, entry := range allowed {
		if allowedEntryOrigin(entry, via[0].URL) == got {
			return nil
		}
	}
	return fmt.Errorf("redirect to untrusted host %q", host)
}

func schemeDowngrade(from, to string) bool {
	from, to = strings.ToLower(from), strings.ToLower(to)
	return (from == "https" || from == "wss") && (to == "http" || to == "ws")
}

func allowedEntryOrigin(entry string, orig *url.URL) string {
	if strings.Contains(entry, "://") {
		return originOf(entry)
	}
	if orig == nil || orig.Scheme == "" {
		return ""
	}
	return originOfURL(&url.URL{Scheme: orig.Scheme, Host: entry})
}

// WithUserAgent overrides the default User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.UserAgent = ua }
}

// WithMaxResponseBytes limits each buffered JSON response to n bytes.
// A positive n replaces the default cap. When n is zero or this option is
// omitted, the cap is max(32<<20, the session core maxSizeRequest when that
// capability is present). Using maxSizeRequest as a response limit is a
// heuristic (RFC 8620 §2 defines the value for requests).
func WithMaxResponseBytes(n int64) Option {
	return func(c *Client) { c.maxResponseBytes = n }
}
