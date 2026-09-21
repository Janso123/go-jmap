package email_test

import (
	"testing"
	"time"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func TestSMIMEVerifyCapabilityURI(t *testing.T) {
	t.Parallel()
	require.Equal(t, jmap.URI("urn:ietf:params:jmap:smimeverify"), email.SMIMEVerify)
}

func TestEmailSMIMEFieldsUnmarshal(t *testing.T) {
	t.Parallel()
	raw := `{
		"id":"e1",
		"smimeStatus":"good",
		"smimeStatusAtDelivery":"unknown",
		"smimeErrors":["expired"],
		"smimeVerifiedAt":"2026-09-21T10:00:00Z"
	}`
	var e email.Email
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &e))
	require.Equal(t, "good", e.SMIMEStatus)
	require.Equal(t, "unknown", e.SMIMEStatusAtDelivery)
	require.Equal(t, []string{"expired"}, e.SMIMEErrors)
	require.NotNil(t, e.SMIMEVerifiedAt)
	require.True(t, time.Time(*e.SMIMEVerifiedAt).UTC().Equal(
		time.Date(2026, 9, 21, 10, 0, 0, 0, time.UTC)))
}
