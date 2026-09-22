package websocket

import (
	"context"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core"
	_ "github.com/Janso123/go-jmap/core"
	cws "github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnDoAndPush(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "jmap", r.Header.Get("Sec-WebSocket-Protocol"))
		assert.True(t, strings.HasPrefix(strings.ToLower(r.Header.Get("Authorization")), "basic "))

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
			if err := jsonv2.Unmarshal(data, &probe); err != nil {
				return
			}
			switch probe.Type {
			case "WebSocketPushEnable":
				msg := `{"@type":"StateChange","changed":{"a1":{"Email":"e1"}},"pushState":"bbb"}`
				_ = c.Write(ctx, cws.MessageText, []byte(msg))
			case "WebSocketPushDisable":
				// no reply
			case "Request":
				resp := fmt.Sprintf(
					`{"@type":"Response","requestId":%q,"methodResponses":[["Core/echo",{"Hello":"world"},"0"]],"sessionState":"s"}`,
					probe.ID,
				)
				_ = c.Write(ctx, cws.MessageText, []byte(resp))
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/jmap/ws/"

	jc := (&jmap.Client{}).WithBasicAuth("user", "pass")
	jc.Session = &jmap.Session{
		RawCapabilities: map[jmap.URI]jsontext.Value{
			jmap.CoreURI: []byte(`{}`),
			URI:          []byte(`{"url":"` + wsURL + `","supportsPush":true}`),
		},
		Capabilities: map[jmap.URI]jmap.Capability{
			URI: &WebSocket{URL: wsURL, SupportsPush: true},
		},
	}

	ctx := context.Background()
	conn, err := Dial(ctx, jc)
	require.NoError(t, err)
	defer conn.Close()

	var (
		mu     sync.Mutex
		gotSC  *jmap.StateChange
		scDone = make(chan struct{})
	)
	conn.SetHandler(func(sc *jmap.StateChange) {
		mu.Lock()
		defer mu.Unlock()
		gotSC = sc
		select {
		case <-scDone:
		default:
			close(scDone)
		}
	})

	require.NoError(t, conn.EnablePush(ctx, []jmap.EventType{"Email", "Mailbox"}, "aaa"))

	select {
	case <-scDone:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for StateChange")
	}

	mu.Lock()
	require.NotNil(t, gotSC)
	assert.Equal(t, "bbb", gotSC.PushState)
	assert.Equal(t, "e1", gotSC.Changed["a1"]["Email"])
	mu.Unlock()
	assert.Equal(t, "bbb", conn.PushState())

	req := &jmap.Request{}
	req.Invoke(core.Echo{"hello": "world"})
	resp, err := conn.Do(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Responses, 1)
	assert.Equal(t, "Core/echo", resp.Responses[0].Name)
	assert.Equal(t, "s", resp.SessionState)

	require.NoError(t, conn.DisablePush(ctx))
}

func TestConnDoMarksSessionStale(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "jmap", r.Header.Get("Sec-WebSocket-Protocol"))
		assert.True(t, strings.HasPrefix(strings.ToLower(r.Header.Get("Authorization")), "basic "))

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
			if err := jsonv2.Unmarshal(data, &probe); err != nil {
				return
			}
			switch probe.Type {
			case "WebSocketPushEnable":
				msg := `{"@type":"StateChange","changed":{"a1":{"Email":"e1"}},"pushState":"bbb"}`
				_ = c.Write(ctx, cws.MessageText, []byte(msg))
			case "WebSocketPushDisable":
				// no reply
			case "Request":
				resp := fmt.Sprintf(
					`{"@type":"Response","requestId":%q,"methodResponses":[["Core/echo",{"Hello":"world"},"0"]],"sessionState":"s"}`,
					probe.ID,
				)
				_ = c.Write(ctx, cws.MessageText, []byte(resp))
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/jmap/ws/"

	jc := (&jmap.Client{}).WithBasicAuth("user", "pass")
	jc.Session = &jmap.Session{
		State: "old",
		RawCapabilities: map[jmap.URI]jsontext.Value{
			jmap.CoreURI: []byte(`{}`),
			URI:          []byte(`{"url":"` + wsURL + `","supportsPush":true}`),
		},
		Capabilities: map[jmap.URI]jmap.Capability{
			URI: &WebSocket{URL: wsURL, SupportsPush: true},
		},
	}

	ctx := context.Background()
	conn, err := Dial(ctx, jc)
	require.NoError(t, err)
	defer conn.Close()

	req := &jmap.Request{}
	req.Invoke(core.Echo{"hello": "world"})
	resp, err := conn.Do(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "s", resp.SessionState)
	require.True(t, jc.SessionStale())
}

func TestConnDoRequestError(t *testing.T) {
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
			_ = jsonv2.Unmarshal(data, &probe)
			if probe.Type != "Request" {
				continue
			}
			resp := fmt.Sprintf(
				`{"@type":"RequestError","requestId":%q,"type":"urn:ietf:params:jmap:error:notJSON","status":400,"detail":"bad"}`,
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
	}

	conn, err := DialURL(context.Background(), jc, wsURL)
	require.NoError(t, err)
	defer conn.Close()

	req := &jmap.Request{}
	req.Invoke(core.Echo{"hello": "x"})
	_, err = conn.Do(context.Background(), req)
	require.Error(t, err)
	var re *jmap.RequestError
	require.ErrorAs(t, err, &re)
	assert.Equal(t, 400, re.Status)
}

func TestDialMissingCapability(t *testing.T) {
	jc := &jmap.Client{
		Session: &jmap.Session{
			RawCapabilities: map[jmap.URI]jsontext.Value{},
		},
	}
	_, err := Dial(context.Background(), jc)
	require.Error(t, err)
}

func TestOnFrameErrorCalled(t *testing.T) {
	frameErrCh := make(chan struct {
		err error
		raw []byte
	}, 1)

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
			_ = c.Write(ctx, cws.MessageText, []byte(`{"@type":"Nope","requestId":`+fmt.Sprintf("%q", probe.ID)+`}`))
			resp := fmt.Sprintf(
				`{"@type":"Response","requestId":%q,"methodResponses":[["Core/echo",{"Hello":"ok"},"0"]],"sessionState":"s"}`,
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
	}

	conn, err := DialURL(context.Background(), jc, wsURL, Options{
		OnFrameError: func(err error, raw []byte) {
			select {
			case frameErrCh <- struct {
				err error
				raw []byte
			}{err, append([]byte(nil), raw...)}:
			default:
			}
		},
	})
	require.NoError(t, err)
	defer conn.Close()

	req := &jmap.Request{}
	req.Invoke(core.Echo{"hello": "ok"})
	doCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := conn.Do(doCtx, req)

	select {
	case fe := <-frameErrCh:
		require.Error(t, fe.err)
		assert.Contains(t, string(fe.raw), `"Nope"`)
	case <-time.After(2 * time.Second):
		t.Fatal("OnFrameError was not called")
	}

	// Unknown @type with requestId must unblock Do (error), not hang.
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestPingPublic(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := cws.Accept(w, r, &cws.AcceptOptions{
			Subprotocols:       []string{"jmap"},
			InsecureSkipVerify: true,
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

	conn, err := DialURL(context.Background(), jc, wsURL)
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	require.NoError(t, conn.Ping(ctx))
}

func TestConnDoUnblocksOnNullRequestIdError(t *testing.T) {
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
			_, _, err := c.Read(ctx)
			if err != nil {
				return
			}
			_ = c.Write(ctx, cws.MessageText, []byte(`{"@type":"RequestError","requestId":null,"type":"urn:ietf:params:jmap:error:notJSON","status":400,"detail":"bad"}`))
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := &jmap.Client{
		Session: &jmap.Session{
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: []byte(`{}`),
			},
		},
	}
	conn, err := DialURL(context.Background(), jc, wsURL)
	require.NoError(t, err)
	defer conn.Close()

	req := &jmap.Request{}
	req.Invoke(core.Echo{"hello": "x"})
	doCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = conn.Do(doCtx, req)
	require.Error(t, err)
	require.NotErrorIs(t, err, context.DeadlineExceeded)
	var re *jmap.RequestError
	require.ErrorAs(t, err, &re)
	require.Equal(t, 400, re.Status)
}

func TestConnDoLargeFrame(t *testing.T) {
	big := strings.Repeat("x", 40<<10)
	reply := func(reqID string) string {
		return `{"@type":"Response","requestId":"` + reqID + `","methodResponses":[["Core/echo",{"pad":"` + big + `"},"c0"]],"sessionState":"s"}`
	}

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
				return
			}
			_ = c.Write(ctx, cws.MessageText, []byte(reply(probe.ID)))
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := &jmap.Client{
		Session: &jmap.Session{
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: []byte(`{}`),
			},
		},
	}
	conn, err := DialURL(context.Background(), jc, wsURL)
	require.NoError(t, err)
	defer conn.Close()

	req := &jmap.Request{}
	req.Invoke(core.Echo{"hello": "x"})
	doCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := conn.Do(doCtx, req)
	require.NoError(t, err)
	require.Len(t, resp.Responses, 1)
}

func TestUnmatchedRequestIdDoesNotFailOthers(t *testing.T) {
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
				return
			}
			_ = c.Write(ctx, cws.MessageText, []byte(`{"@type":"Response","requestId":"not-pending","methodResponses":[],"sessionState":"s"}`))
			resp := fmt.Sprintf(
				`{"@type":"Response","requestId":%q,"methodResponses":[["Core/echo",{"hello":"x"},"0"]],"sessionState":"s"}`,
				probe.ID,
			)
			_ = c.Write(ctx, cws.MessageText, []byte(resp))
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := &jmap.Client{
		Session: &jmap.Session{
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: []byte(`{}`),
			},
		},
	}

	var gotFrameErr atomic.Bool
	conn, err := DialURL(context.Background(), jc, wsURL, Options{
		OnFrameError: func(error, []byte) { gotFrameErr.Store(true) },
	})
	require.NoError(t, err)
	defer conn.Close()

	req := &jmap.Request{}
	req.Invoke(core.Echo{"hello": "x"})
	doCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := conn.Do(doCtx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, gotFrameErr.Load())
}

