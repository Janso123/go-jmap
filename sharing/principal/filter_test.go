package principal

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterConditionMarshal(t *testing.T) {
	data, err := jsonv2.Marshal(&FilterCondition{
		AccountIDs:      []jmap.ID{"a1", "a2"},
		Email:           "jane@example.com",
		Name:            "Jane Doe",
		Text:            "jane",
		Type:            TypeIndividual,
		CalendarAddress: "mailto:jane@example.com",
	})
	require.NoError(t, err)
	assert.Equal(t,
		`{"accountIds":["a1","a2"],"email":"jane@example.com","name":"Jane Doe","text":"jane","type":"individual","calendarAddress":"mailto:jane@example.com"}`,
		string(data))
}

func TestFilterOperatorMarshal(t *testing.T) {
	data, err := jsonv2.Marshal(And(
		&FilterCondition{Name: "Jane Doe"},
		&FilterCondition{Type: TypeGroup},
	))
	require.NoError(t, err)
	assert.Equal(t,
		`{"operator":"AND","conditions":[{"name":"Jane Doe"},{"type":"group"}]}`,
		string(data))
}
