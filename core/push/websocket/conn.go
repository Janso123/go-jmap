package websocket

import (
	"context"
	jsonv2 "encoding/json/v2"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core"
	cws "github.com/coder/websocket"
)

// Conn is an open JMAP-over-WebSocket connection (RFC 8887).
//
// With Options.Reconnect set, Conn redials after unexpected disconnect and
// re-sends WebSocketPushEnable with the last pushState when push was enabled.
// Close() disables reconnect. Binary frames are ignored. Blob upload/download
// remains on HTTPS (RFC 8620 §6 / RFC 8887 §4).
type Conn struct {
	client *jmap.Client
	wsURL  string
	opts   Options

	mu      sync.Mutex
	ws      *cws.Conn
	pending map[string]chan doResult
	handler func(*jmap.StateChange)

	pushState     string
	pushEnabled   bool
	pushDataTypes []jmap.EventType

	nextID atomic.Uint64
	gate   *requestGate

	readCancel context.CancelFunc
	closed     chan struct{}
	closeOnce  sync.Once
	userClosed bool
}

type doResult struct {
	resp *jmap.Response
	err  error
}

// requestGate limits concurrent Do calls (RFC 8887 §4.3.2).
type requestGate struct {
	ch chan struct{}
}

func newRequestGate(n uint64) *requestGate {
	if n == 0 {
		return nil
	}
	g := &requestGate{ch: make(chan struct{}, n)}
	for range n {
		g.ch <- struct{}{}
	}
	return g
}

func (g *requestGate) acquire(ctx context.Context) error {
	if g == nil {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-g.ch:
		return nil
	}
}

func (g *requestGate) release() {
	if g != nil {
		g.ch <- struct{}{}
	}
}

// Dial opens a WebSocket to the URL advertised in the session capability
// urn:ietf:params:jmap:websocket. Optional Options configure compression,
// concurrency limits, and auto-reconnect.
func Dial(ctx context.Context, client *jmap.Client, opts ...Options) (*Conn, error) {
	if client == nil {
		return nil, fmt.Errorf("websocket: nil client")
	}
	if err := ensureSession(ctx, client); err != nil {
		return nil, err
	}
	cap, ok := client.Session.Capabilities[URI].(*WebSocket)
	if !ok || cap.URL == "" {
		return nil, fmt.Errorf("websocket: session missing %s capability url", URI)
	}
	return DialURL(ctx, client, cap.URL, opts...)
}

// DialURL opens a WebSocket at wsURL with the jmap subprotocol.
func DialURL(ctx context.Context, client *jmap.Client, wsURL string, opts ...Options) (*Conn, error) {
	if client == nil {
		return nil, fmt.Errorf("websocket: nil client")
	}
	if wsURL == "" {
		return nil, fmt.Errorf("websocket: empty url")
	}
	if err := ensureSession(ctx, client); err != nil {
		return nil, err
	}

	o := mergeOptions(opts)
	c := &Conn{
		client:  client,
		wsURL:   wsURL,
		opts:    o,
		pending: make(map[string]chan doResult),
		closed:  make(chan struct{}),
		gate:    newRequestGate(resolveMaxConcurrent(client, o)),
	}
	if err := c.dialUnderlying(ctx); err != nil {
		return nil, err
	}
	c.startReadLoop()
	return c, nil
}

func resolveMaxConcurrent(client *jmap.Client, o Options) uint64 {
	if o.MaxConcurrentRequests > 0 {
		return o.MaxConcurrentRequests
	}
	if client.Session != nil {
		if coreCap, ok := client.Session.Capabilities[jmap.CoreURI].(*core.Core); ok {
			return coreCap.MaxConcurrentRequests
		}
	}
	return 0
}

func (c *Conn) dialUnderlying(ctx context.Context) error {
	hc := c.client.HttpClient
	if hc == nil {
		hc = http.DefaultClient
	}

	ws, resp, err := cws.Dial(ctx, c.wsURL, &cws.DialOptions{
		HTTPClient:      hc,
		Subprotocols:    []string{Subprotocol},
		CompressionMode: c.opts.CompressionMode,
	})
	if err != nil {
		return err
	}
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	if ws.Subprotocol() != Subprotocol {
		_ = ws.Close(cws.StatusProtocolError, "expected jmap subprotocol")
		return fmt.Errorf("websocket: server did not negotiate %q subprotocol (got %q)", Subprotocol, ws.Subprotocol())
	}

	c.mu.Lock()
	c.ws = ws
	c.mu.Unlock()
	return nil
}

