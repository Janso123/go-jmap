package principal_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap/sharing"
	"github.com/Janso123/go-jmap/sharing/principal"
	"github.com/stretchr/testify/require"
)

func TestPrincipalOptionalStringsNullRoundTrip(t *testing.T) {
	t.Parallel()
	var p principal.Principal
	raw := `{"id":"p1","type":"individual","name":"Ada","description":null,"email":null,"timeZone":null}`
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &p))
	require.True(t, p.Description.IsNull())
	require.True(t, p.Email.IsNull())
	require.True(t, p.TimeZone.IsNull())

	b, err := jsonv2.Marshal(&p)
	require.NoError(t, err)
	require.Contains(t, string(b), `"description":null`)
	require.Contains(t, string(b), `"email":null`)
	require.Contains(t, string(b), `"timeZone":null`)
}

func TestPrincipalOptionalStringsEmptyDistinctFromNull(t *testing.T) {
	t.Parallel()
	var p principal.Principal
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"description":"","email":"","timeZone":""}`), &p))
	for _, field := range []struct {
		name string
		opt  interface{ IsNull() bool }
	}{
		{"description", p.Description},
		{"email", p.Email},
		{"timeZone", p.TimeZone},
	} {
		require.False(t, field.opt.IsNull(), field.name)
	}
	desc, ok := p.Description.Value()
	require.True(t, ok)
	require.Equal(t, "", desc)
	email, ok := p.Email.Value()
	require.True(t, ok)
	require.Equal(t, "", email)
	tz, ok := p.TimeZone.Value()
	require.True(t, ok)
	require.Equal(t, "", tz)
}

func TestCapabilityCurrentUserPrincipalIDNull(t *testing.T) {
	t.Parallel()
	var c sharing.AccountCapability
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"currentUserPrincipalId":null}`), &c))
	require.True(t, c.CurrentUserPrincipalID.IsNull())

	b, err := jsonv2.Marshal(&c)
	require.NoError(t, err)
	require.JSONEq(t, `{"currentUserPrincipalId":null}`, string(b))
}

func TestPrincipalAccountsNull(t *testing.T) {
	t.Parallel()
	var p principal.Principal
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accounts":null}`), &p))
	require.True(t, p.Accounts.IsNull())

	b, err := jsonv2.Marshal(&p)
	require.NoError(t, err)
	require.JSONEq(t, `{"accounts":null}`, string(b))
}
