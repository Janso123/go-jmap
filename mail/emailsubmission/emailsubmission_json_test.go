package emailsubmission_test

import (
	jsonv2 "encoding/json/v2"
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
	require.Equal(t, []jmap.URI{emailsubmission.URI}, (&emailsubmission.Get{}).Requires())
	require.Equal(t, []jmap.URI{emailsubmission.URI}, (&emailsubmission.Set{}).Requires())
	require.Equal(t, []jmap.URI{emailsubmission.URI, mail.URI}, (&emailsubmission.Set{
		OnSuccessDestroyEmail: jmap.Some([]jmap.ID{"e1"}),
	}).Requires())
}

func TestEmailSubmissionJSONRoundTrip(t *testing.T) {
	sendAt := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	sub := emailsubmission.EmailSubmission{
		ID:         "sub1",
		IdentityID: "id1",
		EmailID:    "em1",
		ThreadID:   "th1",
		Envelope: jmap.Some(emailsubmission.Envelope{
			MailFrom: &emailsubmission.Address{Email: "from@example.com"},
			RcptTo: []*emailsubmission.Address{
				{Email: "to@example.com"},
			},
		}),
		SendAt:     jmap.Some(jmap.UTCDate(sendAt)),
		UndoStatus: emailsubmission.UndoPending,
		DeliveryStatus: jmap.Some(map[string]*emailsubmission.DeliveryStatus{
			"to@example.com": {
				SMTPReply: "250 2.1.5 Ok",
				Delivered: emailsubmission.DeliveredYes,
				Displayed: emailsubmission.DisplayedUnknown,
			},
		}),
		DSNBlobIDs: []jmap.ID{"blob1"},
		MDNBlobIDs: []jmap.ID{"blob2"},
	}

	data, err := jsonv2.Marshal(&sub)
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
	require.NoError(t, jsonv2.Unmarshal(data, &got))
	require.Equal(t, sub.ID, got.ID)
	require.Equal(t, sub.IdentityID, got.IdentityID)
	require.Equal(t, sub.EmailID, got.EmailID)
	require.Equal(t, sub.UndoStatus, got.UndoStatus)
	env, ok := sub.Envelope.Value()
	require.True(t, ok)
	gotEnv, ok := got.Envelope.Value()
	require.True(t, ok)
	require.Equal(t, env.MailFrom.Email, gotEnv.MailFrom.Email)
	status, ok := sub.DeliveryStatus.Value()
	require.True(t, ok)
	gotStatus, ok := got.DeliveryStatus.Value()
	require.True(t, ok)
	require.Equal(t, status["to@example.com"].Delivered, gotStatus["to@example.com"].Delivered)
	gotSendAt, ok := got.SendAt.Value()
	require.True(t, ok)
	require.Equal(t, sendAt.UTC(), time.Time(gotSendAt).UTC())
}

func TestSendAtAlwaysZDoesNotMutate(t *testing.T) {
	loc := time.FixedZone("CEST", 2*3600)
	tm := time.Date(2026, 9, 21, 1, 0, 0, 0, loc)
	s := &emailsubmission.EmailSubmission{SendAt: jmap.Some(jmap.UTCDate(tm))}
	orig, ok := s.SendAt.Value()
	require.True(t, ok)
	b, err := jsonv2.Marshal(s)
	require.NoError(t, err)
	require.Contains(t, string(b), `"sendAt":"2026-09-20T23:00:00Z"`)
	got, ok := s.SendAt.Value()
	require.True(t, ok)
	require.Equal(t, orig, got)
}

func TestEmailSubmissionSendAtUTC(t *testing.T) {
	loc := time.FixedZone("CEST", 2*3600)
	sendAt := time.Date(2026, 9, 21, 14, 0, 0, 0, loc)
	sub := emailsubmission.EmailSubmission{SendAt: jmap.Some(jmap.UTCDate(sendAt))}

	data, err := jsonv2.Marshal(&sub)
	require.NoError(t, err)
	require.Contains(t, string(data), "Z")
	require.NotContains(t, string(data), "+02:00")
	require.Contains(t, string(data), "2026-09-21T12:00:00Z")
}

func TestEmailSubmissionGetResponseRegistered(t *testing.T) {
	raw := []byte(`{"sessionState":"s1","methodResponses":[["EmailSubmission/get",{"accountId":"u1","state":"s0","list":[],"notFound":[]},"0"]]}`)
	var resp jmap.Response
	require.NoError(t, jsonv2.Unmarshal(raw, &resp))
	require.Len(t, resp.Responses, 1)
	_, ok := resp.Responses[0].Args.(*emailsubmission.GetResponse)
	require.True(t, ok)
}

func TestEmptyMailFromOnWire(t *testing.T) {
	t.Parallel()
	env := emailsubmission.Envelope{
		MailFrom: &emailsubmission.Address{Email: ""},
		RcptTo:   []*emailsubmission.Address{{Email: "a@b.c"}},
	}
	b, err := jsonv2.Marshal(env)
	require.NoError(t, err)
	require.JSONEq(t, `{"mailFrom":{"email":""},"rcptTo":[{"email":"a@b.c"}]}`, string(b))
}
