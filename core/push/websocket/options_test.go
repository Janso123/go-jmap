package websocket

import (
	"context"
	"encoding/json"
	"fmt"
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
			if err := json.Unmarshal(data, &probe); err != nil || probe.Type != "Request" {
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
	for i := 0; i < 2; i++ {
		go func() {
			req := &jmap.Request{}
			req.Invoke(&core.Echo{Hello: "x"})
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
			if err := json.Unmarshal(data, &probe); err != nil {
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
	require.NoError(t, conn.EnablePush([]jmap.EventType{"Email"}, "bbb"))

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
	req.Invoke(&core.Echo{Hello: "after"})
	resp, err := conn.Do(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.GreaterOrEqual(t, dials.Load(), int32(2))
}
