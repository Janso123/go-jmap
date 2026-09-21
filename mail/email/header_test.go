package email_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func TestEmailHeaderURLsEmbed(t *testing.T) {
	raw := `{"id":"e1","header:List-Unsubscribe:asURLs":["https://ex/unsub"]}`
	var e email.Email
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &e))
	urls, ok := e.HeaderURLs("List-Unsubscribe")
	require.True(t, ok)
	require.Equal(t, []string{"https://ex/unsub"}, urls)
}
