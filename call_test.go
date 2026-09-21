package jmap_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestByCallIDAndAs(t *testing.T) {
	resp := &jmap.Response{Responses: []*jmap.Invocation{
		{Name: "error", Args: &jmap.MethodError{Type: "serverFail"}, CallID: "c0"},
		{Name: "Core/echo", Args: map[string]any{"ok": true}, CallID: "c1"},
	}}
	inv, ok := resp.ByCallID("c1")
	require.True(t, ok)
	require.Equal(t, "Core/echo", inv.Name)

	_, err := jmap.As[map[string]any](resp, "c0")
	require.Error(t, err)
}

func TestCreationRefAndPaths(t *testing.T) {
	require.Equal(t, jmap.ID("#c0"), jmap.CreationRef("c0"))
	require.Equal(t, "/ids", jmap.PathIDs)
	require.Equal(t, "/list/*/subject", jmap.ListPath("subject"))
}

func TestSetCreatedOrError(t *testing.T) {
	created := map[jmap.ID]string{"c0": "ok"}
	not := map[jmap.ID]*jmap.SetError{"c1": {Type: string(jmap.SetErrNotFound)}}
	v, err := jmap.SetCreated(created, not, "c0")
	require.NoError(t, err)
	require.Equal(t, "ok", v)
	_, err = jmap.SetCreated(created, not, "c1")
	require.Error(t, err)
}
