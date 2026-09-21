package jmap_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

type bareFilter struct {
	jmap.FilterBase
	Name string `json:"name"`
}

func TestAndOrNotMarshal(t *testing.T) {
	f := jmap.And(bareFilter{Name: "a"}, bareFilter{Name: "b"})
	b, err := jsonv2.Marshal(f)
	require.NoError(t, err)
	require.Contains(t, string(b), `"operator":"AND"`)
}

func TestDescComparator(t *testing.T) {
	c := jmap.Desc("receivedAt")
	require.Equal(t, "receivedAt", c.Property)
	require.False(t, c.IsAscending)
}

func TestBoolPtrFalse(t *testing.T) {
	require.False(t, *jmap.Bool(false))
}
