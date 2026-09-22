package emailsubmission_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap/mail/emailsubmission"
	"github.com/stretchr/testify/require"
)

func roundTrip[T any](t *testing.T, in string) string {
	t.Helper()
	var v T
	require.NoError(t, jsonv2.Unmarshal([]byte(in), &v))
	b, err := jsonv2.Marshal(&v)
	require.NoError(t, err)
	return string(b)
}

func TestSubmissionCapabilityZeroMaxDelayedSend(t *testing.T) {
	t.Parallel()
	out := roundTrip[emailsubmission.Capability](t, `{"maxDelayedSend":0}`)
	require.Contains(t, out, `"maxDelayedSend":0`)
}

func TestEmailSubmissionNullDeliveryStatus(t *testing.T) {
	t.Parallel()
	out := roundTrip[emailsubmission.EmailSubmission](t, `{"deliveryStatus":null}`)
	require.Contains(t, out, `"deliveryStatus":null`)
}
