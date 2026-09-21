package email_test

import (
	jsonv2 "encoding/json/v2"
	"testing"
	"time"

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

func TestHeaderPropForms(t *testing.T) {
	t.Parallel()
	require.Equal(t, "header:From:asText", email.HeaderProp("From", email.FormText, false))
	require.Equal(t, "header:From:asAddresses:all", email.HeaderProp("From", email.FormAddresses, true))
	require.Equal(t, "header:Date:asDate", email.HeaderProp("Date", email.FormDate, false))
	require.Equal(t, "header:Message-ID:asMessageIds", email.HeaderProp("Message-ID", email.FormMessageIds, false))
	require.Equal(t, "header:List-Unsubscribe", email.HeaderProp("List-Unsubscribe", email.FormRaw, false))
}

func TestEmailHeaderFormsEmbed(t *testing.T) {
	t.Parallel()
	raw := `{
		"id":"e1",
		"header:Subject:asText":"Hello",
		"header:From:asAddresses":[{"name":"Ada","email":"ada@ex"}],
		"header:To:asGroupedAddresses":[{"name":"Team","addresses":[{"email":"a@ex"}]}],
		"header:Date:asDate":"2026-09-21T10:00:00Z",
		"header:Message-ID:asMessageIds":["<id@ex>"],
		"header:List-Post:asURLs:all":["mailto:list@ex"]
	}`
	var e email.Email
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &e))

	subj, ok := e.HeaderText("Subject")
	require.True(t, ok)
	require.Equal(t, "Hello", subj)

	from, ok := e.HeaderAddresses("From")
	require.True(t, ok)
	require.Equal(t, "ada@ex", from[0].Email)

	groups, ok := e.HeaderGroupedAddresses("To")
	require.True(t, ok)
	require.Equal(t, "Team", groups[0].Name)
	require.Equal(t, "a@ex", groups[0].Addresses[0].Email)

	dt, ok := e.HeaderDate("Date")
	require.True(t, ok)
	require.True(t, dt.UTC().Equal(time.Date(2026, 9, 21, 10, 0, 0, 0, time.UTC)))

	ids, ok := e.HeaderMessageIDs("Message-ID")
	require.True(t, ok)
	require.Equal(t, []string{"<id@ex>"}, ids)

	urls, ok := e.HeaderURLs("List-Post")
	require.True(t, ok)
	require.Equal(t, []string{"mailto:list@ex"}, urls)
}
