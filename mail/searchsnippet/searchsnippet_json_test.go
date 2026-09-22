package searchsnippet_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/searchsnippet"
	"github.com/stretchr/testify/require"
)

func TestSearchSnippetJSONRoundTrip(t *testing.T) {
	subj := "Hello <mark>world</mark>"
	preview := "... <mark>world</mark> ..."
	s := searchsnippet.SearchSnippet{
		Email:   "em1",
		Subject: jmap.Some(subj),
		Preview: jmap.Some(preview),
	}

	data, err := jsonv2.Marshal(s)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"emailId": "em1",
		"subject": "Hello <mark>world</mark>",
		"preview": "... <mark>world</mark> ..."
	}`, string(data))

	var got searchsnippet.SearchSnippet
	require.NoError(t, jsonv2.Unmarshal(data, &got))
	require.Equal(t, s.Email, got.Email)
	gotSubj, ok := got.Subject.Value()
	require.True(t, ok)
	require.Equal(t, subj, gotSubj)
	gotPreview, ok := got.Preview.Value()
	require.True(t, ok)
	require.Equal(t, preview, gotPreview)
}

func TestSearchSnippetNullSubject(t *testing.T) {
	var got searchsnippet.SearchSnippet
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"emailId":"e1","subject":null,"preview":null}`), &got))
	require.True(t, got.Subject.IsNull())
	require.True(t, got.Preview.IsNull())
}

func TestSearchSnippetGetResponseNotFoundNull(t *testing.T) {
	t.Parallel()
	var got searchsnippet.GetResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a","list":[],"notFound":null}`), &got))
	require.True(t, got.NotFound.IsNull())

	b, err := jsonv2.Marshal(&got)
	require.NoError(t, err)
	require.Contains(t, string(b), `"notFound":null`)
}

func TestSearchSnippetEncodeJSONNull(t *testing.T) {
	t.Parallel()
	var got searchsnippet.SearchSnippet
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"emailId":"e1","subject":null,"preview":null}`), &got))
	b, err := jsonv2.Marshal(got)
	require.NoError(t, err)
	require.Contains(t, string(b), `"subject":null`)
	require.Contains(t, string(b), `"preview":null`)
}
