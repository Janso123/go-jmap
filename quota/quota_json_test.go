package quota_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
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
	require.Equal(t, jmap.UnsignedInt(10), q.Used)
	require.Equal(t, jmap.UnsignedInt(100), q.HardLimit)
	warn, ok := q.WarnLimit.Value()
	require.True(t, ok)
	require.Equal(t, jmap.UnsignedInt(80), warn)
	soft, ok := q.SoftLimit.Value()
	require.True(t, ok)
	require.Equal(t, jmap.UnsignedInt(90), soft)

	out, err := jsonv2.Marshal(&q)
	require.NoError(t, err)
	require.JSONEq(t, raw, string(out))
}

func TestQuotaNullLimits(t *testing.T) {
	t.Parallel()
	var q quota.Quota
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"q","used":0,"warnLimit":null,"softLimit":null}`), &q))
	require.Equal(t, jmap.UnsignedInt(0), q.HardLimit)
	require.True(t, q.WarnLimit.IsNull())
	require.True(t, q.SoftLimit.IsNull())

	b, err := jsonv2.Marshal(&q)
	require.NoError(t, err)
	require.Contains(t, string(b), `"hardLimit":0`)
	require.NotContains(t, string(b), `"hardLimit":null`)
	require.Contains(t, string(b), `"warnLimit":null`)
	require.Contains(t, string(b), `"softLimit":null`)
}

func TestQuotaWarnSoftNullRemarshals(t *testing.T) {
	t.Parallel()
	var q quota.Quota
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"q","used":1,"warnLimit":null,"softLimit":null}`), &q))
	require.True(t, q.WarnLimit.IsNull())
	require.True(t, q.SoftLimit.IsNull())

	b, err := jsonv2.Marshal(&q)
	require.NoError(t, err)
	require.Contains(t, string(b), `"warnLimit":null`)
	require.Contains(t, string(b), `"softLimit":null`)
}

func TestQuotaWarnSoftUnsetOmitted(t *testing.T) {
	t.Parallel()
	q := quota.Quota{ID: "q", Used: 1}
	b, err := jsonv2.Marshal(&q)
	require.NoError(t, err)
	require.NotContains(t, string(b), "warnLimit")
	require.NotContains(t, string(b), "softLimit")
}

func TestQuotaHardLimitExplicitZeroOnWire(t *testing.T) {
	t.Parallel()
	q := quota.Quota{ID: "q1", Used: 0, HardLimit: 0}
	b, err := jsonv2.Marshal(&q)
	require.NoError(t, err)
	require.Contains(t, string(b), `"used":0`)
	require.Contains(t, string(b), `"hardLimit":0`)
	require.NotContains(t, string(b), `"hardLimit":null`)
}

func TestQuotaHardLimitAbsentRemarshalsZero(t *testing.T) {
	t.Parallel()
	var q quota.Quota
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"id":"q","used":5}`), &q))
	require.Equal(t, jmap.UnsignedInt(0), q.HardLimit)

	b, err := jsonv2.Marshal(&q)
	require.NoError(t, err)
	require.Contains(t, string(b), `"hardLimit":0`)
	require.Contains(t, string(b), `"used":5`)
}

func TestQuotaHardLimitUsedAndNullDescription(t *testing.T) {
	t.Parallel()
	var q quota.Quota
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"hardLimit":5,"used":0,"description":null}`), &q))
	b, err := jsonv2.Marshal(&q)
	require.NoError(t, err)
	out := string(b)
	require.Contains(t, out, `"hardLimit":5`)
	require.Contains(t, out, `"used":0`)
	require.Contains(t, out, `"description":null`)
}

func TestQuotaHardLimitZeroRemarshals(t *testing.T) {
	t.Parallel()
	var q quota.Quota
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"hardLimit":0}`), &q))
	b, err := jsonv2.Marshal(&q)
	require.NoError(t, err)
	require.Contains(t, string(b), `"hardLimit":0`)
}

func TestQuotaZeroUsedHardLimitOnWire(t *testing.T) {
	t.Parallel()
	q := quota.Quota{ID: "q1", ResourceType: "count", Scope: "account"}
	b, err := jsonv2.Marshal(&q)
	require.NoError(t, err)
	require.Contains(t, string(b), `"used":0`)
	require.Contains(t, string(b), `"hardLimit":0`)
}
