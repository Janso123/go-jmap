package contactcard

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/contacts/jscontact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactCardJSON(t *testing.T) {
	card := ContactCard{
		ID: "c1",
		AddressBookIDs: map[jmap.ID]bool{
			"ab1": true,
		},
		Card: jscontact.Card{
			UID: "urn:uuid:ada",
			Name: &jscontact.Name{
				Full: "Ada Lovelace",
			},
		},
	}

	data, err := json.Marshal(card)
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
