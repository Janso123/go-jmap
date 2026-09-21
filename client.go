package jmap

import (
	"bytes"
	"context"
	"encoding/base64"
	jsonv2 "encoding/json/v2"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

	"golang.org/x/oauth2"
)

// A JMAP Client
type Client struct {
	sync.Mutex
	// HttpClient is used for requests. It should typically handle authentication.
	// WithBasicAuth and WithAccessToken wrap the existing HttpClient transport
	// with oauth2 while preserving Timeout, CheckRedirect, and Jar. If nil,
	// http.DefaultClient is used via httpClient().
	HttpClient *http.Client

	// SessionEndpoint is the JMAP Session Resource URL (RFC 8620 §2).
	// Do marks the session stale when a response's sessionState differs from
	// the cached Session; it does not refetch. Check SessionStale and call
	// RefreshSession(ctx) to reload apiUrl / capabilities.
	SessionEndpoint string

	// UserAgent is sent on Authenticate and Do requests. NewClient sets a
	// default of "go-jmap/" + Version.
	UserAgent string

	// Session is the cached JMAP Session object from Authenticate / RefreshSession.
	Session *Session

	sessionStale bool

	// OnSessionChange is called after a successful RefreshSession when a
	// previous Session existed. old is the Session before refresh; new is
	// the Session after Authenticate.
	OnSessionChange func(old, new *Session)
}

func (c *Client) httpClient() *http.Client {
	if c.HttpClient == nil {
		return http.DefaultClient
	}
	return c.HttpClient
}

// HTTPClient returns the HTTP client used for JMAP requests, or
// http.DefaultClient when none is set.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient()
}

func (c *Client) userAgent() string {
	if c.UserAgent != "" {
		return c.UserAgent
	}
	return "go-jmap/" + Version
}

// EffectiveUserAgent returns Client.UserAgent, or "go-jmap/"+Version when empty.
func (c *Client) EffectiveUserAgent() string {
	return c.userAgent()
}

func authHeader(tok *oauth2.Token) string {
	typ := tok.TokenType
	switch {
	case strings.EqualFold(typ, "bearer"), typ == "":
		typ = "Bearer"
	case strings.EqualFold(typ, "basic"):
		typ = "Basic"
	}
	return typ + " " + tok.AccessToken
}

func (c *Client) applyToken(tok *oauth2.Token) {
	hc := c.HttpClient
	if hc == nil || hc == http.DefaultClient {
		hc = &http.Client{}
	} else {
		clone := *hc
		hc = &clone
	}
	t := &originAuthTransport{
		base:   unwrapOriginAuth(hc.Transport),
		header: authHeader(tok),
	}
	t.allow(hostOf(c.SessionEndpoint))
	hc.Transport = t
	c.HttpClient = hc
}

func (c *Client) allowAuthHosts(s *Session) {
	if c.HttpClient == nil {
		return
	}
	t, ok := c.HttpClient.Transport.(*originAuthTransport)
	if !ok {
		return
	}
	t.allow(hostOf(c.SessionEndpoint))
	if s == nil {
		return
	}
	t.allow(hostOf(s.APIURL))
	t.allow(hostOf(s.UploadURL))
	t.allow(hostOf(s.DownloadURL))
	t.allow(hostOf(s.EventSourceURL))
}

func resolveRef(base *url.URL, ref string) string {
	if ref == "" || base == nil {
		return ref
	}
	u, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return base.ResolveReference(u).String()
}

func (c *Client) resolveSessionURLs(s *Session) {
	base, err := url.Parse(c.SessionEndpoint)
	if err != nil {
		return
	}
	s.APIURL = resolveRef(base, s.APIURL)
	s.UploadURL = resolveRef(base, s.UploadURL)
	s.DownloadURL = resolveRef(base, s.DownloadURL)
	s.EventSourceURL = resolveRef(base, s.EventSourceURL)
}

// WithBasicAuth configures HttpClient to send HTTP Basic auth via oauth2.Transport.
// Existing Timeout, CheckRedirect, Jar, and Transport.Base are preserved.
func (c *Client) WithBasicAuth(username string, password string) *Client {
	auth := username + ":" + password
	c.applyToken(&oauth2.Token{
		AccessToken: base64.StdEncoding.EncodeToString([]byte(auth)),
		TokenType:   "basic",
	})
	return c
}

