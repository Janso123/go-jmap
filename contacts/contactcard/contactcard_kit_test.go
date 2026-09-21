package contactcard_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts"
	"github.com/Janso123/go-jmap/contacts/contactcard"
	"github.com/stretchr/testify/require"
)

func TestContactCardIsObject(t *testing.T) {
	var _ jmap.Object = contactcard.ContactCard{}
	require.Equal(t, "ContactCard", contactcard.ContactCard{}.JMAPType())
	require.Equal(t, []jmap.URI{contacts.URI}, contactcard.ContactCard{}.Requires())
}

func TestContactCardMethodNames(t *testing.T) {
	require.Equal(t, "ContactCard/get", (&contactcard.Get{}).Name())
	require.Equal(t, "ContactCard/changes", (&contactcard.Changes{}).Name())
	require.Equal(t, "ContactCard/query", (&contactcard.Query{}).Name())
	require.Equal(t, "ContactCard/queryChanges", (&contactcard.QueryChanges{}).Name())
	require.Equal(t, "ContactCard/set", (&contactcard.Set{}).Name())
	require.Equal(t, "ContactCard/copy", (&contactcard.Copy{}).Name())
}
