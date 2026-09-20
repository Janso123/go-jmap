package sieve

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(v bool) *bool { return &v }

func TestFilterConditionMarshal(t *testing.T) {
	data, err := json.Marshal(&FilterCondition{
		Name:     "vacation",
		IsActive: boolPtr(true),
	})
	require.NoError(t, err)
	assert.Equal(t, `{"name":"vacation","isActive":true}`, string(data))
}

func TestFilterOperatorMarshal(t *testing.T) {
	data, err := json.Marshal(&FilterOperator{
		Operator: jmap.OperatorAND,
		Conditions: []Filter{
			&FilterCondition{Name: "vacation"},
			&FilterCondition{IsActive: boolPtr(true)},
		},
	})
	require.NoError(t, err)
	assert.Equal(t,
		`{"operator":"AND","conditions":[{"name":"vacation"},{"isActive":true}]}`,
		string(data))
}

func TestFilterConditionMarshalInactive(t *testing.T) {
	data, err := json.Marshal(&FilterCondition{
		IsActive: boolPtr(false),
	})
	require.NoError(t, err)
	assert.Equal(t, `{"isActive":false}`, string(data))
}
