package contactcard

import (
	jsonv2 "encoding/json/v2"
	"strings"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts/jscontact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactCardJSON(t *testing.T) {
	card := ContactCard{
		ID: "c1",
		AddressBookIDs: map[jmap.ID]bool{
			"ab1": true,
		},
		UID: "urn:uuid:ada",
		Name: &jscontact.Name{
			Full: "Ada Lovelace",
		},
	}

	data, err := jsonv2.Marshal(card)
	require.NoError(t, err)

	assert.JSONEq(t, `{
	  "id": "c1",
	  "addressBookIds": {
	    "ab1": true
	  },
	  "uid": "urn:uuid:ada",
	  "name": {
	    "full": "Ada Lovelace"
	  }
	}`, string(data))
}

func TestContactCardExtraRoundTrip(t *testing.T) {
	const input = `{"id":"c1","addressBookIds":{"ab1":true},"uid":"urn:uuid:test","foo":1}`
	var card ContactCard
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &card))
	require.Equal(t, jmap.ID("c1"), card.ID)
	require.Equal(t, map[jmap.ID]bool{"ab1": true}, card.AddressBookIDs)
	require.Equal(t, "urn:uuid:test", card.UID)
	raw, ok := card.Extra["foo"]
	require.True(t, ok)
	require.Equal(t, "1", string(raw))
	_, hasUID := card.Extra["uid"]
	require.False(t, hasUID)
	_, hasID := card.Extra["id"]
	require.False(t, hasID)
	_, hasAddressBooks := card.Extra["addressBookIds"]
	require.False(t, hasAddressBooks)

	out, err := jsonv2.Marshal(card)
	require.NoError(t, err)
	require.JSONEq(t, input, string(out))
	require.False(t, strings.Contains(string(out), `"Extra"`))
}
