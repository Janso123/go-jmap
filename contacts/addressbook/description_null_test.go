package addressbook_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap/contacts/addressbook"
	"github.com/stretchr/testify/require"
)

func TestAddressBookDescriptionNullRoundTrip(t *testing.T) {
	t.Parallel()
	var book addressbook.AddressBook
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"ab","name":"Friends","description":null}`), &book))
	require.True(t, book.Description.IsNull())
	_, ok := book.Description.Value()
	require.False(t, ok)

	b, err := jsonv2.Marshal(&book)
	require.NoError(t, err)
	require.Contains(t, string(b), `"description":null`)
}

func TestAddressBookDescriptionEmptyDistinctFromNull(t *testing.T) {
	t.Parallel()
	var book addressbook.AddressBook
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"description":""}`), &book))
	require.False(t, book.Description.IsNull())
	v, ok := book.Description.Value()
	require.True(t, ok)
	require.Equal(t, "", v)
}
