package searchsnippet_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap/mail/searchsnippet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchSnippetJSONRoundTrip(t *testing.T) {
	s := searchsnippet.SearchSnippet{
		Email:   "em1",
		Subject: "Hello <mark>world</mark>",
		Preview: "... <mark>world</mark> ...",
	}

	data, err := json.Marshal(s)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"emailId": "em1",
		"subject": "Hello <mark>world</mark>",
		"preview": "... <mark>world</mark> ..."
	}`, string(data))

	var got searchsnippet.SearchSnippet
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, s, got)
}