func (c *Conn) startReadLoop() {
	readCtx, cancel := context.WithCancel(context.Background())
	c.mu.Lock()
	if c.readCancel != nil {
		c.readCancel()
	}
	c.readCancel = cancel
	c.mu.Unlock()
	go c.readLoop(readCtx)
	if c.opts.PingInterval > 0 {
		go c.pingLoop(readCtx)
	}
}

func (c *Conn) pingLoop(ctx context.Context) {
	ticker := time.NewTicker(c.opts.PingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.closed:
			return
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, c.opts.PingInterval)
			_ = c.Ping(pingCtx)
			cancel()
		}
	}
}

func ensureSession(ctx context.Context, client *jmap.Client) error {
	client.Lock()
	needAuth := client.Session == nil
	client.Unlock()
	if needAuth {
		return client.Authenticate(ctx)
	}
	return nil
}

// SetHandler sets the callback for StateChange push frames. It may be called
// before or after Dial; the handler is invoked from the read goroutine.
func (c *Conn) SetHandler(fn func(*jmap.StateChange)) {
	c.mu.Lock()
	c.handler = fn
	c.mu.Unlock()
}

// PushState returns the last pushState token received on this connection.
func (c *Conn) PushState() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pushState
}

// EnablePush sends WebSocketPushEnable. Pass nil dataTypes to subscribe to all
// types. pushState may be empty; when set, the server SHOULD send changes since
// that token (RFC 8887 §4.3.5.2).
func (c *Conn) EnablePush(dataTypes []jmap.EventType, pushState string) error {
	c.mu.Lock()
	c.pushEnabled = true
	c.pushDataTypes = append([]jmap.EventType(nil), dataTypes...)
	if pushState != "" {
		c.pushState = pushState
	}
	c.mu.Unlock()
	return c.writeJSON(PushEnable{
		Type:      "WebSocketPushEnable",
		DataTypes: dataTypes,
		PushState: pushState,
	})
}

// DisablePush sends WebSocketPushDisable.
func (c *Conn) DisablePush() error {
	c.mu.Lock()
	c.pushEnabled = false
	c.mu.Unlock()
	return c.writeJSON(PushDisable{Type: "WebSocketPushDisable"})
}

// Do sends a JMAP Request over the WebSocket and waits for the matching
// Response (or RequestError). Requests may complete out of order; correlation
// uses the RFC 8887 id / requestId fields.
func (c *Conn) Do(ctx context.Context, req *jmap.Request) (*jmap.Response, error) {
	if req == nil {
		return nil, fmt.Errorf("websocket: nil request")
	}

	found := slices.Contains(req.Using, jmap.CoreURI)
	if !found {
		req.Using = append(req.Using, jmap.CoreURI)
	}

	c.client.Lock()
	for _, uri := range req.Using {
		if _, ok := c.client.Session.RawCapabilities[uri]; !ok {
			c.client.Unlock()
			return nil, fmt.Errorf("server doesn't support required capability '%s'", uri)
		}
	}
	c.client.Unlock()

	if err := c.gate.acquire(ctx); err != nil {
		return nil, err
	}
	defer c.gate.release()

	id := strconv.FormatUint(c.nextID.Add(1), 10)
	ch := make(chan doResult, 1)

	c.mu.Lock()
	if c.userClosed {
		c.mu.Unlock()
		return nil, fmt.Errorf("websocket: connection closed")
	}
	if c.ws == nil {
		c.mu.Unlock()
		return nil, fmt.Errorf("websocket: connection closed")
	}
	c.pending[id] = ch
	c.mu.Unlock()

	raw, err := marshalRequest(id, req)
	if err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, err
	}
	if err := c.writeRaw(ctx, raw); err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, err
	}

	select {
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, ctx.Err()
	case <-c.closed:
		return nil, fmt.Errorf("websocket: connection closed")
	case res := <-ch:
		return res.resp, res.err
	}
}

// Ping sends a WebSocket ping and waits for a pong (RFC 6455).
// Must be used while the connection read loop is running.
func (c *Conn) Ping(ctx context.Context) error {
	c.mu.Lock()
	ws := c.ws
	userClosed := c.userClosed
	c.mu.Unlock()
	if userClosed || ws == nil {
		return fmt.Errorf("websocket: connection closed")
	}
	return ws.Ping(ctx)
}

// Close performs a normal WebSocket close handshake, cancels in-flight Do calls,
// and disables auto-reconnect.
func (c *Conn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		c.mu.Lock()
		c.userClosed = true
		if c.readCancel != nil {
			c.readCancel()
		}
		ws := c.ws
		c.ws = nil
		for id, ch := range c.pending {
			ch <- doResult{err: fmt.Errorf("websocket: connection closed")}
			delete(c.pending, id)
		}
		c.mu.Unlock()
		if ws != nil {
			err = ws.Close(cws.StatusNormalClosure, "")
		}
		close(c.closed)
	})
	return err
}

