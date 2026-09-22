package subscription

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestOpenPushDecryptsVerification(t *testing.T) {
	keys, priv, err := GenerateKeys()
	require.NoError(t, err)
	auth, err := base64.RawURLEncoding.DecodeString(keys.Auth)
	require.NoError(t, err)
	require.Len(t, auth, 16)

	plain := []byte(`{"@type":"PushVerification","pushSubscriptionId":"p","verificationCode":"v"}`)
	body, err := sealPush(plain, priv.PublicKey().Bytes(), auth)
	require.NoError(t, err)

	v, err := OpenPush(http.Header{"Content-Encoding": {"aes128gcm"}}, body, priv, auth)
	require.NoError(t, err)
	got, ok := v.(*Verification)
	require.True(t, ok)
	require.Equal(t, "PushVerification", got.Type)
	require.Equal(t, "p", got.SubscriptionID)
	require.Equal(t, "v", got.Code)
}

func TestOpenPushDecryptsStateChange(t *testing.T) {
	keys, priv, err := GenerateKeys()
	require.NoError(t, err)
	auth, err := base64.RawURLEncoding.DecodeString(keys.Auth)
	require.NoError(t, err)

	plain := []byte(`{"@type":"StateChange","changed":{"acct":{"Email":"s1"}}}`)
	body, err := sealPush(plain, priv.PublicKey().Bytes(), auth)
	require.NoError(t, err)

	v, err := OpenPush(http.Header{"Content-Encoding": {"aes128gcm"}}, body, priv, auth)
	require.NoError(t, err)
	got, ok := v.(*jmap.StateChange)
	require.True(t, ok)
	require.Equal(t, "StateChange", got.Type)
	require.Equal(t, "s1", got.Changed["acct"]["Email"])
}

func TestOpenPushRejectsPlaintextWhenPrivateKeySet(t *testing.T) {
	_, priv, err := GenerateKeys()
	require.NoError(t, err)
	_, err = OpenPush(http.Header{}, []byte(`{"@type":"PushVerification","pushSubscriptionId":"p","verificationCode":"v"}`), priv, []byte("0123456789abcdef"))
	require.EqualError(t, err, "subscription: unencrypted push rejected; keys were set")
}

func TestOpenPushPlaintextWhenNoPrivateKey(t *testing.T) {
	v, err := OpenPush(nil, []byte(`{"@type":"PushVerification","pushSubscriptionId":"p","verificationCode":"v"}`), nil, nil)
	require.NoError(t, err)
	got, ok := v.(*Verification)
	require.True(t, ok)
	require.Equal(t, "v", got.Code)
}
