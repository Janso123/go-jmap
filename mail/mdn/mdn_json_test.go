package mdn_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/mdn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMDNJSONRoundTrip(t *testing.T) {
	m := mdn.MDN{
		ForEmailID:             jmap.Some(jmap.ID("em1")),
		Subject:                jmap.Some("Read: hello"),
		TextBody:               jmap.Some("The message was displayed."),
		IncludeOriginalMessage: new(true),
		ReportingUA:            jmap.Some("go-jmap; 1.0"),
		Disposition: &mdn.Disposition{
			ActionMode:  "manual-action",
			SendingMode: "mdn-sent-manually",
			Type:        "displayed",
		},
		FinalRecipient:    jmap.Some("rfc822; me@example.com"),
		OriginalMessageID: jmap.Some("<msg@example.com>"),
		Error:             jmap.Some([]string{"smtp; 550"}),
		ExtensionFields:   jmap.Some(map[string]string{"X-Custom": "1"}),
	}

	data, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"forEmailId": "em1",
		"subject": "Read: hello",
		"textBody": "The message was displayed.",
		"includeOriginalMessage": true,
		"reportingUA": "go-jmap; 1.0",
		"disposition": {
			"actionMode": "manual-action",
			"sendingMode": "mdn-sent-manually",
			"type": "displayed"
		},
		"finalRecipient": "rfc822; me@example.com",
		"originalMessageId": "<msg@example.com>",
		"error": ["smtp; 550"],
		"extensionFields": {"X-Custom": "1"}
	}`, string(data))

	var got mdn.MDN
	require.NoError(t, jsonv2.Unmarshal(data, &got))
	require.Equal(t, m.ForEmailID, got.ForEmailID)
	require.Equal(t, m.Subject, got.Subject)
	require.Equal(t, m.Disposition.Type, got.Disposition.Type)
	require.Equal(t, m.FinalRecipient, got.FinalRecipient)
	ext, ok := m.ExtensionFields.Value()
	require.True(t, ok)
	gotExt, ok := got.ExtensionFields.Value()
	require.True(t, ok)
	require.Equal(t, ext["X-Custom"], gotExt["X-Custom"])
}

func TestMDNForEmailIdNull(t *testing.T) {
	t.Parallel()
	var got mdn.MDN
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"forEmailId":null,"subject":"x"}`), &got))
	b, err := jsonv2.Marshal(got)
	require.NoError(t, err)
	require.Contains(t, string(b), `"forEmailId":null`)
}

func TestMDNParseResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["MDN/parse",{"accountId":"u1","parsed":{},"notFound":[],"notParsable":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*mdn.ParseResponse)
	require.True(t, ok)
}
