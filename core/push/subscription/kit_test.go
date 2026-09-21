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
		Set: jmap.Set[subscription.PushSubscription]{
			Update: map[jmap.ID]jmap.Patch{
				"p1": {"verificationCode": "x"},
			},
		},
	}
	require.Equal(t, "PushSubscription/set", m.Name())
	b, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	require.JSONEq(t, `{"update":{"p1":{"verificationCode":"x"}}}`, string(b))

	var resp subscription.GetResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"list":[{"id":"p1","url":"https://push.example"}],"notFound":[]}`), &resp))
	require.Equal(t, []subscription.PushSubscription{{ID: "p1", URL: "https://push.example"}}, resp.List)
}
