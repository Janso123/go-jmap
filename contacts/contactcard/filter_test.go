package contactcard

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterConditionJSONNameGiven(t *testing.T) {
	filter := &FilterCondition{
		NameGiven: "Ada",
	}

	data, err := json.Marshal(filter)
	require.NoError(t, err)
	assert.Equal(t, `{"name/given":"Ada"}`, string(data))
}
