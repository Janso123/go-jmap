package jmap_test

import (
	"testing"
	"time"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

func TestUTCDateRejectsOffsetAndLowercase(t *testing.T) {
	var d jmap.UTCDate
	require.Error(t, jsonv2.Unmarshal([]byte(`"2026-09-21T00:00:00+02:00"`), &d))
	require.Error(t, jsonv2.Unmarshal([]byte(`"2026-09-21t00:00:00z"`), &d))
	require.NoError(t, jsonv2.Unmarshal([]byte(`"2026-09-21T00:00:00Z"`), &d))
}

func TestUTCDateAlwaysZ(t *testing.T) {
	loc := time.FixedZone("CEST", 2*3600)
	d := jmap.UTCDate(time.Date(2026, 9, 21, 1, 0, 0, 0, loc))
	b, err := jsonv2.Marshal(d)
	require.NoError(t, err)
	require.Equal(t, `"2026-09-20T23:00:00Z"`, string(b))
}

func TestUTCDateUnmarshal(t *testing.T) {
	var d jmap.UTCDate
	err := jsonv2.Unmarshal([]byte(`"2026-09-20T23:00:00Z"`), &d)
	require.NoError(t, err)
	require.True(t, time.Time(d).Equal(time.Date(2026, 9, 20, 23, 0, 0, 0, time.UTC)))
}

func TestUTCDateKeepsNonZeroSecfrac(t *testing.T) {
	d := jmap.UTCDate(time.Date(2026, 9, 21, 1, 0, 0, 500000000, time.UTC))
	b, err := jsonv2.Marshal(d)
	require.NoError(t, err)
	require.Equal(t, `"2026-09-21T01:00:00.5Z"`, string(b))
}

func TestDateRFC3339Offset(t *testing.T) {
	loc := time.FixedZone("CST", 8*3600)
	d := jmap.Date(time.Date(2014, 10, 30, 14, 12, 0, 0, loc))
	b, err := jsonv2.Marshal(d)
	require.NoError(t, err)
	require.Equal(t, `"2014-10-30T14:12:00+08:00"`, string(b))
}

func TestDateUnmarshalRFCExample(t *testing.T) {
	var d jmap.Date
	err := jsonv2.Unmarshal([]byte(`"2014-10-30T14:12:00+08:00"`), &d)
	require.NoError(t, err)
	require.True(t, time.Time(d).Equal(time.Date(2014, 10, 30, 14, 12, 0, 0, time.FixedZone("", 8*3600))))
}
