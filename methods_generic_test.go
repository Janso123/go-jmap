package jmap_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

type stubObj struct{}

func (stubObj) JMAPType() string     { return "Stub" }
func (stubObj) Requires() []jmap.URI { return nil }

func TestGetNameFromObject(t *testing.T) {
	m := &jmap.Get[stubObj]{}
	require.Equal(t, "Stub/get", m.Name())
}

func TestKitArgumentNames(t *testing.T) {
	t.Parallel()
	g := &jmap.Get[stubObj]{Account: "a", IDs: []jmap.ID{"1"}}
	b, err := jsonv2.Marshal(g)
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"a","ids":["1"]}`, string(b))

	s := &jmap.Set[stubObj]{Account: "a", Destroy: []jmap.ID{"x"}}
	b, err = jsonv2.Marshal(s)
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"a","destroy":["x"]}`, string(b))

	ref := &jmap.ResultReference{ResultOf: "0", Name: "Stub/query", Path: "/ids"}
	s.ReferenceDestroy = ref
	b, err = jsonv2.Marshal(s)
	require.NoError(t, err)
	require.Contains(t, string(b), `"#destroy"`)

	q := &jmap.Query[stubObj]{Account: "a", Limit: jmap.Uint64Ptr(0), ReferenceAnchor: ref}
	b, err = jsonv2.Marshal(q)
	require.NoError(t, err)
	require.Contains(t, string(b), `"limit":0`)
	require.Contains(t, string(b), `"#anchor"`)

	qc := &jmap.QueryChanges[stubObj]{Account: "a", ReferenceUpToID: ref}
	b, err = jsonv2.Marshal(qc)
	require.NoError(t, err)
	require.Contains(t, string(b), `"#upToId"`)
}
