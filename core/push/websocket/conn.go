package websocket

import (
	"context"
	jsonv2 "encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/core"
	cws "github.com/coder/websocket"
)

// Conn is an open JMAP-over-WebSocket connection (RFC 8887).
//
// With Options.Reconnect set, Conn redials after unexpected disconnect and
// re-sends WebSocketPushEnable with the last pushState when push was enabled.
// Close() disables reconnect. Binary frames are reported to Options.OnFrameError
// and otherwise ignored (RFC 8887 §4.2). Blob upload/download
// remains on HTTPS (RFC 8620 §6 / RFC 8887 §4).
type Conn struct {
	client *jmap.Client
	wsURL  string
	opts   Options

	mu      sync.Mutex
	ws      *cws.Conn
	pending map[string]chan doResult
	handler func(*jmap.StateChange)

	pushState        string
	pushWanted       bool
	pushEnabled      bool
	pushDataTypes    []jmap.EventType
	pushDataTypesNil bool

	// reconnected is closed to wake Do calls waiting out a reconnect.
	// Waiters snapshot it under mu; broadcast replaces it with a new channel.
	reconnected chan struct{}

	nextID atomic.Uint64
	gate   *requestGate

	readCancel context.CancelFunc
	closed     chan struct{}
	closeOnce  sync.Once
	userClosed bool

	// events hands push frames to dispatchLoop so handlers never run on the
	// read goroutine and may therefore call Do. Buffered 1; a full buffer
	// means the stale event is dropped in favour of the newer one.
	events chan wsEvent
}

type doResult struct {
	resp *jmap.Response
	err  error
}

// wsEvent is one push frame queued for the dispatcher goroutine.
type wsEvent struct {
	state *jmap.StateChange
	alert *calendar.CalendarAlert
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
	if err := client.CheckWebSocketURL(wsURL, o.AllowForeignOrigin); err != nil {
		return nil, err
	}
	client.AllowAuthOrigin(wsURL)

	c := &Conn{
		client:      client,
		wsURL:       wsURL,
		opts:        o,
		pending:     make(map[string]chan doResult),
		closed:      make(chan struct{}),
		events:      make(chan wsEvent, 1),
		reconnected: make(chan struct{}),
		gate:        newRequestGate(resolveMaxConcurrent(client, o)),
	}
	if err := c.dialUnderlying(ctx); err != nil {
		return nil, err
	}
	go c.dispatchLoop()
	c.startReadLoop()
	return c, nil
}

