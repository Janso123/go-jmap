package email_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func roundTrip[T any](t *testing.T, in string) string {
	t.Helper()
	var v T
	require.NoError(t, jsonv2.Unmarshal([]byte(in), &v))
	b, err := jsonv2.Marshal(&v)
	require.NoError(t, err)
	return string(b)
}

func TestEmailWireRoundTrip(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "null subject with zero size and false hasAttachment",
			in:   `{"subject":null,"size":0,"hasAttachment":false}`,
			want: []string{`"subject":null`, `"size":0`, `"hasAttachment":false`},
		},
		{
			name: "null header slices",
			in:   `{"messageId":null,"from":null}`,
			want: []string{`"messageId":null`, `"from":null`},
		},
		{
			name: "null smime fields",
			in:   `{"smimeStatus":null,"smimeErrors":null,"smimeVerifiedAt":null}`,
			want: []string{`"smimeStatus":null`, `"smimeErrors":null`, `"smimeVerifiedAt":null`},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out := roundTrip[email.Email](t, tc.in)
			for _, w := range tc.want {
				require.Contains(t, out, w)
			}
		})
	}
}

func TestBodyPartWireRoundTrip(t *testing.T) {
	t.Parallel()
	in := `{"partId":null,"blobId":null,"size":0,"type":"multipart/mixed","subParts":null}`
	out := roundTrip[email.BodyPart](t, in)
	for _, w := range []string{
		`"partId":null`,
		`"blobId":null`,
		`"size":0`,
		`"type":"multipart/mixed"`,
		`"subParts":null`,
	} {
		require.Contains(t, out, w)
	}
}

func TestAddressNameNullRoundTrip(t *testing.T) {
	t.Parallel()
	out := roundTrip[mail.Address](t, `{"name":null,"email":"a@b"}`)
	require.Contains(t, out, `"name":null`)
	require.Contains(t, out, `"email":"a@b"`)
}

func TestHeaderEmptyValueRoundTrip(t *testing.T) {
	t.Parallel()
	out := roundTrip[email.Header](t, `{"name":"X","value":""}`)
	require.Contains(t, out, `"value":""`)
}

func TestBodyValueFalseFlagsRoundTrip(t *testing.T) {
	t.Parallel()
	out := roundTrip[email.BodyValue](t, `{"value":"x"}`)
	require.Contains(t, out, `"isEncodingProblem":false`)
	require.Contains(t, out, `"isTruncated":false`)
}

func TestEmailCreateOmitsServerSetProperties(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(&email.Set{
		Set: jmap.Set[email.Email]{
			Create: jmap.Some(map[jmap.ID]*email.Email{"c1": {Subject: jmap.Some("s")}}),
		},
	})
	require.NoError(t, err)
	out := string(b)
	require.Contains(t, out, `"subject":"s"`)
	require.NotContains(t, out, "size")
	require.NotContains(t, out, "hasAttachment")
}

func TestEmailHeaderAccessorNull(t *testing.T) {
	t.Parallel()
	var e email.Email
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"header:Subject:asText":null}`), &e))
	got := e.HeaderText("Subject")
	require.True(t, got.IsNull())

	require.True(t, e.HeaderText("Missing").IsZero())
}

func TestEmailParseResponseNullArrays(t *testing.T) {
	t.Parallel()
	var r email.ParseResponse
	in := `{"accountId":"a","parsed":null,"notParsable":null,"notFound":null}`
	require.NoError(t, jsonv2.Unmarshal([]byte(in), &r))
	require.True(t, r.Parsed.IsNull())
	require.True(t, r.NotParsable.IsNull())
	require.True(t, r.NotFound.IsNull())
	b, err := jsonv2.Marshal(&r)
	require.NoError(t, err)
	for _, w := range []string{`"parsed":null`, `"notParsable":null`, `"notFound":null`} {
		require.Contains(t, string(b), w)
	}
}
