package mailbox_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/mailbox"
	"github.com/stretchr/testify/require"
)

func TestMailboxCreateTopLevelParentNull(t *testing.T) {
	m := mailbox.Mailbox{Name: "x", ParentID: nil}
	b, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	require.Contains(t, string(b), `"parentId":null`)
}

func TestMailboxIsSubscribedFalse(t *testing.T) {
	m := mailbox.Mailbox{IsSubscribed: jmap.Bool(false)}
	b, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	require.Contains(t, string(b), `"isSubscribed":false`)
}
