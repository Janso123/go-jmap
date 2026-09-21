package jmap_test

import (
	"errors"
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestRequestErrorTypeURIs(t *testing.T) {
	t.Parallel()
	require.Equal(t, "urn:ietf:params:jmap:error:unknownCapability", jmap.ErrTypeUnknownCapability)
	require.Equal(t, "urn:ietf:params:jmap:error:notJSON", jmap.ErrTypeNotJSON)
	require.Equal(t, "urn:ietf:params:jmap:error:notRequest", jmap.ErrTypeNotRequest)
	require.Equal(t, "urn:ietf:params:jmap:error:limit", jmap.ErrTypeLimit)
}

func TestMethodErrorTypesTable(t *testing.T) {
	t.Parallel()
	cases := []struct {
		c jmap.MethodErrorType
		s string
	}{
		{jmap.MethodErrUnknownMethod, "unknownMethod"},
		{jmap.MethodErrInvalidArguments, "invalidArguments"},
		{jmap.MethodErrStateMismatch, "stateMismatch"},
		{jmap.MethodErrAccountNotFound, "accountNotFound"},
		{jmap.MethodErrUnsupportedFilter, "unsupportedFilter"},
		{jmap.MethodErrUnsupportedSort, "unsupportedSort"},
	}
	for _, tt := range cases {
		err := &jmap.MethodError{Type: string(tt.c)}
		require.Equal(t, tt.s, err.Type)
		b, e := jsonv2.Marshal(err)
		require.NoError(t, e)
		require.Contains(t, string(b), tt.s)
	}
	require.True(t, errors.Is(
		&jmap.MethodError{Type: string(jmap.MethodErrStateMismatch)},
		jmap.ErrStateMismatch,
	))
}

func TestSetErrorTypesRoundTrip(t *testing.T) {
	t.Parallel()
	props := []string{"keywords"}
	desc := "bad"
	se := &jmap.SetError{
		Type:        string(jmap.SetErrInvalidProperties),
		Description: &desc,
		Properties:  &props,
	}
	b, err := jsonv2.Marshal(se)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"type":"invalidProperties",
		"description":"bad",
		"properties":["keywords"]
	}`, string(b))

	var got jmap.SetError
	require.NoError(t, jsonv2.Unmarshal(b, &got))
	require.Equal(t, string(jmap.SetErrInvalidProperties), got.Type)
	require.Equal(t, []string{"keywords"}, *got.Properties)
}

func TestRequestErrorAs(t *testing.T) {
	t.Parallel()
	lim := "maxSizeRequest"
	err := error(&jmap.RequestError{
		Type:   jmap.ErrTypeLimit,
		Status: 400,
		Title:  "Limit",
		Detail: "too big",
		Limit:  &lim,
	})
	var re *jmap.RequestError
	require.ErrorAs(t, err, &re)
	require.Equal(t, jmap.ErrTypeLimit, re.Type)
	require.Equal(t, 400, re.Status)
}