func resolveMaxConcurrent(client *jmap.Client, o Options) uint64 {
	if o.MaxConcurrentRequests > 0 {
		return o.MaxConcurrentRequests
	}
	if client.Session != nil {
		if coreCap, ok := client.Session.Capabilities[jmap.CoreURI].(*core.Core); ok {
			return uint64(coreCap.MaxConcurrentRequests)
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
	if err := jmap.WebSocketResponseOriginOK(c.wsURL, resp); err != nil {
		_ = ws.Close(cws.StatusPolicyViolation, "origin")
		return err
	}
	if ws.Subprotocol() != Subprotocol {
		_ = ws.Close(cws.StatusProtocolError, "expected jmap subprotocol")
		return fmt.Errorf("websocket: server did not negotiate %q subprotocol (got %q)", Subprotocol, ws.Subprotocol())
	}

	limit := c.opts.ReadLimit
	if limit == 0 {
		limit = 32 << 20
	}
	ws.SetReadLimit(limit)

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
// before or after Dial; the handler is invoked from a dispatcher goroutine, so
// it may call Do.
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

// ErrPushUnsupported is returned by EnablePush when the session advertises
// urn:ietf:params:jmap:websocket with supportsPush false.
var ErrPushUnsupported = errors.New("websocket: push not supported")

// EnablePush sends WebSocketPushEnable. Pass nil dataTypes to subscribe to all
// types. pushState may be empty; when set, the server SHOULD send changes since
// that token (RFC 8887 §4.3.5.2).
//
// EnablePush returns ErrPushUnsupported when the websocket capability is
// present and SupportsPush is false. pushEnabled becomes true only after the
// frame is written.
func (c *Conn) EnablePush(ctx context.Context, dataTypes []jmap.EventType, pushState string) error {
	if err := c.errIfPushUnsupported(); err != nil {
		return err
	}
	c.mu.Lock()
	c.pushWanted = true
	c.pushDataTypes = append([]jmap.EventType(nil), dataTypes...)
	c.pushDataTypesNil = dataTypes == nil
	if pushState != "" {
		c.pushState = pushState
	}
	c.mu.Unlock()
	err := c.writeJSON(ctx, pushEnableFrame(dataTypes, pushState))
	c.mu.Lock()
	c.pushEnabled = err == nil && c.ws != nil
	c.mu.Unlock()
	return err
}

func (c *Conn) errIfPushUnsupported() error {
	if c.client == nil {
		return nil
	}
	c.client.Lock()
	defer c.client.Unlock()
	if c.client.Session == nil {
		return nil
	}
	cap, ok := c.client.Session.Capabilities[URI].(*WebSocket)
	if !ok {
		return nil
	}
	if !cap.SupportsPush {
		return ErrPushUnsupported
	}
	return nil
}

// DisablePush sends WebSocketPushDisable.
func (c *Conn) DisablePush(ctx context.Context) error {
	c.mu.Lock()
	c.pushWanted = false
	c.pushEnabled = false
	c.mu.Unlock()
	return c.writeJSON(ctx, PushDisable{Type: "WebSocketPushDisable"})
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
		using := append([]jmap.URI(nil), req.Using...)
		using = append(using, jmap.CoreURI)
		cp := *req
		cp.Using = using
		req = &cp
	}

	c.client.Lock()
	for _, uri := range req.Using {
		if _, ok := c.client.Session.RawCapabilities[uri]; !ok {
			c.client.Unlock()
			return nil, fmt.Errorf("server doesn't support required capability '%s'", uri)
		}
	}
	c.client.Unlock()

	if err := c.waitReady(ctx); err != nil {
		return nil, err
	}

	if err := c.gate.acquire(ctx); err != nil {
		return nil, err
	}
	defer c.gate.release()

	id := strconv.FormatUint(c.nextID.Add(1), 10)
	ch := make(chan doResult, 1)

	for {
		c.mu.Lock()
		if c.userClosed {
			c.mu.Unlock()
			return nil, fmt.Errorf("websocket: connection closed")
		}
		if c.ws == nil {
			c.mu.Unlock()
			if err := c.waitReady(ctx); err != nil {
				return nil, err
			}
			continue
		}
		c.pending[id] = ch
		c.mu.Unlock()
		break
	}

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

func (c *Conn) writeJSON(ctx context.Context, v any) error {
	raw, err := jsonv2.Marshal(v)
	if err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return c.writeRaw(ctx, raw)
}

// waitReady blocks while auto-reconnect is in progress. A nil socket with
// Reconnect unset, or Close, returns immediately.
func (c *Conn) waitReady(ctx context.Context) error {
	for {
		c.mu.Lock()
		userClosed := c.userClosed
		wsUp := c.ws != nil
		reconnect := c.opts.Reconnect != nil && !userClosed
		wait := c.reconnected
		c.mu.Unlock()
		if userClosed {
			return fmt.Errorf("websocket: connection closed")
		}
		if wsUp {
			return nil
		}
		if !reconnect {
			return fmt.Errorf("websocket: connection closed")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.closed:
			return fmt.Errorf("websocket: connection closed")
		case <-wait:
		}
	}
}

func (c *Conn) broadcastReconnected() {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch := c.reconnected
	c.reconnected = make(chan struct{})
	if ch != nil {
		close(ch)
	}
}

func (c *Conn) writeRaw(ctx context.Context, raw []byte) error {
	c.mu.Lock()
	ws := c.ws
	userClosed := c.userClosed
	c.mu.Unlock()
	if userClosed || ws == nil {
		return fmt.Errorf("websocket: connection closed")
	}
	if c.opts.WriteTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.opts.WriteTimeout)
		defer cancel()
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
			if c.opts.OnFrameError != nil {
				c.opts.OnFrameError(errors.New("websocket: binary frame (RFC 8887 §4.2 requires text)"), data)
			}
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
				c.deliver(probe.RequestID, doResult{err: err}, data)
			}
			continue
		}
		switch {
		case fr.StateChange != nil:
			c.mu.Lock()
			if fr.StateChange.PushState != "" {
				c.pushState = fr.StateChange.PushState
			}
			c.mu.Unlock()
			c.dispatch(wsEvent{state: fr.StateChange})
		case fr.CalendarAlert != nil:
			c.dispatch(wsEvent{alert: fr.CalendarAlert})
		case fr.Response != nil:
			c.client.ObserveSessionState(fr.Response.SessionState)
			if fr.RequestID == "" {
				// No id to correlate with: report it, but leave every other
				// in-flight Do alone.
				if c.opts.OnFrameError != nil {
					c.opts.OnFrameError(fmt.Errorf("websocket: Response missing requestId"), data)
				}
				break
			}
			c.deliver(fr.RequestID, doResult{resp: fr.Response}, data)
		case fr.RequestError != nil:
			if fr.RequestID == "" {
				// A null requestId is only unambiguous with a single in-flight
				// request; otherwise it cannot be attributed to any of them.
				if c.deliverToSolePending(doResult{err: fr.RequestError}) {
					break
				}
				if c.opts.OnFrameError != nil {
					c.opts.OnFrameError(fr.RequestError, data)
				}
				break
			}
			c.deliver(fr.RequestID, doResult{err: fr.RequestError}, data)
		}
	}

	c.mu.Lock()
	if c.ws != nil {
		_ = c.ws.CloseNow()
		c.ws = nil
	}
	c.pushEnabled = false
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

		c.mu.Lock()
		if c.userClosed {
			c.mu.Unlock()
			return
		}
		c.mu.Unlock()

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
		if c.userClosed {
			ws := c.ws
			c.ws = nil
			c.mu.Unlock()
			if ws != nil {
				_ = ws.CloseNow()
			}
			return
		}
		pushWanted := c.pushWanted
		dataTypes := append([]jmap.EventType(nil), c.pushDataTypes...)
		if c.pushDataTypesNil {
			dataTypes = nil
		}
		pushState := c.pushState
		var onRe func()
		if c.opts.Reconnect != nil {
			onRe = c.opts.Reconnect.OnReconnect
		}
		c.mu.Unlock()

		if pushWanted {
			if err := c.writeJSON(context.Background(), pushEnableFrame(dataTypes, pushState)); err != nil {
				c.mu.Lock()
				c.pushEnabled = false
				ws := c.ws
				c.ws = nil
				c.mu.Unlock()
				if ws != nil {
					_ = ws.CloseNow()
				}
				if onDisc != nil {
					onDisc(err)
				}
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				continue
			}
			c.mu.Lock()
			c.pushEnabled = true
			c.mu.Unlock()
		}
		if onRe != nil {
			onRe()
		}
		c.startReadLoop()
		c.broadcastReconnected()
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

