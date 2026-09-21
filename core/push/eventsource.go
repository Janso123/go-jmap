package push

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
)

// ErrClosed is returned by Listen when Close is called.
var ErrClosed = errors.New("eventsource: closed")

const defaultMaxEventSize = 1 << 20 // 1 MiB

// ReconnectOptions controls optional EventSource auto-reconnect.
type ReconnectOptions struct {
	// MinBackoff is the initial wait after disconnect. Default: 1s.
	MinBackoff time.Duration
	// MaxBackoff caps exponential backoff. Default: 30s.
	MaxBackoff time.Duration
	// OnDisconnect is called when the stream ends unexpectedly (before backoff).
	OnDisconnect func(error)
	// OnReconnect is called after a successful re-connect.
	OnReconnect func()
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

// A subscription to an event stream
type EventSource struct {
	// The JMAP client to use for the stream
	Client *jmap.Client

	// The function to pass state change events to
	Handler func(*jmap.StateChange)

	// OnPing is called when the server sends an event: ping frame.
	OnPing func()

	// OnCalendarAlert is called for SSE event: calendarAlert with a decoded
	// CalendarAlert object (draft-ietf-jmap-calendars §6.4).
	OnCalendarAlert func(*calendar.CalendarAlert)

	// The events to subscribe to. If left unset, will default to AllEvents
	Events []jmap.EventType

	// Interval the server should ping the client at, in seconds. The server
	// may choose to ignore this value. Set to 0 to disable pinging (which
	// the server may also ignore)
	Ping uint

	// Whether to close the connection after a state event
	CloseAfterState bool

	// LastEventID is sent as Last-Event-ID on connect and updated from id: fields.
	LastEventID string

	// MaxEventSize is the maximum SSE line/event buffer size. Default: 1 MiB.
	MaxEventSize int

	// Retry is updated from retry: fields (milliseconds → duration).
	Retry time.Duration

	// Reconnect, if set, redials after clean EOF while preserving LastEventID.
	Reconnect *ReconnectOptions

	mu     sync.Mutex
	resp   *http.Response
	closed bool
}

func (e *EventSource) httpClient() *http.Client {
	hc := e.Client.HttpClient
	if hc == nil {
		return http.DefaultClient
	}
	return hc
}

func (e *EventSource) maxEventSize() int {
	if e.MaxEventSize > 0 {
		return e.MaxEventSize
	}
	return defaultMaxEventSize
}

// Connect to the server
func (e *EventSource) connect(ctx context.Context) error {
	// Create the URL for the subscription
	u, err := url.Parse(e.Client.Session.EventSourceURL)
	if err != nil {
		return err
	}
	q := u.Query()

	if len(e.Events) == 0 {
		e.Events = []jmap.EventType{jmap.AllEvents}
	}
	// types field
	types := []string{}
	for _, ev := range e.Events {
		types = append(types, string(ev))
	}
	typeStr := strings.Join(types, ",")
	q.Set("types", typeStr)

	// ping field
	q.Set("ping", fmt.Sprintf("%d", e.Ping))

	// close after field
	closeAfter := "no"
	if e.CloseAfterState {
		closeAfter = "state"
	}
	q.Set("closeafter", closeAfter)

	// set the query string
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")
	if e.LastEventID != "" {
		req.Header.Set("Last-Event-ID", e.LastEventID)
	}

	resp, err := e.httpClient().Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return fmt.Errorf("invalid request, response code: %d", resp.StatusCode)
	}

	e.mu.Lock()
	if e.closed {
		e.mu.Unlock()
		resp.Body.Close()
		return ErrClosed
	}
	e.resp = resp
	e.mu.Unlock()
	return nil
}

// Listen starts listening for events. It blocks until the stream ends, the
// context is canceled, or Close is called.
//
// On a clean server close with Reconnect unset, Listen returns io.EOF.
// Close yields ErrClosed. Context cancel yields ctx.Err().
func (e *EventSource) Listen(ctx context.Context) error {
	e.mu.Lock()
	e.closed = false
	e.mu.Unlock()

	backoff := time.Duration(0)
	if e.Reconnect != nil {
		backoff = e.Reconnect.minBackoff()
	}
	maxBackoff := time.Duration(0)
	if e.Reconnect != nil {
		maxBackoff = e.Reconnect.maxBackoff()
	}

	attempt := 0
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		e.mu.Lock()
		closed := e.closed
		e.mu.Unlock()
		if closed {
			return ErrClosed
		}

		err := e.connect(ctx)
		if err != nil {
			if errors.Is(err, ErrClosed) {
				return ErrClosed
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if e.Reconnect == nil {
				return err
			}
			if e.Reconnect.OnDisconnect != nil {
				e.Reconnect.OnDisconnect(err)
			}
			if err := e.waitBackoff(ctx, backoff); err != nil {
				return err
			}
			backoff = nextBackoff(backoff, maxBackoff)
			continue
		}

		if attempt > 0 && e.Reconnect != nil && e.Reconnect.OnReconnect != nil {
			e.Reconnect.OnReconnect()
		}
		attempt++
		if e.Reconnect != nil {
			backoff = e.Reconnect.minBackoff()
		}

		err = e.readStream(ctx)
		e.closeResp()

		if errors.Is(err, ErrClosed) {
			return ErrClosed
		}
		if err != nil && !errors.Is(err, io.EOF) {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if e.Reconnect == nil {
				return err
			}
			if e.Reconnect.OnDisconnect != nil {
				e.Reconnect.OnDisconnect(err)
			}
			if waitErr := e.waitBackoff(ctx, backoff); waitErr != nil {
				return waitErr
			}
			backoff = nextBackoff(backoff, maxBackoff)
			continue
		}

		// Clean EOF
		if e.Reconnect == nil {
			if err == nil {
				return io.EOF
			}
			return err
		}
		if e.Reconnect.OnDisconnect != nil {
			e.Reconnect.OnDisconnect(io.EOF)
		}
		if waitErr := e.waitBackoff(ctx, backoff); waitErr != nil {
			return waitErr
		}
		backoff = nextBackoff(backoff, maxBackoff)
	}
}

