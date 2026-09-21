package mailbox

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	assert := assert.New(t)
	query := &Query{
		Sort: []*jmap.Comparator{
			{
				Property: "name",
			},
		},
	}
	data, err := json.Marshal(query)
	assert.NoError(err)
	expected := `{"sort":[{"property":"name","isAscending":false}]}`
	assert.Equal(expected, string(data))
}