// WithAccessToken configures HttpClient to send a Bearer token via oauth2.Transport.
// Existing Timeout, CheckRedirect, Jar, and Transport.Base are preserved.
func (c *Client) WithAccessToken(token string) *Client {
	c.applyToken(&oauth2.Token{AccessToken: token, TokenType: "bearer"})
	return c
}

// PrimaryAccount returns the primary account ID for the given capability URI.
func (c *Client) PrimaryAccount(uri URI) (ID, error) {
	if c.Session == nil {
		return "", fmt.Errorf("session not loaded")
	}
	id, ok := c.Session.PrimaryAccounts[uri]
	if !ok || id == "" {
		return "", fmt.Errorf("no primary account for %s", uri)
	}
	return id, nil
}

// Authenticate authenticates the client and retrieves the Session object.
// Authenticate will be called automatically when Do is called if the Session
// object hasn't already been initialized. Call Authenticate before any requests
// if you need to access information from the Session object prior to the first
// request
func (c *Client) Authenticate(ctx context.Context) error {
	c.Lock()
	if c.SessionEndpoint == "" {
		c.Unlock()
		return fmt.Errorf("no session url is set")
	}
	endpoint := c.SessionEndpoint
	c.Unlock()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent())

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return decodeHttpError(resp)
	}

	data, err := readJSONBody(resp.Body)
	if err != nil {
		return err
	}

	s := &Session{}
	err = jsonv2.Unmarshal(data, s)
	if err != nil {
		return fmt.Errorf("jmap: decoding session: %w", err)
	}
	c.resolveSessionURLs(s)

	c.Lock()
	c.Session = s
	c.sessionStale = false
	c.allowAuthHosts(s)
	c.Unlock()

	return nil
}

// SessionStale reports whether the last Do response indicated a different
// sessionState than the cached Session. The client does not refresh
// automatically; call RefreshSession when this is true.
func (c *Client) SessionStale() bool {
	c.Lock()
	defer c.Unlock()
	return c.sessionStale
}

// ObserveSessionState marks the cached Session stale when state is non-empty
// and differs from Session.State. HTTPS Do and WebSocket Response frames call this.
func (c *Client) ObserveSessionState(state string) {
	if c == nil {
		return
	}
	c.Lock()
	defer c.Unlock()
	if c.Session != nil && state != "" && state != c.Session.State {
		c.sessionStale = true
	}
}

// RefreshSession re-fetches the Session via Authenticate, clears the stale
// flag, and invokes OnSessionChange when a previous Session existed.
// Call this when SessionStale is true (or whenever you need a fresh Session).
func (c *Client) RefreshSession(ctx context.Context) error {
	c.Lock()
	old := c.Session
	cb := c.OnSessionChange
	c.Unlock()
	if err := c.Authenticate(ctx); err != nil {
		return err
	}
	c.Lock()
	c.sessionStale = false
	neu := c.Session
	c.Unlock()
	if cb != nil && old != nil {
		cb(old, neu)
	}
	return nil
}

// The core capabilty must be included in all method calls
const CoreURI URI = "urn:ietf:params:jmap:core"

// Do performs a JMAP request and returns the response. Prefer the ctx
// parameter for cancellation; if req.Context is also set and ctx is
// background-only, req.Context is used as a fallback.
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	c.Lock()
	if c.Session == nil {
		c.Unlock()
		err := c.Authenticate(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		c.Unlock()
	}
	using := append([]URI(nil), req.Using...)
	if !slices.Contains(using, CoreURI) {
		using = append(using, CoreURI)
	}

	// Check the required capabilities before making the request
	c.Lock()
	for _, uri := range using {
		// Check RawCapabilities in case we have asked for unparsed
		// capabilities, or the core capability
		_, ok := c.Session.RawCapabilities[uri]
		if !ok {
			c.Unlock()
			return nil, fmt.Errorf("server doesn't support required capability '%s'", uri)
		}
	}
	apiURL := c.Session.APIURL
	c.Unlock()

	wire := *req
	wire.Using = using
	body, err := jsonv2.Marshal(&wire)
	if err != nil {
		return nil, err
	}
	reqCtx := ctx
	if reqCtx == nil {
		if req.Context != nil {
			reqCtx = req.Context
		} else {
			reqCtx = context.Background()
		}
	}
	httpReq, err := http.NewRequestWithContext(reqCtx, "POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", c.userAgent())

	httpResp, err := c.httpClient().Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return nil, decodeHttpError(httpResp)
	}

	data, err := readJSONBody(httpResp.Body)
	if err != nil {
		return nil, err
	}
	resp := &Response{}
	err = jsonv2.Unmarshal(data, resp)
	if err != nil {
		return nil, fmt.Errorf("jmap: decoding response: %w", err)
	}

	c.ObserveSessionState(resp.SessionState)

	return resp, nil
}

