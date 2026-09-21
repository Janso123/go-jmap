package push

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventSourcePingAndLastEventID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "text/event-stream", r.Header.Get("Accept"))
		require.Equal(t, "42", r.Header.Get("Last-Event-ID"))
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, ": comment\n\nevent: ping\ndata: {}\n\nid: 43\nevent: state\ndata: {\"changed\":{}}\n\n")
	}))
	t.Cleanup(srv.Close)

	var pinged atomic.Bool
	var got *jmap.StateChange
	es := &EventSource{
		Client: &jmap.Client{
			HttpClient: srv.Client(),
			Session:    &jmap.Session{EventSourceURL: srv.URL},
		},
		LastEventID: "42",
		OnPing:      func() { pinged.Store(true) },
		Handler: func(sc *jmap.StateChange) {
			got = sc
		},
	}

	err := es.Listen(context.Background())
	require.ErrorIs(t, err, io.EOF)
	require.True(t, pinged.Load(), "OnPing should be called")
	require.NotNil(t, got)
	assert.Equal(t, "43", es.LastEventID)
}

func TestEventSourceMultilineData(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "event: state\ndata: {\"changed\":{\"a\":\ndata: {\"Email\":\"s1\"}}}\n\n")
	}))
	t.Cleanup(srv.Close)

	var got *jmap.StateChange
	es := &EventSource{
		Client: &jmap.Client{
			HttpClient: srv.Client(),
			Session:    &jmap.Session{EventSourceURL: srv.URL},
		},
		Handler: func(sc *jmap.StateChange) { got = sc },
	}

	err := es.Listen(context.Background())
	require.ErrorIs(t, err, io.EOF)
	require.NotNil(t, got)
	require.Equal(t, "s1", got.Changed["a"]["Email"])
}

func TestEventSourceRetryAndLargeEvent(t *testing.T) {
	payload := `{"changed":{"a":{"Email":"` + strings.Repeat("x", 70*1024) + `"}}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintf(w, "retry: 1500\nevent: state\ndata: %s\n\n", payload)
	}))
	t.Cleanup(srv.Close)

	var got *jmap.StateChange
	es := &EventSource{
		Client: &jmap.Client{
			HttpClient: srv.Client(),
			Session:    &jmap.Session{EventSourceURL: srv.URL},
		},
		MaxEventSize: 1 << 20,
		Handler:      func(sc *jmap.StateChange) { got = sc },
	}

	err := es.Listen(context.Background())
	require.ErrorIs(t, err, io.EOF)
	require.NotNil(t, got)
	assert.Equal(t, 1500*time.Millisecond, es.Retry)
	require.Contains(t, got.Changed["a"]["Email"], "xxx")
}

func TestEventSourceCalendarAlert(t *testing.T) {
	raw := `{"@type":"CalendarAlert","accountId":"a1","calendarEventId":"e1","uid":"u1","alertId":"al1"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintf(w, "event: calendarAlert\ndata: %s\n\n", raw)
	}))
	t.Cleanup(srv.Close)

	var alert *calendar.CalendarAlert
	es := &EventSource{
		Client: &jmap.Client{
			HttpClient: srv.Client(),
			Session:    &jmap.Session{EventSourceURL: srv.URL},
		},
		OnCalendarAlert: func(a *calendar.CalendarAlert) { alert = a },
	}

	err := es.Listen(context.Background())
	require.ErrorIs(t, err, io.EOF)
	require.NotNil(t, alert)
	assert.Equal(t, "CalendarAlert", alert.Type)
	assert.Equal(t, jmap.ID("a1"), alert.AccountID)
	assert.Equal(t, jmap.ID("e1"), alert.CalendarEventID)
	assert.Equal(t, "u1", alert.UID)
	assert.Equal(t, "al1", alert.AlertID)
}

func TestEventSourceContextCancel(t *testing.T) {
	started := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.(http.Flusher).Flush()
		close(started)
		<-r.Context().Done()
	}))
	t.Cleanup(srv.Close)

	es := &EventSource{
		Client: &jmap.Client{
			HttpClient: srv.Client(),
			Session:    &jmap.Session{EventSourceURL: srv.URL},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- es.Listen(ctx) }()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("server did not start stream")
	}
	cancel()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("Listen did not return after cancel")
	}
}

func TestEventSourceCloseReturnsErrClosed(t *testing.T) {
	started := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.(http.Flusher).Flush()
		close(started)
		<-r.Context().Done()
	}))
	t.Cleanup(srv.Close)

	es := &EventSource{
		Client: &jmap.Client{
			HttpClient: srv.Client(),
			Session:    &jmap.Session{EventSourceURL: srv.URL},
		},
	}

	errCh := make(chan error, 1)
	go func() { errCh <- es.Listen(context.Background()) }()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("server did not start stream")
	}
	es.Close()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, ErrClosed)
	case <-time.After(2 * time.Second):
		t.Fatal("Listen did not return after Close")
	}
}

func TestEventSourceNilHTTPClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "text/event-stream", r.Header.Get("Accept"))
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "event: ping\ndata: {}\n\n")
	}))
	t.Cleanup(srv.Close)

	var pinged atomic.Bool
	es := &EventSource{
		Client: &jmap.Client{
			HttpClient: nil, // must not panic
			Session:    &jmap.Session{EventSourceURL: srv.URL},
		},
		OnPing: func() { pinged.Store(true) },
	}

	err := es.Listen(context.Background())
	require.ErrorIs(t, err, io.EOF)
	require.True(t, pinged.Load())
}

func TestEventSourceReconnectPreservesLastEventID(t *testing.T) {
	var mu sync.Mutex
	var lastIDs []string
	var n atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		lastIDs = append(lastIDs, r.Header.Get("Last-Event-ID"))
		mu.Unlock()

		w.Header().Set("Content-Type", "text/event-stream")
		switch n.Add(1) {
		case 1:
			fmt.Fprint(w, "id: 10\nevent: state\ndata: {\"changed\":{}}\n\n")
		default:
			fmt.Fprint(w, "event: ping\ndata: {}\n\n")
		}
	}))
	t.Cleanup(srv.Close)

	var pinged atomic.Bool
	es := &EventSource{
		Client: &jmap.Client{
			HttpClient: srv.Client(),
			Session:    &jmap.Session{EventSourceURL: srv.URL},
		},
		LastEventID: "1",
		Reconnect: &ReconnectOptions{
			MinBackoff: 5 * time.Millisecond,
			MaxBackoff: 20 * time.Millisecond,
		},
		Handler: func(*jmap.StateChange) {},
	}
	es.OnPing = func() {
		pinged.Store(true)
		es.Close()
	}

	err := es.Listen(context.Background())
	require.ErrorIs(t, err, ErrClosed)
	require.True(t, pinged.Load())

	mu.Lock()
	defer mu.Unlock()
	require.GreaterOrEqual(t, len(lastIDs), 2)
	assert.Equal(t, "1", lastIDs[0])
	assert.Equal(t, "10", lastIDs[1])
}
