package addressbook_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts"
	"github.com/Janso123/go-jmap/contacts/addressbook"
	"github.com/stretchr/testify/require"
)

func TestAddressBookIsObject(t *testing.T) {
	var _ jmap.Object = addressbook.AddressBook{}
	require.Equal(t, "AddressBook", addressbook.AddressBook{}.JMAPType())
	require.Equal(t, []jmap.URI{contacts.URI}, addressbook.AddressBook{}.Requires())
}

func TestAddressBookMethodNames(t *testing.T) {
	require.Equal(t, "AddressBook/get", (&addressbook.Get{}).Name())
	require.Equal(t, "AddressBook/changes", (&addressbook.Changes{}).Name())
	require.Equal(t, "AddressBook/set", (&addressbook.Set{}).Name())
}