func nextBackoff(cur, max time.Duration) time.Duration {
	if cur <= 0 {
		return time.Second
	}
	next := cur * 2
	if max > 0 && next > max {
		return max
	}
	return next
}

func (e *EventSource) waitBackoff(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		d = time.Second
	}
	// Prefer Retry from server if set and larger.
	if e.Retry > d {
		d = e.Retry
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		e.mu.Lock()
		closed := e.closed
		e.mu.Unlock()
		if closed {
			return ErrClosed
		}
		return nil
	}
}

func (e *EventSource) readStream(ctx context.Context) error {
	e.mu.Lock()
	resp := e.resp
	e.mu.Unlock()
	if resp == nil || resp.Body == nil {
		return io.EOF
	}

	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			e.closeResp()
		case <-done:
		}
	}()

	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, e.maxEventSize())

	var eventType string
	var dataBuf bytes.Buffer
	var id string
	hasID := false

	dispatch := func() error {
		defer func() {
			eventType = ""
			dataBuf.Reset()
			id = ""
			hasID = false
		}()

		if hasID {
			e.LastEventID = id
		}

		data := dataBuf.Bytes()
		// No data and no event type means empty comment-only block — skip.
		if len(data) == 0 && eventType == "" {
			return nil
		}

		switch eventType {
		case "", "message":
			// ignore unnamed / generic
		case "ping":
			if e.OnPing != nil {
				e.OnPing()
			}
		case "state":
			if e.Handler != nil && len(data) > 0 {
				state := &jmap.StateChange{}
				if err := json.Unmarshal(data, state); err != nil {
					return err
				}
				e.Handler(state)
			}
		case "calendarAlert":
			if e.OnCalendarAlert != nil && len(data) > 0 {
				alert := &calendar.CalendarAlert{}
				if err := json.Unmarshal(data, alert); err != nil {
					return err
				}
				e.OnCalendarAlert(alert)
			}
		}
		return nil
	}

	for scanner.Scan() {
		e.mu.Lock()
		closed := e.closed
		e.mu.Unlock()
		if closed {
			return ErrClosed
		}
		if err := ctx.Err(); err != nil {
			return err
		}

		line := scanner.Text()
		switch {
		case line == "":
			if err := dispatch(); err != nil {
				return err
			}
		case strings.HasPrefix(line, ":"):
			// comment — ignore
		case strings.Contains(line, ":"):
			field, value, _ := strings.Cut(line, ":")
			if strings.HasPrefix(value, " ") {
				value = value[1:]
			}
			switch field {
			case "event":
				eventType = value
			case "data":
				if dataBuf.Len() > 0 {
					dataBuf.WriteByte('\n')
				}
				dataBuf.WriteString(value)
			case "id":
				if !strings.ContainsRune(value, '\x00') {
					id = value
					hasID = true
				}
			case "retry":
				if ms, err := strconv.Atoi(value); err == nil && ms >= 0 {
					e.Retry = time.Duration(ms) * time.Millisecond
				}
			}
		default:
			// line with no colon: field name with empty value (WHATWG)
			switch line {
			case "event":
				eventType = ""
			case "data":
				if dataBuf.Len() > 0 {
					dataBuf.WriteByte('\n')
				}
			case "id":
				id = ""
				hasID = true
			}
		}
	}

	e.mu.Lock()
	closed := e.closed
	e.mu.Unlock()
	if closed {
		return ErrClosed
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return io.EOF
}

func (e *EventSource) closeResp() {
	e.mu.Lock()
	resp := e.resp
	e.resp = nil
	e.mu.Unlock()
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

// Close closes the stream. Concurrent Listen returns ErrClosed.
func (e *EventSource) Close() {
	e.mu.Lock()
	e.closed = true
	resp := e.resp
	e.resp = nil
	e.mu.Unlock()
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}
