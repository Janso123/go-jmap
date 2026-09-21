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

func TestVacationIsEnabledFalse(t *testing.T) {
	t.Parallel()
	v := vacationresponse.VacationResponse{IsEnabled: new(false)}
	b, err := jsonv2.Marshal(v)
	require.NoError(t, err)
	require.JSONEq(t, `{"isEnabled":false}`, string(b))
}

func TestVacationResponseMethodsRegistered(t *testing.T) {
	t.Parallel()
	raw := []byte(`{"sessionState":"s","methodResponses":[
		["VacationResponse/get",{"list":[],"notFound":[]},"0"],
		["VacationResponse/set",{"updated":{}},"1"]
	]}`)
	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 2)
	_, ok := resp.Responses[0].Args.(*vacationresponse.GetResponse)
	require.True(t, ok)
	_, ok = resp.Responses[1].Args.(*vacationresponse.SetResponse)
	require.True(t, ok)
}
