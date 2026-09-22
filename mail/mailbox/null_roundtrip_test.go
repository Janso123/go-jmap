package mailbox_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/mailbox"
	"github.com/stretchr/testify/require"
)

func roundTrip[T any](t *testing.T, in string) string {
	t.Helper()
	var v T
	require.NoError(t, jsonv2.Unmarshal([]byte(in), &v))
	b, err := jsonv2.Marshal(&v)
	require.NoError(t, err)
	return string(b)
}

func TestMailboxWireRoundTrip(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		in     string
		want   []string
		absent []string
	}{
		{
			name: "explicit nulls survive",
			in:   `{"parentId":null,"name":"Inbox","role":null}`,
			want: []string{`"parentId":null`, `"role":null`},
		},
		{
			name:   "unset stays omitted",
			in:     `{"name":"Inbox"}`,
			absent: []string{`"parentId"`, `"role"`},
		},
		{
			name: "server-set zeroes survive",
			in:   `{"totalEmails":0,"unreadEmails":0,"totalThreads":0,"unreadThreads":0,"sortOrder":0}`,
			want: []string{
				`"sortOrder":0`,
				`"totalEmails":0`,
				`"unreadEmails":0`,
				`"totalThreads":0`,
				`"unreadThreads":0`,
			},
		},
		{
			name: "false rights survive",
			in:   `{"myRights":{"mayReadItems":false}}`,
			want: []string{`"mayReadItems":false`},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out := roundTrip[mailbox.Mailbox](t, tc.in)
			for _, w := range tc.want {
				require.Contains(t, out, w)
			}
			for _, a := range tc.absent {
				require.NotContains(t, out, a)
			}
		})
	}
}

func TestMailboxChangesResponseWireRoundTrip(t *testing.T) {
	t.Parallel()
	in := `{"updatedProperties":null,"created":[],"updated":[],"destroyed":[],"hasMoreChanges":false}`
	out := roundTrip[mailbox.ChangesResponse](t, in)
	require.Contains(t, out, `"updatedProperties":null`)
	require.Contains(t, out, `"created":[]`)
	require.Contains(t, out, `"updated":[]`)
	require.Contains(t, out, `"destroyed":[]`)
	require.Contains(t, out, `"hasMoreChanges":false`)
}

func TestMailboxCreateOmitsServerSetProperties(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(&mailbox.Set{
		Set: jmap.Set[mailbox.Mailbox]{
			Create: jmap.Some(map[jmap.ID]*mailbox.Mailbox{"c1": {Name: "New"}}),
		},
	})
	require.NoError(t, err)
	out := string(b)
	require.Contains(t, out, `"name":"New"`)
	for _, absent := range []string{"totalEmails", "unreadEmails", "totalThreads", "unreadThreads", "myRights"} {
		require.NotContains(t, out, absent)
	}
}
