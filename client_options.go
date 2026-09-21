package jmap

import (
	"fmt"
	"net/http"
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

// WithTrustedHosts restricts HTTP redirects to the given hosts (req.URL.Host).
// Redirect chains longer than 5 are rejected.
func WithTrustedHosts(hosts ...string) Option {
	return func(c *Client) {
		allowed := map[string]struct{}{}
		for _, h := range hosts {
			allowed[h] = struct{}{}
		}
		hc := c.httpClient()
		if hc == http.DefaultClient {
			hc = &http.Client{}
			c.HttpClient = hc
		}
		hc.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("stopped after 5 redirects")
			}
			if _, ok := allowed[req.URL.Host]; !ok {
				return fmt.Errorf("redirect to untrusted host %q", req.URL.Host)
			}
			return nil
		}
	}
}

// WithUserAgent overrides the default User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.UserAgent = ua }
}
