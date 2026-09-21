package jmap_test

import (
	"testing"

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
