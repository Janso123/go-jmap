package jmap_test

import (
	"testing"
	"time"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

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

func TestDateLocalCalendar(t *testing.T) {
	loc := time.FixedZone("CEST", 2*3600)
	d := jmap.Date(time.Date(2026, 9, 21, 1, 0, 0, 0, loc))
	b, err := jsonv2.Marshal(d)
	require.NoError(t, err)
	require.Equal(t, `"2026-09-21"`, string(b))
}

func TestDateUnmarshal(t *testing.T) {
	var d jmap.Date
	err := jsonv2.Unmarshal([]byte(`"2026-09-21"`), &d)
	require.NoError(t, err)
	require.Equal(t, "2026-09-21", time.Time(d).Format("2006-01-02"))
}
