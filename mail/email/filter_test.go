package email

import (
	jsonv2 "encoding/json/v2"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	filter := &FilterCondition{}
	data, err := jsonv2.Marshal(filter)
	assert.NoError(t, err)
	assert.Equal(t, "{}", string(data))
}

func TestFilterWithAttachmentMarshal(t *testing.T) {
	data, err := jsonv2.Marshal(WithAttachment())
	assert.NoError(t, err)
	assert.Equal(t, `{"hasAttachment":true}`, string(data))
}

func TestFilterConstructorsWire(t *testing.T) {
	t.Parallel()

	loc := time.FixedZone("CEST", 2*3600)
	before := jmap.UTCDate(time.Date(2026, 9, 21, 12, 0, 0, 0, loc))

	tests := []struct {
		name string
		f    jmap.Filter
		want string
	}{
		{"hasKeyword", HasKeyword(KeywordSeen), `{"hasKeyword":"$seen"}`},
		{"notKeyword", NotKeyword(KeywordDraft), `{"notKeyword":"$draft"}`},
		{"inMailbox", InMailbox(jmap.ID("mb1")), `{"inMailbox":"mb1"}`},
		{"text", Text("invoice"), `{"text":"invoice"}`},
		{"beforeZ", Before(before), `{"before":"2026-09-21T10:00:00Z"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := jsonv2.Marshal(tt.f)
			require.NoError(t, err)
			require.JSONEq(t, tt.want, string(b))
		})
	}
}

func TestFilterAndOperatorWire(t *testing.T) {
	t.Parallel()
	f := And(InMailbox("mb1"), HasKeyword(KeywordFlagged))
	b, err := jsonv2.Marshal(f)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"operator":"AND",
		"conditions":[
			{"inMailbox":"mb1"},
			{"hasKeyword":"$flagged"}
		]
	}`, string(b))
}

func TestEmailSMIMEFilterMarshal(t *testing.T) {
	t.Parallel()
	f := &FilterCondition{
		HasSMIME:                   new(true),
		HasVerifiedSMIME:           new(true),
		HasVerifiedSMIMEAtDelivery: new(true),
	}
	b, err := jsonv2.Marshal(f)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"hasSmime":true,
		"hasVerifiedSmime":true,
		"hasVerifiedSmimeAtDelivery":true
	}`, string(b))
}

func TestFilterHasAttachmentFalseOnWire(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(WithoutAttachment())
	require.NoError(t, err)
	require.JSONEq(t, `{"hasAttachment":false}`, string(b))
	b, err = jsonv2.Marshal(&FilterCondition{HasSMIME: new(false)})
	require.NoError(t, err)
	require.JSONEq(t, `{"hasSmime":false}`, string(b))
}

func TestFilterMaxSizeZeroOnWire(t *testing.T) {
	t.Parallel()
	b, err := jsonv2.Marshal(&FilterCondition{MaxSize: jmap.Uint64Ptr(0)})
	require.NoError(t, err)
	require.JSONEq(t, `{"maxSize":0}`, string(b))
	b, err = jsonv2.Marshal(&FilterCondition{MinSize: jmap.Uint64Ptr(0)})
	require.NoError(t, err)
	require.JSONEq(t, `{"minSize":0}`, string(b))
	b, err = jsonv2.Marshal(&FilterCondition{})
	require.NoError(t, err)
	require.JSONEq(t, `{}`, string(b))
}