// DownloadOptions controls URL template substitution for Download.
// Empty Name defaults to "filename"; empty Type defaults to
// "application/octet-stream".
type DownloadOptions struct {
	Name string
	Type string
}

// Upload sends binary data to the server and returns blob ID and some
// associated meta-data.
//
// contentType is sent as the request Content-Type. If empty, it defaults
// to "application/octet-stream".
//
// There are some caveats to keep in mind:
// - Server may return the same blob ID for multiple uploads of the same blob.
// - Blob ID may become invalid after some time if it is unused.
// - Blob ID is usable only by the uploader until it is used, even for shared accounts.
func (c *Client) Upload(ctx context.Context, accountID ID, blob io.Reader, contentType string) (*UploadResponse, error) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Lock()
	if c.SessionEndpoint == "" {
		c.Unlock()
		return nil, fmt.Errorf("jmap/client: SessionEndpoint is empty")
	}
	if c.Session == nil {
		c.Unlock()
		err := c.Authenticate(ctx)
		if err != nil {
			return nil, err
		}

		c.Lock()
	}

	uploadURL := ExpandURITemplateLevel1(c.Session.UploadURL, map[string]string{
		"accountId": string(accountID),
	})
	c.Unlock()
	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, blob)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent())

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, decodeHttpError(resp)
	}

	data, err := readJSONBody(resp.Body)
	if err != nil {
		return nil, err
	}

	info := &UploadResponse{}
	err = jsonv2.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func expandDownloadURL(tmpl, accountID, blobID, typ, name string) string {
	return ExpandURITemplateLevel1(tmpl, map[string]string{
		"accountId": accountID,
		"blobId":    blobID,
		"type":      typ,
		"name":      name,
	})
}

const maxJSONBody = 32 << 20

func readJSONBody(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxJSONBody+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxJSONBody {
		return nil, fmt.Errorf("jmap: response body exceeds %d bytes", maxJSONBody)
	}
	return data, nil
}

// Download downloads binary data by its Blob ID from the server.
// opts.Name and opts.Type are substituted into the session downloadUrl
// template ({name}, {type}); see DownloadOptions for defaults.
func (c *Client) Download(ctx context.Context, accountID ID, blobID ID, opts DownloadOptions) (io.ReadCloser, error) {
	name := opts.Name
	if name == "" {
		name = "filename"
	}
	typ := opts.Type
	if typ == "" {
		typ = "application/octet-stream"
	}

	c.Lock()
	if c.SessionEndpoint == "" {
		c.Unlock()
		return nil, fmt.Errorf("jmap/client: SessionEndpoint is empty")
	}
	if c.Session == nil {
		c.Unlock()
		err := c.Authenticate(ctx)
		if err != nil {
			return nil, err
		}

		c.Lock()
	}

	tgtUrl := expandDownloadURL(c.Session.DownloadURL, string(accountID), string(blobID), typ, name)
	c.Unlock()
	req, err := http.NewRequestWithContext(ctx, "GET", tgtUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent())

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		return nil, decodeHttpError(resp)
	}

	return resp.Body, nil
}

func decodeHttpError(resp *http.Response) error {
	const maxBody = 4096
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	truncated := len(body) > maxBody
	if truncated {
		body = body[:maxBody]
	}
	mt, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if mt == "application/problem+json" || mt == "application/json" {
		var re RequestError
		if err := jsonv2.Unmarshal(body, &re); err == nil {
			if re.Type == "" {
				re.Type = "about:blank"
			}
			re.Status = resp.StatusCode
			return &re
		}
	}
	return &HTTPError{
		Status:      resp.StatusCode,
		StatusText:  resp.Status,
		Body:        strings.TrimSpace(string(body)),
		ContentType: mt,
	}
}

// UploadResponse is the object returned in response to blob upload.
type UploadResponse struct {
	// The id of the account used for the call.
	Account ID `json:"accountId"`

	// The id representing the binary data uploaded. The data for this id is
	// immutable. The id only refers to the binary data, not any metadata.
	ID ID `json:"blobId"`

	// The media type of the file (as specified in RFC 6838, section 4.2) as
	// set in the Content-Type header of the upload HTTP request.
	Type string `json:"type"`

	// The size of the file in octets.
	Size uint64 `json:"size"`
}
