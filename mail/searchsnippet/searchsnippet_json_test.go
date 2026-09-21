package searchsnippet_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap/mail/searchsnippet"
	"github.com/stretchr/testify/require"
)

func TestSearchSnippetJSONRoundTrip(t *testing.T) {
	subj := "Hello <mark>world</mark>"
	preview := "... <mark>world</mark> ..."
	s := searchsnippet.SearchSnippet{
		Email:   "em1",
		Subject: &subj,
		Preview: &preview,
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
	require.Equal(t, subj, *got.Subject)
	require.Equal(t, preview, *got.Preview)
}

func TestSearchSnippetNullSubject(t *testing.T) {
	var got searchsnippet.SearchSnippet
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"emailId":"e1","subject":null,"preview":null}`), &got))
	require.Nil(t, got.Subject)
	require.Nil(t, got.Preview)
}
