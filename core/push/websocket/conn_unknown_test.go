package websocket

import (
	"context"
	"encoding/json"
	"encoding/json/jsontext"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core"
	_ "github.com/Janso123/go-jmap/core"
	cws "github.com/coder/websocket"
	"github.com/stretchr/testify/require"
)

func TestConnDoUnknownMethodInMethodResponses(t *testing.T) {
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
			if err := json.Unmarshal(data, &probe); err != nil {
				return
			}
			if probe.Type != "Request" {
				continue
			}
			resp := fmt.Sprintf(
				`{"@type":"Response","requestId":%q,"methodResponses":[["Vendor/unknown",{"x":1},"0"]],"sessionState":"s"}`,
				probe.ID,
			)
			_ = c.Write(ctx, cws.MessageText, []byte(resp))
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := Dial(ctx, jc)
	require.NoError(t, err)
	defer conn.Close()

	req := &jmap.Request{}
	req.Invoke(&core.Echo{Hello: "x"})
	resp, err := conn.Do(ctx, req)
	require.NoError(t, err)
	require.Len(t, resp.Responses, 1)
	u, ok := resp.Responses[0].Args.(*jmap.UnknownResponse)
	require.True(t, ok)
	require.Equal(t, "Vendor/unknown", u.Name)
	require.Contains(t, string(u.Raw), `"x":1`)
}
