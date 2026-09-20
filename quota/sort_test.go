package quota

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortMarshal(t *testing.T) {
	query := &Query{
		Sort: []*SortComparator{
			{Property: "used"},
		},
	}

	data, err := json.Marshal(query)
	require.NoError(t, err)
	assert.Equal(t, `{"sort":[{"property":"used","isAscending":false}]}`, string(data))
}
