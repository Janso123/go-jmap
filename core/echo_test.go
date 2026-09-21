package core

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	echo := &Echo{
		Hello: "world",
	}

	req := jmap.Request{}
	req.Invoke(echo)
	data, err := json.Marshal(req)
	assert.NoError(t, err)
	exp := `{"using":[],"methodCalls":[["Core/echo",{"Hello":"world"},"0"]]}`
	assert.Equal(t, exp, string(data))
}
