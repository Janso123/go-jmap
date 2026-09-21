package sieve

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:fix inline
func boolPtr(v bool) *bool { return new(v) }

func TestFilterConditionMarshal(t *testing.T) {
	data, err := json.Marshal(&FilterCondition{
		Name:     "vacation",
		IsActive: new(true),
	})
	require.NoError(t, err)
	assert.Equal(t, `{"name":"vacation","isActive":true}`, string(data))
}

func TestFilterOperatorMarshal(t *testing.T) {
	data, err := json.Marshal(And(
		&FilterCondition{Name: "vacation"},
		&FilterCondition{IsActive: new(true)},
	))
	require.NoError(t, err)
	assert.Equal(t,
		`{"operator":"AND","conditions":[{"name":"vacation"},{"isActive":true}]}`,
		string(data))
}

func TestFilterConditionMarshalInactive(t *testing.T) {
	data, err := json.Marshal(&FilterCondition{
		IsActive: new(false),
	})
	require.NoError(t, err)
	assert.Equal(t, `{"isActive":false}`, string(data))
}