func TestStateChangeHandlerMayCallDo(t *testing.T) {
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
			if err := jsonv2.Unmarshal(data, &probe); err != nil {
				return
			}
			switch probe.Type {
			case "WebSocketPushEnable":
				_ = c.Write(ctx, cws.MessageText, []byte(`{"@type":"StateChange","changed":{"a1":{"Email":"e1"}},"pushState":"p1"}`))
			case "Request":
				resp := fmt.Sprintf(
					`{"@type":"Response","requestId":%q,"methodResponses":[["Core/echo",{"hello":"x"},"0"]],"sessionState":"s"}`,
					probe.ID,
				)
				_ = c.Write(ctx, cws.MessageText, []byte(resp))
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := &jmap.Client{
		Session: &jmap.Session{
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: []byte(`{}`),
			},
		},
	}
	conn, err := DialURL(context.Background(), jc, wsURL)
	require.NoError(t, err)
	defer conn.Close()

	echoRequest := func() *jmap.Request {
		req := &jmap.Request{}
		req.Invoke(core.Echo{"hello": "x"})
		return req
	}

	done := make(chan struct{})
	var handlerErr error
	conn.SetHandler(func(sc *jmap.StateChange) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, handlerErr = conn.Do(ctx, echoRequest())
		close(done)
	})
	require.NoError(t, conn.EnablePush(context.Background(), nil, ""))

	select {
	case <-done:
		require.NoError(t, handlerErr)
	case <-time.After(3 * time.Second):
		t.Fatal("handler deadlocked")
	}
}

