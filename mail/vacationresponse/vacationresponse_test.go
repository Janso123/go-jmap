package vacationresponse_test

import (
	jsonv2 "encoding/json/v2"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/vacationresponse"
	"github.com/stretchr/testify/require"
)

func TestVacationFromDateUTC(t *testing.T) {
	loc := time.FixedZone("CEST", 2*3600)
	v := vacationresponse.VacationResponse{
		FromDate: jmap.UTCDatePtr(time.Date(2026, 9, 21, 1, 0, 0, 0, loc)),
	}
	b, err := jsonv2.Marshal(v)
	require.NoError(t, err)
	require.Contains(t, string(b), "Z")
	require.NotContains(t, string(b), "+02:00")
}
