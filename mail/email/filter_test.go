package email

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	filter := &FilterCondition{}
	data, err := json.Marshal(filter)
	assert.NoError(t, err)
	assert.Equal(t, "{}", string(data))
}

func TestFilterWithAttachmentMarshal(t *testing.T) {
	data, err := json.Marshal(WithAttachment())
	assert.NoError(t, err)
	assert.Equal(t, `{"hasAttachment":true}`, string(data))
}