func TestDialRejectsForeignCapabilityOrigin(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := cws.Accept(w, r, &cws.AcceptOptions{
			Subprotocols:       []string{"jmap"},
			InsecureSkipVerify: true,
		})
		if err != nil {
			return
		}
		defer c.CloseNow()
		_, _, _ = c.Read(context.Background())
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/jmap/ws/"
	jc := (&jmap.Client{}).WithAccessToken("tok")
	jc.SessionEndpoint = "https://mail.example.com/jmap/session"
	jc.Session = &jmap.Session{
		RawCapabilities: map[jmap.URI]jsontext.Value{
			jmap.CoreURI: []byte(`{}`),
			URI:          []byte(`{"url":"` + wsURL + `","supportsPush":true}`),
		},
		Capabilities: map[jmap.URI]jmap.Capability{
			URI: &WebSocket{URL: wsURL, SupportsPush: true},
		},
	}

	_, err := Dial(context.Background(), jc)
	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "origin")

	conn, err := Dial(context.Background(), jc, Options{AllowForeignOrigin: true})
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func TestDialURLRejectsNonLoopbackCleartext(t *testing.T) {
	jc := &jmap.Client{
		Session: &jmap.Session{
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: []byte(`{}`),
			},
		},
	}
	_, err := DialURL(context.Background(), jc, "ws://example.invalid/jmap/ws/")
	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "wss")
}

