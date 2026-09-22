package quota

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterConditionMarshal(t *testing.T) {
	filter := &FilterCondition{
		Name:         "Main quota",
		Scope:        "account",
		ResourceType: "octets",
		Type:         "Email",
	}

	data, err := jsonv2.Marshal(filter)
	require.NoError(t, err)
	assert.Equal(t,
		`{"name":"Main quota","scope":"account","resourceType":"octets","type":"Email"}`,
		string(data))
}
