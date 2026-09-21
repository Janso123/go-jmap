package email_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func TestEmailImportParseResponsesRegistered(t *testing.T) {
	t.Parallel()
	raw := []byte(`{"sessionState":"s","methodResponses":[
		["Email/import",{"accountId":"a1","created":{"c1":{"id":"e1"}},"notCreated":{}},"0"],
		["Email/parse",{"accountId":"a1","parsed":{"b1":{"id":"e2","subject":"hi"}},"notParsable":[],"notFound":[]},"1"]
	]}`)
	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 2)

	imp, ok := resp.Responses[0].Args.(*email.ImportResponse)
	require.True(t, ok)
	require.Equal(t, jmap.ID("e1"), imp.Created[jmap.ID("c1")].ID)

	par, ok := resp.Responses[1].Args.(*email.ParseResponse)
	require.True(t, ok)
	require.Equal(t, "hi", par.Parsed[jmap.ID("b1")].Subject)
}
