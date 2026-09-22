package vacationresponse_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap/mail/vacationresponse"
	"github.com/stretchr/testify/require"
)

func roundTrip[T any](t *testing.T, in string) string {
	t.Helper()
	var v T
	require.NoError(t, jsonv2.Unmarshal([]byte(in), &v))
	b, err := jsonv2.Marshal(&v)
	require.NoError(t, err)
	return string(b)
}

func TestVacationNullFromDateAndSubject(t *testing.T) {
	t.Parallel()
	out := roundTrip[vacationresponse.VacationResponse](t, `{"fromDate":null,"subject":null}`)
	require.Contains(t, out, `"fromDate":null`)
	require.Contains(t, out, `"subject":null`)
}
