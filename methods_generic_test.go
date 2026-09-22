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
func (stubObj) JMAPCreatable()       {}

func TestGetNameFromObject(t *testing.T) {
	m := &jmap.Get[stubObj]{}
	require.Equal(t, "Stub/get", m.Name())
}

func TestKitArgumentNames(t *testing.T) {
	t.Parallel()
	g := &jmap.Get[stubObj]{Account: "a", IDs: jmap.Some([]jmap.ID{"1"})}
	b, err := jsonv2.Marshal(g)
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"a","ids":["1"]}`, string(b))

	s := &jmap.Set[stubObj]{Account: "a", Destroy: jmap.Some([]jmap.ID{"x"})}
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

func TestSetResponseOldStateNull(t *testing.T) {
	t.Parallel()
	var resp jmap.SetResponse[stubObj]
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"oldState":null,"newState":"n1"}`), &resp))
	b, err := jsonv2.Marshal(resp)
	require.NoError(t, err)
	require.Contains(t, string(b), `"oldState":null`)
	require.Contains(t, string(b), `"newState":"n1"`)
}

func TestCopyResponseOldStateNull(t *testing.T) {
	t.Parallel()
	var resp jmap.CopyResponse[stubObj]
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"oldState":null,"newState":"n1"}`), &resp))
	b, err := jsonv2.Marshal(resp)
	require.NoError(t, err)
	require.Contains(t, string(b), `"oldState":null`)
}

func TestMaxChangesOnWire(t *testing.T) {
	t.Parallel()
	one := jmap.Uint64Ptr(1)

	b, err := jsonv2.Marshal(&jmap.Changes[stubObj]{Account: "a", SinceState: "s", MaxChanges: one})
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"a","sinceState":"s","maxChanges":1}`, string(b))

	b, err = jsonv2.Marshal(&jmap.QueryChanges[stubObj]{Account: "a", SinceQueryState: "q", MaxChanges: one})
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"a","sinceQueryState":"q","maxChanges":1}`, string(b))

	// RFC 8620 §5.6: a supplied maxChanges MUST be > 0.
	qc := &jmap.QueryChanges[stubObj]{Account: "a", SinceQueryState: "q", MaxChanges: new(jmap.UnsignedInt(0))}
	require.Error(t, qc.ValidateArgs())
}

func TestMaxChangesUnsetOmitted(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(&jmap.Changes[stubObj]{Account: "a", SinceState: "s"})
	require.NoError(t, err)
	require.NotContains(t, string(b), "maxChanges")

	b, err = jsonv2.Marshal(&jmap.QueryChanges[stubObj]{Account: "a", SinceQueryState: "q"})
	require.NoError(t, err)
	require.NotContains(t, string(b), "maxChanges")
}

func TestChangesResponseHasMoreChangesFalseOnWire(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(jmap.ChangesResponse{NewState: "s"})
	require.NoError(t, err)
	require.Contains(t, string(b), `"hasMoreChanges":false`)
}

func TestQueryResponseCanCalculateChangesFalseOnWire(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(jmap.QueryResponse{QueryState: "q"})
	require.NoError(t, err)
	require.Contains(t, string(b), `"canCalculateChanges":false`)
}

func TestGetIDsThreeStates(t *testing.T) {
	b, err := jsonv2.Marshal(&jmap.Get[stubObj]{Account: "a"})
	require.NoError(t, err)
	require.NotContains(t, string(b), `"ids"`)

	b, err = jsonv2.Marshal(&jmap.Get[stubObj]{Account: "a", IDs: jmap.Null[[]jmap.ID]()})
	require.NoError(t, err)
	require.Contains(t, string(b), `"ids":null`)

	b, err = jsonv2.Marshal(&jmap.Get[stubObj]{Account: "a", IDs: jmap.Some([]jmap.ID{})})
	require.NoError(t, err)
	require.Contains(t, string(b), `"ids":[]`)
}

func TestQueryResponseZeroFields(t *testing.T) {
	b, err := jsonv2.Marshal(jmap.QueryResponse{})
	require.NoError(t, err)
	require.Contains(t, string(b), `"position":0`)
	require.Contains(t, string(b), `"ids":[]`)
	require.NotContains(t, string(b), `"total"`)

	b, err = jsonv2.Marshal(jmap.QueryResponse{Total: new(jmap.UnsignedInt(0)), Limit: new(jmap.UnsignedInt(0))})
	require.NoError(t, err)
	require.Contains(t, string(b), `"total":0`)
	require.Contains(t, string(b), `"limit":0`)
}

func TestSetIfInStateNull(t *testing.T) {
	b, err := jsonv2.Marshal(&jmap.Set[stubObj]{Account: "a", IfInState: jmap.Null[string]()})
	require.NoError(t, err)
	require.Contains(t, string(b), `"ifInState":null`)
}

func TestMaxChangesZeroRejected(t *testing.T) {
	c := &jmap.Changes[stubObj]{Account: "a", MaxChanges: new(jmap.UnsignedInt(0))}
	require.Error(t, c.ValidateArgs()) // RFC 8620 §5.2: MUST be > 0
}

func TestRefsDuplicateRejected(t *testing.T) {
	g := &jmap.Get[stubObj]{Account: "a", IDs: jmap.Some([]jmap.ID{"1"}),
		Refs: map[string]*jmap.ResultReference{"#ids": {ResultOf: "0", Name: "Stub/query", Path: "/ids"}}}
	require.Error(t, g.ValidateArgs())
}

// kitWrapper stands in for the thin packages: they embed a kit type and add
// their own arguments, which must survive marshaling.
type kitWrapper struct {
	jmap.Get[stubObj]

	BodyProperties []string `json:"bodyProperties,omitzero"`
}

func TestEmbeddedKitKeepsWrapperFields(t *testing.T) {
	t.Parallel()
	w := &kitWrapper{BodyProperties: []string{"subject"}}
	w.Account = "a"
	w.IDs = jmap.Some([]jmap.ID{"1"})

	b, err := jsonv2.Marshal(w)
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"a","ids":["1"],"bodyProperties":["subject"]}`, string(b))
}

func TestRequestValidatesArgs(t *testing.T) {
	t.Parallel()
	req := &jmap.Request{}
	req.Invoke(&jmap.Changes[stubObj]{Account: "a", MaxChanges: new(jmap.UnsignedInt(0))})
	_, err := jsonv2.Marshal(req)
	require.Error(t, err)

	ok := &jmap.Request{}
	ok.Invoke(&jmap.Changes[stubObj]{Account: "a", MaxChanges: new(jmap.UnsignedInt(1))})
	b, err := jsonv2.Marshal(ok)
	require.NoError(t, err)
	require.Contains(t, string(b), `"maxChanges":1`)
}
