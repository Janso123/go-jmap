package subscription_test

import (
	jsonv2 "encoding/json/v2"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/push/subscription"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPushSubscriptionMethodNames(t *testing.T) {
	require.Equal(t, "PushSubscription/get", (&subscription.Get{}).Name())
	require.Equal(t, "PushSubscription/set", (&subscription.Set{}).Name())
}

func TestPushSubscriptionJSONRoundTrip(t *testing.T) {
	expires := time.Date(2026, 12, 1, 15, 30, 0, 0, time.UTC)
	sub := subscription.PushSubscription{
		ID:             "ps1",
		DeviceClientID: "device-a",
		URL:            "https://push.example/endpoint",
		Keys: jmap.Some(subscription.Key{
			Public: "BNcRd...",
			Auth:   "tBH...",
		}),
		VerificationCode: jmap.Some("abc123"),
		Expires:          jmap.Some(jmap.UTCDate(expires)),
		Types:            jmap.Some([]string{"Email", "Mailbox"}),
	}

	data, err := jsonv2.Marshal(sub)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"id": "ps1",
		"deviceClientId": "device-a",
		"url": "https://push.example/endpoint",
		"keys": {"p256dh": "BNcRd...", "auth": "tBH..."},
		"verificationCode": "abc123",
		"expires": "2026-12-01T15:30:00Z",
		"types": ["Email", "Mailbox"]
	}`, string(data))

	var got subscription.PushSubscription
	require.NoError(t, jsonv2.Unmarshal(data, &got))
	require.Equal(t, sub.ID, got.ID)
	require.Equal(t, sub.DeviceClientID, got.DeviceClientID)
	require.Equal(t, sub.URL, got.URL)
	wantKeys, ok := sub.Keys.Value()
	require.True(t, ok)
	gotKeys, ok := got.Keys.Value()
	require.True(t, ok)
	require.Equal(t, wantKeys.Public, gotKeys.Public)
	require.Equal(t, wantKeys.Auth, gotKeys.Auth)
	wantCode, ok := sub.VerificationCode.Value()
	require.True(t, ok)
	gotCode, ok := got.VerificationCode.Value()
	require.True(t, ok)
	require.Equal(t, wantCode, gotCode)
	gotExp, ok := got.Expires.Value()
	require.True(t, ok)
	require.Equal(t, jmap.UTCDate(expires), gotExp)
	wantTypes, ok := sub.Types.Value()
	require.True(t, ok)
	gotTypes, ok := got.Types.Value()
	require.True(t, ok)
	require.Equal(t, wantTypes, gotTypes)
}

func TestPushSubscriptionTypesNullVsEmpty(t *testing.T) {
	t.Parallel()
	var got subscription.PushSubscription
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"types":null}`), &got))
	require.True(t, got.Types.IsNull())
	b, err := jsonv2.Marshal(got)
	require.NoError(t, err)
	require.Contains(t, string(b), `"types":null`)
	require.NotContains(t, string(b), `"types":[]`)

	require.NoError(t, jsonv2.Unmarshal([]byte(`{"types":[]}`), &got))
	types, ok := got.Types.Value()
	require.True(t, ok)
	require.Empty(t, types)
	b, err = jsonv2.Marshal(got)
	require.NoError(t, err)
	require.JSONEq(t, `{"types":[]}`, string(b))

	require.NoError(t, jsonv2.Unmarshal([]byte(`{"types":["Email"]}`), &got))
	b, err = jsonv2.Marshal(got)
	require.NoError(t, err)
	require.JSONEq(t, `{"types":["Email"]}`, string(b))
}

func TestPushVerificationJSON(t *testing.T) {
	v := subscription.Verification{
		Type:           "PushVerification",
		SubscriptionID: "ps1",
		Code:           "verify-me",
	}
	data, err := jsonv2.Marshal(v)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"@type": "PushVerification",
		"pushSubscriptionId": "ps1",
		"verificationCode": "verify-me"
	}`, string(data))

	var got subscription.Verification
	require.NoError(t, jsonv2.Unmarshal(data, &got))
	require.Equal(t, v, got)
}

func TestPushSubscriptionGetResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["PushSubscription/get",{"list":[],"notFound":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*subscription.GetResponse)
	require.True(t, ok)
}

func TestPushSubscriptionSetResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["PushSubscription/set",{"created":{},"destroyed":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*subscription.SetResponse)
	require.True(t, ok)
}

func TestPushSubscriptionExpiresNonUTCMarshalsZ(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("CEST", 2*3600)
	expires := time.Date(2026, 12, 1, 16, 30, 0, 0, loc) // 14:30Z
	sub := subscription.PushSubscription{ID: "ps1", Expires: jmap.Some(jmap.UTCDate(expires))}
	data, err := jsonv2.Marshal(sub)
	require.NoError(t, err)
	require.Contains(t, string(data), "2026-12-01T14:30:00Z")
	require.NotContains(t, string(data), "+02:00")
}
