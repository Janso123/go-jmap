package emailsubmission_test

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/emailsubmission"
	"github.com/stretchr/testify/require"
)

func TestEmailSubmissionSetOnSuccessWire(t *testing.T) {
	m := &emailsubmission.Set{
		Set: jmap.Set[emailsubmission.EmailSubmission]{
			Account: "a1",
			Destroy: []jmap.ID{"s1"},
		},
		OnSuccessDestroyEmail: []jmap.ID{"e1"},
	}
	require.Equal(t, "EmailSubmission/set", m.Name())
	b, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"accountId":"a1",
		"destroy":["s1"],
		"onSuccessDestroyEmail":["e1"]
	}`, string(b))
}

func TestEmailSubmissionGetEmbedsKit(t *testing.T) {
	m := &emailsubmission.Get{Get: jmap.Get[emailsubmission.EmailSubmission]{Account: "a1"}}
	require.Equal(t, "EmailSubmission/get", m.Name())
	b, err := jsonv2.Marshal(m)
	require.NoError(t, err)
	require.JSONEq(t, `{"accountId":"a1"}`, string(b))

	var resp emailsubmission.GetResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a1","state":"s0","list":[{"id":"s1","emailId":"e1"}],"notFound":[]}`), &resp))
	require.Equal(t, []emailsubmission.EmailSubmission{{ID: "s1", EmailID: "e1"}}, resp.List)
}

func TestEmailSubmissionQueryResponseTotalUint64(t *testing.T) {
	var resp emailsubmission.QueryResponse
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"accountId":"a1","total":3}`), &resp))
	require.Equal(t, uint64(3), resp.Total)
}