func (c *Conn) deliver(requestID string, res doResult, raw []byte) {
	c.mu.Lock()
	ch, ok := c.pending[requestID]
	if ok {
		delete(c.pending, requestID)
	}
	c.mu.Unlock()
	if ok {
		ch <- res
		return
	}
	if c.opts.OnFrameError != nil {
		c.opts.OnFrameError(fmt.Errorf("websocket: unmatched requestId %q", requestID), raw)
	}
}

// deliverToSolePending hands res to the only in-flight Do call, if there is
// exactly one. It reports whether the result was delivered.
func (c *Conn) deliverToSolePending(res doResult) bool {
	c.mu.Lock()
	if len(c.pending) != 1 {
		c.mu.Unlock()
		return false
	}
	var ch chan doResult
	for id, pending := range c.pending {
		ch = pending
		delete(c.pending, id)
	}
	c.mu.Unlock()
	ch <- res
	return true
}

// dispatch queues a push frame for dispatchLoop without ever blocking the read
// goroutine. A full buffer means the queued event is stale, so it is dropped.
func (c *Conn) dispatch(ev wsEvent) {
	select {
	case c.events <- ev:
		return
	default:
	}
	select {
	case <-c.events:
	default:
	}
	select {
	case c.events <- ev:
	default:
	}
}

// dispatchLoop runs push handlers off the read goroutine so a handler may call
// Do. It runs for the life of the Conn, across reconnects, and stops on Close.
func (c *Conn) dispatchLoop() {
	for {
		select {
		case <-c.closed:
			return
		case ev := <-c.events:
			switch {
			case ev.state != nil:
				c.mu.Lock()
				h := c.handler
				c.mu.Unlock()
				if h != nil {
					h(ev.state)
				}
			case ev.alert != nil:
				if c.opts.OnCalendarAlert != nil {
					c.opts.OnCalendarAlert(ev.alert)
				}
			}
		}
	}
}
