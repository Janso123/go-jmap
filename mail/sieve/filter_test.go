package sieve

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:fix inline
func boolPtr(v bool) *bool { return new(v) }

func TestFilterConditionMarshal(t *testing.T) {
	data, err := jsonv2.Marshal(&FilterCondition{
		Name:     "vacation",
		IsActive: new(true),
	})
	require.NoError(t, err)
	assert.Equal(t, `{"name":"vacation","isActive":true}`, string(data))
}

func TestFilterOperatorMarshal(t *testing.T) {
	data, err := jsonv2.Marshal(And(
		&FilterCondition{Name: "vacation"},
		&FilterCondition{IsActive: new(true)},
	))
	require.NoError(t, err)
	assert.Equal(t,
		`{"operator":"AND","conditions":[{"name":"vacation"},{"isActive":true}]}`,
		string(data))
}

func TestFilterConditionMarshalInactive(t *testing.T) {
	data, err := jsonv2.Marshal(&FilterCondition{
		IsActive: new(false),
	})
	require.NoError(t, err)
	assert.Equal(t, `{"isActive":false}`, string(data))
}
