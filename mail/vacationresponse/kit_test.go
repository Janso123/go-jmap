package vacationresponse_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/vacationresponse"
	"github.com/stretchr/testify/require"
)

func TestVacationResponseGetEmbedsKit(t *testing.T) {
	m := &vacationresponse.Get{Get: jmap.Get[vacationresponse.VacationResponse]{Account: "a1"}}
	require.Equal(t, "VacationResponse/get", m.Name())
	b, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"a1"}`, string(b))

	var resp vacationresponse.GetResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a1","state":"s0","list":[{"id":"singleton"}],"notFound":[]}`), &resp))
	require.Equal(t, []vacationresponse.VacationResponse{{ID: "singleton"}}, resp.List)
}
