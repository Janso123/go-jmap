package quota

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Get{Account: "u1"})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:quota"],"methodCalls":[["Quota/get",{"accountId":"u1"},"0"]]}`,
		string(data))
}
