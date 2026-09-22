package contactcard

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts"
	"github.com/Janso123/go-jmap/contacts/jscontact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Get{
		Account:    "u1",
		IDs:        jmap.Some([]jmap.ID{"c1"}),
		Properties: jmap.Some([]string{"uid", "name"}),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["ContactCard/get",{"accountId":"u1","ids":["c1"],"properties":["uid","name"]},"0"]]}`,
		string(data))
}

func TestGetInvokeWithResultReferences(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Get{
		Account: "u1",
		ReferenceIDs: &jmap.ResultReference{
			ResultOf: "c1",
			Name:     "ContactCard/changes",
			Path:     "/created",
		},
		ReferenceProperties: &jmap.ResultReference{
			ResultOf: "c2",
			Name:     "Core/echo",
			Path:     "/properties",
		},
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["ContactCard/get",{"accountId":"u1","#ids":{"resultOf":"c1","name":"ContactCard/changes","path":"/created"},"#properties":{"resultOf":"c2","name":"Core/echo","path":"/properties"}},"0"]]}`,
		string(data))
}

func TestGetRequiresContactsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{contacts.URI}, (&Get{}).Requires())
}

func TestChangesInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Changes{
		Account:    "u1",
		SinceState: "s1",
		MaxChanges: jmap.Uint64Ptr(50),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["ContactCard/changes",{"accountId":"u1","sinceState":"s1","maxChanges":50},"0"]]}`,
		string(data))
}

func TestChangesRequiresContactsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{contacts.URI}, (&Changes{}).Requires())
}

func TestQueryInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Query{
		Account: "u1",
		Filter: &FilterCondition{
			InAddressBook: "ab1",
			NameGiven:     "Ada",
		},
		Sort: []*jmap.Comparator{
			{Property: "name/given"},
		},
		Limit: jmap.Uint64Ptr(10),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["ContactCard/query",{"accountId":"u1","limit":10,"filter":{"inAddressBook":"ab1","name/given":"Ada"},"sort":[{"property":"name/given"}]},"0"]]}`,
		string(data))
}

func TestQueryRequiresContactsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{contacts.URI}, (&Query{}).Requires())
}

func TestQueryChangesInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&QueryChanges{
		Account:         "u1",
		Filter:          &FilterCondition{Email: "ada@example.com"},
		Sort:            []*jmap.Comparator{{Property: "name/surname"}},
		SinceQueryState: "q1",
		MaxChanges:      jmap.Uint64Ptr(25),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["ContactCard/queryChanges",{"accountId":"u1","sinceQueryState":"q1","maxChanges":25,"filter":{"email":"ada@example.com"},"sort":[{"property":"name/surname"}]},"0"]]}`,
		string(data))
}

func TestQueryChangesRequiresContactsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{contacts.URI}, (&QueryChanges{}).Requires())
}

func TestSetInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Set{
		Account: "u1",
		Create: jmap.Some(map[jmap.ID]*ContactCard{
			"c1": {
				AddressBookIDs: map[jmap.ID]bool{
					"ab1": true,
				},
				Card: jscontact.Card{
					UID: "urn:uuid:ada",
				},
			},
		}),
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["ContactCard/set",{"accountId":"u1","create":{"c1":{"addressBookIds":{"ab1":true},"uid":"urn:uuid:ada"}}},"0"]]}`,
		string(data))
}

func TestSetJSON(t *testing.T) {
	set := &Set{
		Account: "u1",
		Update: jmap.Some(map[jmap.ID]jmap.Patch{
			"c1": {
				"notes/n1": nil,
			},
		}),
	}

	data, err := jsonv2.Marshal(set)
	require.NoError(t, err)
	assert.Equal(t,
		`{"accountId":"u1","update":{"c1":{"notes/n1":null}}}`,
		string(data))
}

func TestSetRequiresContactsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{contacts.URI}, (&Set{}).Requires())
}

func TestCopyInvoke(t *testing.T) {
	req := &jmap.Request{}

	id := req.Invoke(&Copy{
		FromAccount: "u1",
		Account:     "u2",
		Create: jmap.Some(map[jmap.ID]*ContactCard{
			"c1": {
				AddressBookIDs: map[jmap.ID]bool{
					"ab2": true,
				},
			},
		}),
		OnSuccessDestroyOriginal: true,
	})
	assert.Equal(t, "0", id)

	data, err := jsonv2.Marshal(req)
	require.NoError(t, err)
	assert.Equal(t,
		`{"using":["urn:ietf:params:jmap:contacts"],"methodCalls":[["ContactCard/copy",{"fromAccountId":"u1","accountId":"u2","create":{"c1":{"addressBookIds":{"ab2":true}}},"onSuccessDestroyOriginal":true},"0"]]}`,
		string(data))
}

func TestCopyRequiresContactsCapability(t *testing.T) {
	assert.Equal(t, []jmap.URI{contacts.URI}, (&Copy{}).Requires())
}

func TestChangesResponseMatchesKit(t *testing.T) {
	t.Parallel()
	data, err := jsonv2.Marshal(ChangesResponse{NewState: "s"})
	require.NoError(t, err)
	kit, err := jsonv2.Marshal(jmap.ChangesResponse{NewState: "s"})
	require.NoError(t, err)
	assert.JSONEq(t, string(kit), string(data))
}
