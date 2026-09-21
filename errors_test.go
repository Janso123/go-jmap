package jmap_test

import (
	"errors"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestMethodErrorConstantCompare(t *testing.T) {
	err := &jmap.MethodError{Type: string(jmap.MethodErrStateMismatch)}
	require.True(t, errors.Is(err, jmap.ErrStateMismatch))
}
