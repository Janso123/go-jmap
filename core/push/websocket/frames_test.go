package websocket

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/core"
	_ "git.sr.ht/~rockorager/go-jmap/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalWSRequest(t *testing.T) {
	req := &jmap.Request{}
	req.Invoke(&core.Echo{Hello: "world"})
	raw, err := marshalRequest("R1", req)
	require.NoError(t, err)
	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &m))
	assert.Equal(t, "Request", m["@type"])
	assert.Equal(t, "R1", m["id"])
	assert.NotNil(t, m["methodCalls"])
	assert.NotNil(t, m["using"])
}

func TestMarshalPushEnableDisable(t *testing.T) {
	en, err := json.Marshal(PushEnable{
		Type:      "WebSocketPushEnable",
		DataTypes: []jmap.EventType{"Email", "Mailbox"},
		PushState: "aaa",
	})
	require.NoError(t, err)
	assert.Contains(t, string(en), `"WebSocketPushEnable"`)
	assert.Contains(t, string(en), `"pushState":"aaa"`)

	dis, err := json.Marshal(PushDisable{Type: "WebSocketPushDisable"})
	require.NoError(t, err)
	assert.Contains(t, string(dis), `"WebSocketPushDisable"`)
}

func TestDecodeServerFrames(t *testing.T) {
	respRaw := `{"@type":"Response","requestId":"R1","methodResponses":[["Core/echo",{"Hello":"world"},"0"]],"sessionState":"s"}`
	fr, err := decodeServerFrame([]byte(respRaw))
	require.NoError(t, err)
	require.NotNil(t, fr.Response)
	assert.Equal(t, "R1", fr.RequestID)
	assert.Equal(t, "s", fr.Response.SessionState)
	require.Len(t, fr.Response.Responses, 1)
	assert.Equal(t, "Core/echo", fr.Response.Responses[0].Name)

	scRaw := `{"@type":"StateChange","changed":{"a":{"Email":"1"}},"pushState":"p2"}`
	fr, err = decodeServerFrame([]byte(scRaw))
	require.NoError(t, err)
	require.NotNil(t, fr.StateChange)
	assert.Equal(t, "p2", fr.StateChange.PushState)
	assert.Equal(t, "1", fr.StateChange.Changed["a"]["Email"])

	errRaw := `{"@type":"RequestError","requestId":"R2","type":"urn:ietf:params:jmap:error:notJSON","status":400,"detail":"bad"}`
	fr, err = decodeServerFrame([]byte(errRaw))
	require.NoError(t, err)
	require.NotNil(t, fr.RequestError)
	assert.Equal(t, "R2", fr.RequestID)
	assert.Equal(t, 400, fr.RequestError.Status)
	assert.Equal(t, "bad", fr.RequestError.Detail)
}

func TestDecodeUnknownFrame(t *testing.T) {
	_, err := decodeServerFrame([]byte(`{"@type":"Nope"}`))
	require.Error(t, err)
}
