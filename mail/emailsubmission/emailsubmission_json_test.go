package emailsubmission_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
	"github.com/Janso123/go-jmap/mail/emailsubmission"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailSubmissionMethodNames(t *testing.T) {
	require.Equal(t, "EmailSubmission/get", (&emailsubmission.Get{}).Name())
	require.Equal(t, "EmailSubmission/changes", (&emailsubmission.Changes{}).Name())
	require.Equal(t, "EmailSubmission/query", (&emailsubmission.Query{}).Name())
	require.Equal(t, "EmailSubmission/queryChanges", (&emailsubmission.QueryChanges{}).Name())
	require.Equal(t, "EmailSubmission/set", (&emailsubmission.Set{}).Name())
}

func TestEmailSubmissionRequires(t *testing.T) {
	want := []jmap.URI{emailsubmission.URI, mail.URI}
	require.Equal(t, want, (&emailsubmission.Get{}).Requires())
	require.Equal(t, want, (&emailsubmission.Set{}).Requires())
}

func TestEmailSubmissionJSONRoundTrip(t *testing.T) {
	sendAt := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	sub := emailsubmission.EmailSubmission{
		ID:         "sub1",
		IdentityID: "id1",
		EmailID:    "em1",
		ThreadID:   "th1",
		Envelope: &emailsubmission.Envelope{
			MailFrom: &emailsubmission.Address{Email: "from@example.com"},
			RcptTo: []*emailsubmission.Address{
				{Email: "to@example.com"},
			},
		},
		SendAt:     &sendAt,
		UndoStatus: emailsubmission.UndoPending,
		DeliveryStatus: map[string]*emailsubmission.DeliveryStatus{
			"to@example.com": {
				SMTPReply: "250 2.1.5 Ok",
				Delivered: emailsubmission.DeliveredYes,
				Displayed: emailsubmission.DisplayedUnknown,
			},
		},
		DSNBlobIDs: []jmap.ID{"blob1"},
		MDNBlobIDs: []jmap.ID{"blob2"},
	}

	data, err := json.Marshal(&sub)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"id": "sub1",
		"identityId": "id1",
		"emailId": "em1",
		"threadId": "th1",
		"envelope": {
			"mailFrom": {"email": "from@example.com"},
			"rcptTo": [{"email": "to@example.com"}]
		},
		"sendAt": "2026-09-21T12:00:00Z",
		"undoStatus": "pending",
		"deliveryStatus": {
			"to@example.com": {
				"smtpReply": "250 2.1.5 Ok",
				"delivered": "yes",
				"displayed": "unknown"
			}
		},
		"dsnBlobIds": ["blob1"],
		"mdnBlobIds": ["blob2"]
	}`, string(data))

	var got emailsubmission.EmailSubmission
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, sub.ID, got.ID)
	require.Equal(t, sub.IdentityID, got.IdentityID)
	require.Equal(t, sub.EmailID, got.EmailID)
	require.Equal(t, sub.UndoStatus, got.UndoStatus)
	require.Equal(t, sub.Envelope.MailFrom.Email, got.Envelope.MailFrom.Email)
	require.Equal(t, sub.DeliveryStatus["to@example.com"].Delivered, got.DeliveryStatus["to@example.com"].Delivered)
	require.Equal(t, sendAt.UTC(), got.SendAt.UTC())
}

func TestEmailSubmissionSendAtUTC(t *testing.T) {
	loc := time.FixedZone("CEST", 2*3600)
	sendAt := time.Date(2026, 9, 21, 14, 0, 0, 0, loc)
	sub := emailsubmission.EmailSubmission{SendAt: &sendAt}

	data, err := json.Marshal(&sub)
	require.NoError(t, err)
	require.Contains(t, string(data), "Z")
	require.NotContains(t, string(data), "+02:00")
	require.Contains(t, string(data), "2026-09-21T12:00:00Z")
}

func TestEmailSubmissionGetResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["EmailSubmission/get",{"accountId":"u1","state":"s0","list":[],"notFound":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*emailsubmission.GetResponse)
	require.True(t, ok)
}
