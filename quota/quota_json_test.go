package quota_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap/quota"
	"github.com/stretchr/testify/require"
)

func TestQuotaJSONRoundTrip(t *testing.T) {
	t.Parallel()
	raw := `{
		"id":"q1",
		"resourceType":"count",
		"used":10,
		"hardLimit":100,
		"scope":"account",
		"name":"mail",
		"types":["Email"],
		"warnLimit":80,
		"softLimit":90,
		"description":"mailbox count"
	}`
	var q quota.Quota
	require.NoError(t, jsonv2.Unmarshal([]byte(raw), &q))
	require.Equal(t, "q1", string(q.ID))
	require.Equal(t, uint64(10), q.Used)
	require.Equal(t, uint64(100), q.HardLimit)
	require.Equal(t, uint64(80), *q.WarnLimit)
	require.Equal(t, uint64(90), *q.SoftLimit)

	out, err := jsonv2.Marshal(&q)
	require.NoError(t, err)
	require.JSONEq(t, raw, string(out))
}

func TestQuotaNullLimits(t *testing.T) {
	t.Parallel()
	var q quota.Quota
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"q","warnLimit":null,"softLimit":null}`), &q))
	require.Nil(t, q.WarnLimit)
	require.Nil(t, q.SoftLimit)
}
