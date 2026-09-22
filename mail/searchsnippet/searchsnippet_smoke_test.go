package searchsnippet_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/Janso123/go-jmap/mail/searchsnippet"
	"github.com/stretchr/testify/require"
)

func TestSearchSnippetGetName(t *testing.T) {
	require.Equal(t, "SearchSnippet/get", (&searchsnippet.Get{}).Name())
}

func TestSearchSnippetGetRequiresMail(t *testing.T) {
	require.Equal(t, []jmap.URI{mail.URI}, (&searchsnippet.Get{}).Requires())
}

func TestSearchSnippetGetRequiresSMIMEWhenFilterSet(t *testing.T) {
	g := &searchsnippet.Get{Filter: &email.FilterCondition{HasSMIME: new(true)}}
	require.Equal(t, []jmap.URI{mail.URI, email.SMIMEVerify}, g.Requires())
}

func TestSearchSnippetGetResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["SearchSnippet/get",{"accountId":"u1","list":[],"notFound":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*searchsnippet.GetResponse)
	require.True(t, ok)
}