func (c *Conn) writeJSON(v any) error {
	raw, err := jsonv2.Marshal(v)
	if err != nil {
		return err
	}
	return c.writeRaw(context.Background(), raw)
}

func (c *Conn) writeRaw(ctx context.Context, raw []byte) error {
	c.mu.Lock()
	ws := c.ws
	userClosed := c.userClosed
	c.mu.Unlock()
	if userClosed || ws == nil {
		return fmt.Errorf("websocket: connection closed")
	}
	return ws.Write(ctx, cws.MessageText, raw)
}

func (c *Conn) readLoop(ctx context.Context) {
	var readErr error
	for {
		c.mu.Lock()
		ws := c.ws
		c.mu.Unlock()
		if ws == nil {
			readErr = fmt.Errorf("websocket: connection closed")
			break
		}
		typ, data, err := ws.Read(ctx)
		if err != nil {
			readErr = err
			break
		}
		if typ != cws.MessageText {
			continue
		}
		fr, err := decodeServerFrame(data)
		if err != nil {
			if c.opts.OnFrameError != nil {
				c.opts.OnFrameError(err, data)
			}
			// Unblock a matching pending Do when requestId is present so an
			// unknown/@type or decode error cannot leave Do hanging forever.
			var probe struct {
				RequestID string `json:"requestId"`
			}
			if jsonv2.Unmarshal(data, &probe) == nil && probe.RequestID != "" {
				c.deliver(probe.RequestID, doResult{err: err})
			}
			continue
		}
		switch {
		case fr.StateChange != nil:
			c.mu.Lock()
			if fr.StateChange.PushState != "" {
				c.pushState = fr.StateChange.PushState
			}
			h := c.handler
			c.mu.Unlock()
			if h != nil {
				h(fr.StateChange)
			}
		case fr.CalendarAlert != nil:
			if c.opts.OnCalendarAlert != nil {
				c.opts.OnCalendarAlert(fr.CalendarAlert)
			}
		case fr.Response != nil:
			c.deliver(fr.RequestID, doResult{resp: fr.Response})
		case fr.RequestError != nil:
			c.deliver(fr.RequestID, doResult{err: fr.RequestError})
		}
	}

	c.mu.Lock()
	if c.ws != nil {
		_ = c.ws.CloseNow()
		c.ws = nil
	}
	userClosed := c.userClosed
	wantReconnect := c.opts.Reconnect != nil && !userClosed
	var onDisc func(error)
	if c.opts.Reconnect != nil {
		onDisc = c.opts.Reconnect.OnDisconnect
	}
	c.mu.Unlock()

	c.failPending(fmt.Errorf("websocket: connection closed"))

	if !wantReconnect {
		_ = c.Close()
		return
	}

	if onDisc != nil {
		onDisc(readErr)
	}

	backoff := c.opts.Reconnect.minBackoff()
	maxBackoff := c.opts.Reconnect.maxBackoff()
	for {
		c.mu.Lock()
		if c.userClosed {
			c.mu.Unlock()
			return
		}
		c.mu.Unlock()

		timer := time.NewTimer(backoff)
		select {
		case <-c.closed:
			timer.Stop()
			return
		case <-timer.C:
		}

		dialCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := c.dialUnderlying(dialCtx)
		cancel()
		if err != nil {
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		c.mu.Lock()
		pushEnabled := c.pushEnabled
		dataTypes := append([]jmap.EventType(nil), c.pushDataTypes...)
		pushState := c.pushState
		var onRe func()
		if c.opts.Reconnect != nil {
			onRe = c.opts.Reconnect.OnReconnect
		}
		c.mu.Unlock()

		if pushEnabled {
			_ = c.writeJSON(PushEnable{
				Type:      "WebSocketPushEnable",
				DataTypes: dataTypes,
				PushState: pushState,
			})
		}
		if onRe != nil {
			onRe()
		}
		c.startReadLoop()
		return
	}
}

func (c *Conn) failPending(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for id, ch := range c.pending {
		ch <- doResult{err: err}
		delete(c.pending, id)
	}
}

func (c *Conn) deliver(requestID string, res doResult) {
	c.mu.Lock()
	ch, ok := c.pending[requestID]
	if ok {
		delete(c.pending, requestID)
	}
	c.mu.Unlock()
	if ok {
		ch <- res
	}
}
