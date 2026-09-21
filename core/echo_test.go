package core

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestEchoRFCExampleRoundTrip(t *testing.T) {
	echo := Echo{"hello": true, "high": 5}
	req := jmap.Request{}
	req.Invoke(echo)
	data, err := jsonv2.Marshal(&req)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"using":["urn:ietf:params:jmap:core"],
		"methodCalls":[["Core/echo",{"hello":true,"high":5},"0"]]
	}`, string(data))

	raw := `{"hello":true,"high":5}`
	var got Echo
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &got))
	require.Equal(t, true, got["hello"])
}
