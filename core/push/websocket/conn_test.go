package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/core"
	_ "git.sr.ht/~rockorager/go-jmap/core"
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
			if err := json.Unmarshal(data, &probe); err != nil {
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
		RawCapabilities: map[jmap.URI]json.RawMessage{
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

	require.NoError(t, conn.EnablePush([]jmap.EventType{"Email", "Mailbox"}, "aaa"))

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
	req.Invoke(&core.Echo{Hello: "world"})
	resp, err := conn.Do(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Responses, 1)
	assert.Equal(t, "Core/echo", resp.Responses[0].Name)
	assert.Equal(t, "s", resp.SessionState)

	require.NoError(t, conn.DisablePush())
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
			_ = json.Unmarshal(data, &probe)
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
		RawCapabilities: map[jmap.URI]json.RawMessage{
			jmap.CoreURI: []byte(`{}`),
		},
	}

	conn, err := DialURL(context.Background(), jc, wsURL)
	require.NoError(t, err)
	defer conn.Close()

	req := &jmap.Request{}
	req.Invoke(&core.Echo{Hello: "x"})
	_, err = conn.Do(context.Background(), req)
	require.Error(t, err)
	var re *jmap.RequestError
	require.ErrorAs(t, err, &re)
	assert.Equal(t, 400, re.Status)
}

func TestDialMissingCapability(t *testing.T) {
	jc := &jmap.Client{
		Session: &jmap.Session{
			RawCapabilities: map[jmap.URI]json.RawMessage{},
		},
	}
	_, err := Dial(context.Background(), jc)
	require.Error(t, err)
}