func TestEnablePushUnsupported(t *testing.T) {
	sawEnable := make(chan struct{}, 1)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := cws.Accept(w, r, &cws.AcceptOptions{
			Subprotocols:       []string{"jmap"},
			InsecureSkipVerify: true,
		})
		if err != nil {
			return
		}
		defer c.CloseNow()
		_, data, err := c.Read(context.Background())
		if err != nil {
			return
		}
		var probe struct {
			Type string `json:"@type"`
		}
		if jsonv2.Unmarshal(data, &probe) == nil && probe.Type == "WebSocketPushEnable" {
			sawEnable <- struct{}{}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := &jmap.Client{
		Session: &jmap.Session{
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: []byte(`{}`),
			},
			Capabilities: map[jmap.URI]jmap.Capability{
				URI: &WebSocket{URL: wsURL, SupportsPush: false},
			},
		},
	}
	conn, err := DialURL(context.Background(), jc, wsURL)
	require.NoError(t, err)
	defer conn.Close()

	err = conn.EnablePush(context.Background(), nil, "")
	require.ErrorIs(t, err, ErrPushUnsupported)
	select {
	case <-sawEnable:
		t.Fatal("EnablePush wrote a frame when SupportsPush is false")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestDoWaitsForReconnect(t *testing.T) {
	var refuse atomic.Bool
	down := make(chan struct{})
	var downOnce sync.Once

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if refuse.Load() {
			http.Error(w, "down", http.StatusServiceUnavailable)
			return
		}
		c, err := cws.Accept(w, r, &cws.AcceptOptions{
			Subprotocols:       []string{"jmap"},
			InsecureSkipVerify: true,
		})
		if err != nil {
			return
		}
		defer c.CloseNow()
		ctx := context.Background()
		_, data, err := c.Read(ctx)
		if err != nil {
			return
		}
		var probe struct {
			Type string `json:"@type"`
			ID   string `json:"id"`
		}
		if jsonv2.Unmarshal(data, &probe) != nil || probe.Type != "Request" {
			return
		}
		resp := fmt.Sprintf(
			`{"@type":"Response","requestId":%q,"methodResponses":[["Core/echo",{"hello":"back"},"0"]],"sessionState":"s"}`,
			probe.ID,
		)
		_ = c.Write(ctx, cws.MessageText, []byte(resp))
		_, _, _ = c.Read(ctx)
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http") + "/"
	jc := &jmap.Client{
		Session: &jmap.Session{
			RawCapabilities: map[jmap.URI]jsontext.Value{
				jmap.CoreURI: []byte(`{}`),
			},
		},
	}

	conn, err := DialURL(context.Background(), jc, wsURL, Options{
		Reconnect: &ReconnectOptions{
			MinBackoff: 10 * time.Millisecond,
			MaxBackoff: 30 * time.Millisecond,
			OnDisconnect: func(error) {
				refuse.Store(true)
				downOnce.Do(func() { close(down) })
			},
		},
	})
	require.NoError(t, err)
	defer conn.Close()

	// Drop the socket so the next dials fail until the test brings the server back.
	conn.mu.Lock()
	ws := conn.ws
	conn.mu.Unlock()
	require.NotNil(t, ws)
	_ = ws.CloseNow()

	select {
	case <-down:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for disconnect")
	}

	req := &jmap.Request{}
	req.Invoke(core.Echo{"hello": "back"})
	doCtx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	type doOut struct {
		resp *jmap.Response
		err  error
	}
	done := make(chan doOut, 1)
	go func() {
		resp, err := conn.Do(doCtx, req)
		done <- doOut{resp, err}
	}()

	select {
	case out := <-done:
		t.Fatalf("Do returned while the server was down: resp=%v err=%v", out.resp, out.err)
	case <-time.After(200 * time.Millisecond):
	}

	refuse.Store(false)

	select {
	case out := <-done:
		require.NoError(t, out.err)
		require.NotNil(t, out.resp)
		require.Len(t, out.resp.Responses, 1)
		assert.Equal(t, "Core/echo", out.resp.Responses[0].Name)
	case <-time.After(3 * time.Second):
		t.Fatal("Do did not return after the server came back")
	}
}
