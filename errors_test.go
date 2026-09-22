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
	se := &jmap.SetError{
		Type:        string(jmap.SetErrInvalidProperties),
		Description: jmap.Some("bad"),
		Properties:  jmap.Some([]string{"keywords"}),
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
	gotProps, ok := got.Properties.Value()
	require.True(t, ok)
	require.Equal(t, []string{"keywords"}, gotProps)
}

func TestSetErrorRFCExtras(t *testing.T) {
	t.Parallel()
	se := &jmap.SetError{
		Type:       string(jmap.SetErrAlreadyExists),
		ExistingID: "id1",
		NotFound:   []jmap.ID{"b1"},
		MaxSize:    new(jmap.UnsignedInt(1024)),
	}
	b, err := jsonv2.Marshal(se)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"type":"alreadyExists",
		"existingId":"id1",
		"notFound":["b1"],
		"maxSize":1024
	}`, string(b))

	raw := `{"type":"tooManyRecipients","maxRecipients":10,"invalidRecipients":["a@b.c"]}`
	var got jmap.SetError
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &got))
	require.Equal(t, string(jmap.SetErrTooManyRecipients), got.Type)
	require.Equal(t, jmap.UnsignedInt(10), *got.MaxRecipients)
	require.Equal(t, []string{"a@b.c"}, got.InvalidRecipients)
	require.Equal(t, jmap.SetErrInvalidSieve, jmap.SetErrorType("invalidSieve"))
	require.Equal(t, jmap.SetErrSieveIsActive, jmap.SetErrorType("sieveIsActive"))
}

func TestRequestErrorProblemJSON(t *testing.T) {
	t.Parallel()
	raw := `{"type":"about:blank","status":400,"instance":"https://x/e","detail":null,"x-vendor":1}`
	var re jmap.RequestError
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &re))
	inst, ok := re.Instance.Value()
	require.True(t, ok)
	require.Equal(t, "https://x/e", inst)
	require.True(t, re.Detail.IsNull())
	require.Contains(t, re.Extra, "x-vendor")
	require.JSONEq(t, "1", string(re.Extra["x-vendor"]))

	out, err := jsonv2.Marshal(&re)
	require.NoError(t, err)
	require.Contains(t, string(out), `"detail":null`)
	require.Contains(t, string(out), `"x-vendor":1`)

	unset := &jmap.RequestError{Type: "about:blank", Status: 400}
	b, err := jsonv2.Marshal(unset)
	require.NoError(t, err)
	require.NotContains(t, string(b), "detail")
	require.NotContains(t, string(b), "instance")
}

func TestRequestErrorAs(t *testing.T) {
	t.Parallel()
	err := error(&jmap.RequestError{
		Type:   jmap.ErrTypeLimit,
		Status: 400,
		Title:  "Limit",
		Detail: jmap.Some("too big"),
		Limit:  jmap.Some("maxSizeRequest"),
	})
	var re *jmap.RequestError
	require.ErrorAs(t, err, &re)
	require.Equal(t, jmap.ErrTypeLimit, re.Type)
	require.Equal(t, 400, re.Status)
}

func TestSetErrorExtraRoundTrip(t *testing.T) {
	t.Parallel()
	raw := `{"type":"forbidden","x-server":"hint"}`
	var se jmap.SetError
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &se))
	require.Equal(t, "forbidden", se.Type)
	require.Contains(t, se.Extra, "x-server")
	out, err := jsonv2.Marshal(&se)
	require.NoError(t, err)
	require.JSONEq(t, raw, string(out))
}

func TestMethodErrorExtraRoundTrip(t *testing.T) {
	t.Parallel()
	raw := `{"type":"invalidArguments","description":"bad","existingId":"id1"}`
	var me jmap.MethodError
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &me))
	require.Equal(t, "invalidArguments", me.Type)
	require.Contains(t, me.Extra, "existingId")
	out, err := jsonv2.Marshal(&me)
	require.NoError(t, err)
	require.JSONEq(t, raw, string(out))
}
