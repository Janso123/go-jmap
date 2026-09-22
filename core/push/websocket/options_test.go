package websocket

import (
	"bytes"
	"context"
	jsonv2 "encoding/json/v2"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"encoding/json/jsontext"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core"
	_ "github.com/Janso123/go-jmap/core"
	cws "github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDialURLCompressionExtension(t *testing.T) {
	sawExt := make(chan string, 1)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawExt <- r.Header.Get("Sec-WebSocket-Extensions")
		c, err := cws.Accept(w, r, &cws.AcceptOptions{
			Subprotocols:       []string{"jmap"},
			InsecureSkipVerify: true,
			CompressionMode:    cws.CompressionNoContextTakeover,
		})
		if err != nil {
			return
		}
		defer c.CloseNow()
		_, _, _ = c.Read(context.Background())
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := (&jmap.Client{}).WithAccessToken("tok")
	jc.Session = &jmap.Session{
		RawCapabilities: map[jmap.URI]jsontext.Value{
			jmap.CoreURI: []byte(`{}`),
		},
	}

	conn, err := DialURL(context.Background(), jc, wsURL, Options{
		CompressionMode: CompressionNoContextTakeover,
	})
	require.NoError(t, err)
	defer conn.Close()

	ext := <-sawExt
	assert.Contains(t, strings.ToLower(ext), "permessage-deflate")
}

