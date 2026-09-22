package jmap_test

import (
	"testing"

	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

type bareFilter struct {
	jmap.FilterBase
	Name string `json:"name"`
}

type cond struct {
	jmap.FilterBase
	Name string `json:"name,omitzero"`
}

func TestUnmarshalFilterTree(t *testing.T) {
	f, err := jmap.UnmarshalFilter[cond](jsontext.Value(`{"operator":"AND","conditions":[{"name":"a"},{"operator":"NOT","conditions":[{"name":"b"}]}]}`))
	require.NoError(t, err)
	op := f.(*jmap.FilterOperator)
	require.Equal(t, jmap.OperatorAND, op.Operator)
	require.Equal(t, "a", op.Conditions[0].(*cond).Name)
	inner := op.Conditions[1].(*jmap.FilterOperator)
	require.Equal(t, "b", inner.Conditions[0].(*cond).Name)

	f, err = jmap.UnmarshalFilter[cond](jsontext.Value(`null`))
	require.NoError(t, err)
	require.Nil(t, f)
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
	require.NotNil(t, c.IsAscending)
	require.False(t, *c.IsAscending)
}

func TestComparatorZeroOmitsIsAscending(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(&jmap.Comparator{Property: "name"})
	require.NoError(t, err)
	require.JSONEq(t, `{"property":"name"}`, string(b))

	b, err = jsonv2.Marshal(jmap.Desc("name"))
	require.NoError(t, err)
	require.JSONEq(t, `{"property":"name","isAscending":false}`, string(b))
}

func TestBoolPtrFalse(t *testing.T) {
	require.False(t, *new(false))
}
