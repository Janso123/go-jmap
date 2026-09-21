package email_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap/mail/email"
	"github.com/stretchr/testify/require"
)

func TestEmailKeywordsRoundTripAndAccessors(t *testing.T) {
	t.Parallel()
	raw := `{"id":"e1","keywords":{"$seen":true,"$flagged":true,"custom":false}}`
	var e email.Email
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &e))
	require.True(t, e.IsSeen())
	require.True(t, e.HasKeyword(email.KeywordFlagged))
	require.False(t, e.HasKeyword("custom"))
	require.False(t, e.HasKeyword(email.KeywordDraft))

	out, err := jsonv2.Marshal(&e)
	require.NoError(t, err)
	require.Contains(t, string(out), `"$seen":true`)
	require.Contains(t, string(out), `"$flagged":true`)
}
