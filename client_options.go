package jmap

import (
	"fmt"
	"net/http"
	"time"
)

// Option configures a Client created by NewClient.
//
// Option order matters: WithBearer and WithBasic replace HttpClient entirely.
// Apply auth before WithTimeout / WithTrustedHosts so those mutate the auth
// client. WithHTTPClient after WithBearer/WithBasic overwrites auth; WithBearer
// / WithBasic after WithHTTPClient discards the custom client. Prefer either a
// pre-authenticated *http.Client via WithHTTPClient alone, or auth first then
// timeout/hosts.
type Option func(*Client)

// NewClient creates a Client for the given session URL with optional
// configuration. The default User-Agent is "go-jmap/" + Version.
//
// Options are applied in order. WithBearer and WithBasic replace HttpClient;
// see Option for composition rules (auth before timeout/trusted hosts; do not
// interleave WithHTTPClient with auth unless intentional).
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
// Applied after WithBearer/WithBasic replaces the auth client; applied before
// them is discarded when auth runs. Prefer a client that already carries auth,
// or use WithBearer/WithBasic without a prior WithHTTPClient.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.HttpClient = h }
}

// WithBearer configures bearer token authentication via WithAccessToken.
// Replaces HttpClient; apply before WithTimeout / WithTrustedHosts.
func WithBearer(token string) Option {
	return func(c *Client) { c.WithAccessToken(token) }
}

// WithBasic configures HTTP basic authentication via WithBasicAuth.
// Replaces HttpClient; apply before WithTimeout / WithTrustedHosts.
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
