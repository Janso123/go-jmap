package email_test

import (
	"testing"

	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func TestEmailGetEmbedsKitName(t *testing.T) {
	require.Equal(t, "Email/get", (&email.Get{}).Name())
}

func TestEmailManualMethodNames(t *testing.T) {
	require.Equal(t, "Email/import", (&email.Import{}).Name())
	require.Equal(t, "Email/parse", (&email.Parse{}).Name())
}
