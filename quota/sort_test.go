package quota

import (
	"encoding/json"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortMarshal(t *testing.T) {
	query := &Query{
		Sort: []*jmap.Comparator{
			{Property: "used"},
		},
	}

	data, err := json.Marshal(query)
	require.NoError(t, err)
	assert.Equal(t, `{"sort":[{"property":"used","isAscending":false}]}`, string(data))
}
