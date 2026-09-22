package mailbox

import (
	jsonv2 "encoding/json/v2"
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
	data, err := jsonv2.Marshal(query)
	assert.NoError(err)
	expected := `{"sort":[{"property":"name"}]}`
	assert.Equal(expected, string(data))
}
