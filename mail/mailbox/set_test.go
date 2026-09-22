package mailbox

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	assert := assert.New(t)
	set := &Set{
		Account: "xyz",
		Update: jmap.Some(map[jmap.ID]jmap.Patch{
			"mailbox-id": {
				"name": "New Name",
			},
		}),
	}
	data, err := jsonv2.Marshal(set)
	assert.NoError(err)
	expected := `{"accountId":"xyz","update":{"mailbox-id":{"name":"New Name"}}}`
	assert.Equal(expected, string(data))

	set = &Set{
		Account: "xyz",
		Update: jmap.Some(map[jmap.ID]jmap.Patch{
			"mailbox-id": {
				"parentId": nil,
			},
		}),
	}
	data, err = jsonv2.Marshal(set)
	assert.NoError(err)
	expected = `{"accountId":"xyz","update":{"mailbox-id":{"parentId":null}}}`
	assert.Equal(expected, string(data))
}
