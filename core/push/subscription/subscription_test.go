package subscription_test

import (
	"encoding/json"
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
		Keys: &subscription.Key{
			Public: "BNcRd...",
			Auth:   "tBH...",
		},
		VerificationCode: "abc123",
		Expires:          &expires,
		Types:            []string{"Email", "Mailbox"},
	}

	data, err := json.Marshal(sub)
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
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, sub.ID, got.ID)
	require.Equal(t, sub.DeviceClientID, got.DeviceClientID)
	require.Equal(t, sub.URL, got.URL)
	require.Equal(t, sub.Keys.Public, got.Keys.Public)
	require.Equal(t, sub.Keys.Auth, got.Keys.Auth)
	require.Equal(t, sub.VerificationCode, got.VerificationCode)
	require.Equal(t, expires.UTC(), got.Expires.UTC())
	require.Equal(t, sub.Types, got.Types)
}

func TestPushVerificationJSON(t *testing.T) {
	v := subscription.Verification{
		Type:           "PushVerification",
		SubscriptionID: "ps1",
		Code:           "verify-me",
	}
	data, err := json.Marshal(v)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"@type": "PushVerification",
		"pushSubscriptionId": "ps1",
		"verificationCode": "verify-me"
	}`, string(data))

	var got subscription.Verification
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, v, got)
}

func TestPushSubscriptionGetResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["PushSubscription/get",{"list":[],"notFound":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*subscription.GetResponse)
	require.True(t, ok)
}

func TestPushSubscriptionSetResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["PushSubscription/set",{"created":{},"destroyed":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*subscription.SetResponse)
	require.True(t, ok)
}
