package subscription_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/push/subscription"
	"github.com/stretchr/testify/require"
)

func TestPushSubscriptionGetOmitsAccountID(t *testing.T) {
	b, err := jsonv2.Marshal(&subscription.Get{})
	require.NoError(t, err)
	require.JSONEq(t, `{}`, string(b))
}

func TestPushSubscriptionSetUpdateIsPatch(t *testing.T) {
	m := &subscription.Set{
		Update: jmap.Some(map[jmap.ID]jmap.Patch{
			"p1": {"verificationCode": "x"},
		}),
	}
	require.Equal(t, "PushSubscription/set", m.Name())
	b, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	require.JSONEq(t, `{"update":{"p1":{"verificationCode":"x"}}}`, string(b))

	var resp subscription.GetResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"list":[{"id":"p1","url":"https://push.example"}],"notFound":[]}`), &resp))
	require.Equal(t, []subscription.PushSubscription{{ID: "p1", URL: "https://push.example"}}, resp.List)
}

func TestPushSubscriptionTypesOmittedNullEmpty(t *testing.T) {
	omitted, err := jsonv2.Marshal(subscription.PushSubscription{ID: "ps1"})
	require.NoError(t, err)
	require.JSONEq(t, `{"id":"ps1"}`, string(omitted))
	require.NotContains(t, string(omitted), "types")

	nulls, err := jsonv2.Marshal(subscription.PushSubscription{Types: jmap.Null[[]string]()})
	require.NoError(t, err)
	require.JSONEq(t, `{"types":null}`, string(nulls))

	empty, err := jsonv2.Marshal(subscription.PushSubscription{Types: jmap.Some([]string{})})
	require.NoError(t, err)
	require.JSONEq(t, `{"types":[]}`, string(empty))
}

func TestPushSubscriptionNullExpiresKeysVerificationCode(t *testing.T) {
	raw := `{"expires":null,"keys":null,"verificationCode":null}`
	var got subscription.PushSubscription
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &got))
	require.True(t, got.Expires.IsNull())
	require.True(t, got.Keys.IsNull())
	require.True(t, got.VerificationCode.IsNull())
	require.True(t, got.Types.IsZero())

	b, err := jsonv2.Marshal(got)
	require.NoError(t, err)
	require.JSONEq(t, raw, string(b))
}
