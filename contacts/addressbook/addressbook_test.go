package addressbook

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressBookJSON(t *testing.T) {
	book := AddressBook{
		ID:           "ab1",
		Name:         "Friends",
		Description:  jmap.Some("People I know"),
		SortOrder:    10,
		IsDefault:    new(true),
		IsSubscribed: new(true),
		ShareWith: jmap.Some(map[jmap.ID]*Rights{
			"principal1": {
				MayRead:  true,
				MayWrite: true,
			},
		}),
		MyRights: &Rights{
			MayRead:   true,
			MayWrite:  true,
			MayShare:  true,
			MayDelete: true,
		},
	}

	data, err := jsonv2.Marshal(book)
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
	      "mayWrite": true,
	      "mayShare": false,
	      "mayDelete": false
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

func TestRightsFalseBoolsRemarshaled(t *testing.T) {
	t.Parallel()
	const in = `{"mayRead":true,"mayWrite":false,"mayShare":false,"mayDelete":false}`
	var rights Rights
	require.NoError(t, jsonv2.Unmarshal([]byte(in), &rights))
	data, err := jsonv2.Marshal(&rights)
	require.NoError(t, err)
	assert.JSONEq(t, in, string(data))
}

func TestShareWithNull(t *testing.T) {
	t.Parallel()
	var book AddressBook
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"shareWith":null}`), &book))
	require.True(t, book.ShareWith.IsNull())
	data, err := jsonv2.Marshal(&book)
	require.NoError(t, err)
	require.Contains(t, string(data), `"shareWith":null`)
}

func TestAddressBookCreateOmitsIsDefault(t *testing.T) {
	t.Parallel()
	data, err := jsonv2.Marshal(&AddressBook{Name: "a"})
	require.NoError(t, err)
	require.NotContains(t, string(data), "isDefault")
	require.Contains(t, string(data), `"name":"a"`)
}

func TestChangesResponseMatchesKit(t *testing.T) {
	t.Parallel()
	data, err := jsonv2.Marshal(ChangesResponse{NewState: "s"})
	require.NoError(t, err)
	kit, err := jsonv2.Marshal(jmap.ChangesResponse{NewState: "s"})
	require.NoError(t, err)
	assert.JSONEq(t, string(kit), string(data))
}
