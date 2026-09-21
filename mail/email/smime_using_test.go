package email_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func TestEmailGetRequiresSMIMEWhenPropertiesAsk(t *testing.T) {
	t.Parallel()
	require.Equal(t, []jmap.URI{mail.URI}, (&email.Get{}).Requires())
	g := &email.Get{}
	g.Properties = []string{email.PropSubject, email.PropSMIMEStatus}
	require.Equal(t, []jmap.URI{mail.URI, email.SMIMEVerify}, g.Requires())
}

func TestEmailQueryRequiresSMIMEWhenFilterSet(t *testing.T) {
	t.Parallel()
	q := &email.Query{Filter: &email.FilterCondition{HasSMIME: jmap.Bool(false)}}
	require.Equal(t, []jmap.URI{mail.URI, email.SMIMEVerify}, q.Requires())
}

func TestEmailSentAtIsRFC8620Date(t *testing.T) {
	t.Parallel()
	raw := `{"sentAt":"2014-10-30T14:12:00+08:00"}`
	var e email.Email
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &e))
	require.NotNil(t, e.SentAt)
	b, err := jsonv2.Marshal(e)
	require.NoError(t, err)
	require.Contains(t, string(b), "2014-10-30T14:12:00+08:00")
}
