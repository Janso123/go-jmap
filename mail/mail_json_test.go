package mail_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/stretchr/testify/require"
)

func TestMailCapabilityNullLimits(t *testing.T) {
	t.Parallel()
	var cap mail.Mail
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"maxMailboxesPerEmail":null,"maxMailboxDepth":null}`), &cap))
	require.True(t, cap.MaxMailboxesPerEmail.IsNull())
	require.True(t, cap.MaxMailboxDepth.IsNull())

	b, err := jsonv2.Marshal(&cap)
	require.NoError(t, err)
	require.Contains(t, string(b), `"maxMailboxesPerEmail":null`)
	require.Contains(t, string(b), `"maxMailboxDepth":null`)

	require.NoError(t, jsonv2.Unmarshal([]byte(`{"maxMailboxesPerEmail":3,"maxMailboxDepth":10}`), &cap))
	perEmail, ok := cap.MaxMailboxesPerEmail.Value()
	require.True(t, ok)
	require.Equal(t, jmap.UnsignedInt(3), perEmail)
	depth, ok := cap.MaxMailboxDepth.Value()
	require.True(t, ok)
	require.Equal(t, jmap.UnsignedInt(10), depth)
}

func TestMailCapabilityRequiredZeros(t *testing.T) {
	t.Parallel()
	var cap mail.Mail
	in := []byte(`{"maxSizeMailboxName":0,"maxSizeAttachmentsPerEmail":0,"emailQuerySortOptions":[],"mayCreateTopLevelMailbox":false}`)
	require.NoError(t, jsonv2.Unmarshal(in, &cap))

	b, err := jsonv2.Marshal(&cap)
	require.NoError(t, err)
	got := string(b)
	require.Contains(t, got, `"maxSizeMailboxName":0`)
	require.Contains(t, got, `"maxSizeAttachmentsPerEmail":0`)
	require.Contains(t, got, `"emailQuerySortOptions":[]`)
	require.Contains(t, got, `"mayCreateTopLevelMailbox":false`)
}
