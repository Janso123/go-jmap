package mail_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap/mail"
	"github.com/stretchr/testify/require"
)

func TestMailCapabilityNullLimits(t *testing.T) {
	t.Parallel()
	var cap mail.Mail
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"maxMailboxesPerEmail":null,"maxMailboxDepth":null}`), &cap))
	require.Nil(t, cap.MaxMailboxesPerEmail)
	require.Nil(t, cap.MaxMailboxDepth)

	require.NoError(t, jsonv2.Unmarshal([]byte(`{"maxMailboxesPerEmail":3,"maxMailboxDepth":10}`), &cap))
	require.Equal(t, uint64(3), *cap.MaxMailboxesPerEmail)
	require.Equal(t, uint64(10), *cap.MaxMailboxDepth)
}