func TestDialURLMaxConcurrentRequests(t *testing.T) {
	var mu sync.Mutex
	inflight := 0
	maxSeen := 0
	release := make(chan struct{})

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := cws.Accept(w, r, &cws.AcceptOptions{
			Subprotocols:       []string{"jmap"},
			InsecureSkipVerify: true,
		})
		if err != nil {
			return
		}
		defer c.CloseNow()
		ctx := context.Background()
		for {
			_, data, err := c.Read(ctx)
			if err != nil {
				return
			}
			var probe struct {
				Type string `json:"@type"`
				ID   string `json:"id"`
			}
			if err := jsonv2.Unmarshal(data, &probe); err != nil || probe.Type != "Request" {
				continue
			}
			mu.Lock()
			inflight++
			if inflight > maxSeen {
				maxSeen = inflight
			}
			mu.Unlock()

			<-release

			mu.Lock()
			inflight--
			mu.Unlock()

			resp := fmt.Sprintf(
				`{"@type":"Response","requestId":%q,"methodResponses":[["Core/echo",{},"0"]],"sessionState":"s"}`,
				probe.ID,
			)
			_ = c.Write(ctx, cws.MessageText, []byte(resp))
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := (&jmap.Client{}).WithAccessToken("tok")
	jc.Session = &jmap.Session{
		RawCapabilities: map[jmap.URI]jsontext.Value{
			jmap.CoreURI: []byte(`{}`),
		},
		Capabilities: map[jmap.URI]jmap.Capability{
			jmap.CoreURI: &core.Core{MaxConcurrentRequests: 1},
		},
	}

	conn, err := DialURL(context.Background(), jc, wsURL, Options{MaxConcurrentRequests: 1})
	require.NoError(t, err)
	defer conn.Close()

	ctx := context.Background()
	done := make(chan error, 2)
	for range 2 {
		go func() {
			req := &jmap.Request{}
			req.Invoke(core.Echo{"hello": "x"})
			_, err := conn.Do(ctx, req)
			done <- err
		}()
	}

	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	assert.Equal(t, 1, maxSeen, "client must not exceed maxConcurrentRequests")
	mu.Unlock()

	close(release)
	require.NoError(t, <-done)
	require.NoError(t, <-done)
	mu.Lock()
	assert.Equal(t, 1, maxSeen)
	mu.Unlock()
}

func TestPingIntervalKeepalive(t *testing.T) {
	sawPing := make(chan struct{}, 1)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := cws.Accept(w, r, &cws.AcceptOptions{
			Subprotocols:       []string{"jmap"},
			InsecureSkipVerify: true,
			OnPingReceived: func(ctx context.Context, payload []byte) bool {
				select {
				case sawPing <- struct{}{}:
				default:
				}
				return true // write pong
			},
		})
		if err != nil {
			return
		}
		defer c.CloseNow()
		_, _, _ = c.Read(context.Background())
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := (&jmap.Client{}).WithAccessToken("tok")
	jc.Session = &jmap.Session{
		RawCapabilities: map[jmap.URI]jsontext.Value{
			jmap.CoreURI: []byte(`{}`),
		},
	}

	conn, err := DialURL(context.Background(), jc, wsURL, Options{
		PingInterval: 30 * time.Millisecond,
	})
	require.NoError(t, err)
	defer conn.Close()

	select {
	case <-sawPing:
	case <-time.After(2 * time.Second):
		t.Fatal("PingInterval keepalive did not send a ping")
	}
}

func TestConnAutoReconnect(t *testing.T) {
	var dials atomic.Int32
	pushStates := make(chan string, 4)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := dials.Add(1)
		c, err := cws.Accept(w, r, &cws.AcceptOptions{
			Subprotocols:       []string{"jmap"},
			InsecureSkipVerify: true,
		})
		if err != nil {
			return
		}
		defer c.CloseNow()
		ctx := context.Background()
		for {
			_, data, err := c.Read(ctx)
			if err != nil {
				return
			}
			var probe struct {
				Type      string `json:"@type"`
				ID        string `json:"id"`
				PushState string `json:"pushState"`
			}
			if err := jsonv2.Unmarshal(data, &probe); err != nil {
				continue
			}
			switch probe.Type {
			case "WebSocketPushEnable":
				pushStates <- probe.PushState
				if n == 1 {
					_ = c.Close(cws.StatusGoingAway, "bye")
					return
				}
			case "Request":
				resp := fmt.Sprintf(
					`{"@type":"Response","requestId":%q,"methodResponses":[["Core/echo",{},"0"]],"sessionState":"s"}`,
					probe.ID,
				)
				_ = c.Write(ctx, cws.MessageText, []byte(resp))
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := (&jmap.Client{}).WithAccessToken("tok")
	jc.Session = &jmap.Session{
		RawCapabilities: map[jmap.URI]jsontext.Value{
			jmap.CoreURI: []byte(`{}`),
		},
	}

	reconnected := make(chan struct{}, 1)
	conn, err := DialURL(context.Background(), jc, wsURL, Options{
		Reconnect: &ReconnectOptions{
			MinBackoff: 10 * time.Millisecond,
			MaxBackoff: 50 * time.Millisecond,
			OnReconnect: func() {
				select {
				case reconnected <- struct{}{}:
				default:
				}
			},
		},
	})
	require.NoError(t, err)
	defer conn.Close()

	// Seed pushState as if prior StateChange arrived, then enable push.
	conn.mu.Lock()
	conn.pushState = "bbb"
	conn.mu.Unlock()
	require.NoError(t, conn.EnablePush(context.Background(), []jmap.EventType{"Email"}, "bbb"))

	select {
	case ps := <-pushStates:
		assert.Equal(t, "bbb", ps)
	case <-time.After(2 * time.Second):
		t.Fatal("first EnablePush not received")
	}

	select {
	case <-reconnected:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for reconnect")
	}

	select {
	case ps := <-pushStates:
		assert.Equal(t, "bbb", ps, "reconnect must re-EnablePush with last pushState")
	case <-time.After(2 * time.Second):
		t.Fatal("re-EnablePush not received after reconnect")
	}

	req := &jmap.Request{}
	req.Invoke(core.Echo{"hello": "after"})
	resp, err := conn.Do(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.GreaterOrEqual(t, dials.Load(), int32(2))
}

func TestReconnectEnablePushFailure(t *testing.T) {
	gotPush := make(chan string, 4)
	var accepts atomic.Int32
	releaseFirst := make(chan struct{})
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(releaseFirst) }) }
	defer release()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := cws.Accept(w, r, &cws.AcceptOptions{
			Subprotocols:       []string{"jmap"},
			InsecureSkipVerify: true,
		})
		if err != nil {
			return
		}
		defer c.CloseNow()
		n := accepts.Add(1)
		ctx := context.Background()
		_, data, err := c.Read(ctx)
		if err != nil {
			return
		}
		var probe struct {
			Type      string `json:"@type"`
			PushState string `json:"pushState"`
		}
		if jsonv2.Unmarshal(data, &probe) == nil && probe.Type == "WebSocketPushEnable" {
			gotPush <- probe.PushState
		}
		if n == 1 {
			<-releaseFirst
			return
		}
		_, _, _ = c.Read(ctx)
	}))
	defer s.Close()

	var dialN atomic.Int32
	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			var d net.Dialer
			c, err := d.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			n := dialN.Add(1)
			return &failPostHandshakeWriteConn{Conn: c, fail: n == 2}, nil
		},
	}

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := (&jmap.Client{HttpClient: &http.Client{Transport: tr}}).WithAccessToken("tok")
	jc.Session = &jmap.Session{
		RawCapabilities: map[jmap.URI]jsontext.Value{
			jmap.CoreURI: []byte(`{}`),
		},
		Capabilities: map[jmap.URI]jmap.Capability{
			URI: &WebSocket{URL: wsURL, SupportsPush: true},
		},
	}

	var conn *Conn
	var err error
	second := make(chan bool, 1)
	reconnected := make(chan struct{}, 1)
	var discs atomic.Int32
	conn, err = DialURL(context.Background(), jc, wsURL, Options{
		Reconnect: &ReconnectOptions{
			MinBackoff: 10 * time.Millisecond,
			MaxBackoff: 40 * time.Millisecond,
			OnReconnect: func() {
				// Fires only after re-EnablePush returns and pushEnabled is set.
				select {
				case reconnected <- struct{}{}:
				default:
				}
			},
			OnDisconnect: func(error) {
				if discs.Add(1) == 2 && conn != nil {
					conn.mu.Lock()
					enabled := conn.pushEnabled
					conn.mu.Unlock()
					second <- enabled
				}
			},
		},
	})
	require.NoError(t, err)
	defer conn.Close()

	require.NoError(t, conn.EnablePush(context.Background(), []jmap.EventType{"Email"}, "bbb"))
	select {
	case ps := <-gotPush:
		assert.Equal(t, "bbb", ps)
	case <-time.After(2 * time.Second):
		t.Fatal("first EnablePush not received")
	}
	conn.mu.Lock()
	assert.True(t, conn.pushEnabled)
	conn.mu.Unlock()
	release()

	select {
	case enabled := <-second:
		assert.False(t, enabled, "pushEnabled must stay false when re-EnablePush write fails")
	case <-time.After(3 * time.Second):
		t.Fatal("re-EnablePush failure did not trigger OnDisconnect again")
	}

	select {
	case <-reconnected:
	case <-time.After(3 * time.Second):
		t.Fatal("push enable did not succeed on a later reconnect")
	}
	select {
	case ps := <-gotPush:
		assert.Equal(t, "bbb", ps)
	case <-time.After(3 * time.Second):
		t.Fatal("push enable did not succeed on a later reconnect")
	}
	conn.mu.Lock()
	assert.True(t, conn.pushEnabled, "pushEnabled must be true once re-EnablePush has completed")
	conn.mu.Unlock()
}

// failPostHandshakeWriteConn fails Write after the HTTP response headers
// have been read, so the WebSocket handshake can succeed and the next
// client frame (re-EnablePush) fails.
type failPostHandshakeWriteConn struct {
	net.Conn
	fail bool

	mu  sync.Mutex
	buf []byte
	hdr bool
}

func (c *failPostHandshakeWriteConn) Read(p []byte) (int, error) {
	n, err := c.Conn.Read(p)
	if n > 0 {
		c.mu.Lock()
		if !c.hdr {
			c.buf = append(c.buf, p[:n]...)
			if bytes.Contains(c.buf, []byte("\r\n\r\n")) {
				c.hdr = true
				c.buf = nil
			}
		}
		c.mu.Unlock()
	}
	return n, err
}

func (c *failPostHandshakeWriteConn) Write(p []byte) (int, error) {
	c.mu.Lock()
	fail := c.fail && c.hdr
	c.mu.Unlock()
	if fail {
		return 0, errors.New("forced websocket write failure")
	}
	return c.Conn.Write(p)
}
