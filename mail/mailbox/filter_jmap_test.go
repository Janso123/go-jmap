package mailbox_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/mailbox"
	"github.com/stretchr/testify/require"
)

func TestJmapAndMailboxParent(t *testing.T) {
	f := jmap.And(&mailbox.FilterCondition{Name: "Inbox"})
	b, err := jsonv2.Marshal(f)
	require.NoError(t, err)
	require.JSONEq(t, `{"operator":"AND","conditions":[{"name":"Inbox"}]}`, string(b))
}

func TestMailboxFilterExplicitFalseAndNull(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(&mailbox.FilterCondition{IsSubscribed: jmap.Bool(false)})
	require.NoError(t, err)
	require.JSONEq(t, `{"isSubscribed":false}`, string(b))

	b, err = jsonv2.Marshal(mailbox.TopLevel())
	require.NoError(t, err)
	require.JSONEq(t, `{"parentId":null}`, string(b))

	b, err = jsonv2.Marshal(&mailbox.FilterCondition{Role: jmap.Null[mailbox.Role]()})
	require.NoError(t, err)
	require.JSONEq(t, `{"role":null}`, string(b))

	b, err = jsonv2.Marshal(&mailbox.FilterCondition{HasAnyRole: jmap.Bool(false)})
	require.NoError(t, err)
	require.JSONEq(t, `{"hasAnyRole":false}`, string(b))
}
