package quota

import (
	jsonv2 "encoding/json/v2"
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

	data, err := jsonv2.Marshal(query)
	require.NoError(t, err)
	assert.Equal(t, `{"sort":[{"property":"used"}]}`, string(data))
}
