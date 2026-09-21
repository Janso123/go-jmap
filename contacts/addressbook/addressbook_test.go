package addressbook

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressBookJSON(t *testing.T) {
	book := AddressBook{
		ID:           "ab1",
		Name:         "Friends",
		Description:  "People I know",
		SortOrder:    10,
		IsDefault:    true,
		IsSubscribed: jmap.Bool(true),
		ShareWith: map[jmap.ID]*Rights{
			"principal1": {
				MayRead:  true,
				MayWrite: true,
			},
		},
		MyRights: &Rights{
			MayRead:   true,
			MayWrite:  true,
			MayShare:  true,
			MayDelete: true,
		},
	}

	data, err := json.Marshal(book)
	require.NoError(t, err)

	assert.JSONEq(t, `{
	  "id": "ab1",
	  "name": "Friends",
	  "description": "People I know",
	  "sortOrder": 10,
	  "isDefault": true,
	  "isSubscribed": true,
	  "shareWith": {
	    "principal1": {
	      "mayRead": true,
	      "mayWrite": true
	    }
	  },
	  "myRights": {
	    "mayRead": true,
	    "mayWrite": true,
	    "mayShare": true,
	    "mayDelete": true
	  }
	}`, string(data))
}
