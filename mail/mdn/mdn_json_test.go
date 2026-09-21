package mdn_test

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/mdn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMDNJSONRoundTrip(t *testing.T) {
	m := mdn.MDN{
		ForEmailID:             "em1",
		Subject:                "Read: hello",
		TextBody:               "The message was displayed.",
		IncludeOriginalmessage: true,
		ReportingUA:            "go-jmap; 1.0",
		Disposition: &mdn.Disposition{
			ActionMode:  "manual-action",
			SendingMode: "mdn-sent-manually",
			Type:        "displayed",
		},
		FinalRecipient:    "rfc822; me@example.com",
		OriginalMessageID: "<msg@example.com>",
		Error:             []string{"smtp; 550"},
		ExtensionFields:   map[string]string{"X-Custom": "1"},
	}

	data, err := json.Marshal(m)
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
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, m.ForEmailID, got.ForEmailID)
	require.Equal(t, m.Subject, got.Subject)
	require.Equal(t, m.Disposition.Type, got.Disposition.Type)
	require.Equal(t, m.FinalRecipient, got.FinalRecipient)
	require.Equal(t, m.ExtensionFields["X-Custom"], got.ExtensionFields["X-Custom"])
}

func TestMDNParseResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["MDN/parse",{"accountId":"u1","parsed":{},"notFound":[],"notParsable":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*mdn.ParseResponse)
	require.True(t, ok)
}
