package mdn_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/mdn"
	"github.com/stretchr/testify/require"
)

func TestMDNMethodNames(t *testing.T) {
	require.Equal(t, "MDN/send", (&mdn.Send{}).Name())
	require.Equal(t, "MDN/parse", (&mdn.Parse{}).Name())
}

func TestMDNSendRequiresMailAndMDN(t *testing.T) {
	require.Equal(t, []jmap.URI{mail.URI, mdn.URI}, (&mdn.Send{}).Requires())
}

func TestMDNSendResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["MDN/send",{"accountId":"u1","sent":{},"notSent":{}},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*mdn.SendResponse)
	require.True(t, ok)
}
