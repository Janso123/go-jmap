package contactcard

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterConditionJSONNameGiven(t *testing.T) {
	filter := &FilterCondition{
		NameGiven: "Ada",
	}

	data, err := jsonv2.Marshal(filter)
	require.NoError(t, err)
	assert.Equal(t, `{"name/given":"Ada"}`, string(data))
}
