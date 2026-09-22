//go:build e2e

package e2e

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts"
	"github.com/Janso123/go-jmap/contacts/addressbook"
	"github.com/Janso123/go-jmap/contacts/contactcard"
	"github.com/Janso123/go-jmap/contacts/jscontact"
)

func TestContacts(t *testing.T) {
	sc := &scenario{name: "Contacts"}
	if alice == nil || alice.Client == nil || bob == nil || bob.Client == nil {
		t.Fatal("account session missing")
	}
	if !hasCap(alice.Client, contacts.URI) {
		skipRest(t, sc, contactSteps, "missing "+string(contacts.URI))
		return
	}

	aliceID, err := alice.Client.PrimaryAccount(contacts.URI)
	if err != nil {
		skipRest(t, sc, contactSteps, err.Error())
		return
	}
	bobID, err := bob.Client.PrimaryAccount(contacts.URI)
	if err != nil {
		skipRest(t, sc, contactSteps, err.Error())
		return
	}

	call[*addressbook.GetResponse](t, sc, step{
		RFC: "RFC 9610", Method: "AddressBook/get", Account: alice.Name, Request: "list",
	}, alice.Client, []jmap.URI{contacts.URI}, false, &addressbook.Get{Account: aliceID}, func(resp *addressbook.GetResponse) (string, error) {
		if resp == nil || len(resp.List) < 1 {
			return "", errString("no address book")
		}
		return "books=" + strconv.Itoa(len(resp.List)), nil
	})

	var bookID jmap.ID
	call[*addressbook.SetResponse](t, sc, step{
		RFC: "RFC 9610", Method: "AddressBook/set", Account: alice.Name, Request: "create name=E2E",
	}, alice.Client, []jmap.URI{contacts.URI}, false, &addressbook.Set{
		Account: aliceID,
		Create: jmap.Some(map[jmap.ID]*addressbook.AddressBook{
			"b1": {Name: "E2E"},
		}),
	}, func(resp *addressbook.SetResponse) (string, error) {
		id, err := createdObjectID(resp, "b1")
		if err != nil {
			return "", err
		}
		bookID = id
		return "id=" + string(bookID), nil
	})

	var stateBefore string
	call[*contactcard.GetResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/get", Account: alice.Name, Request: "ids omitted",
	}, alice.Client, []jmap.URI{contacts.URI}, false, &contactcard.Get{Account: aliceID}, func(resp *contactcard.GetResponse) (string, error) {
		if resp == nil || resp.State == "" {
			return "", errString("card state missing")
		}
		stateBefore = resp.State
		return "state=" + stateBefore, nil
	})

	cardIDs := make([]jmap.ID, 30)
	batch := maxInSet(alice.Client)
	for start := 1; start <= 30; start += batch {
		end := start + batch - 1
		if end > 30 {
			end = 30
		}
		create := map[jmap.ID]*contactcard.ContactCard{}
		var keys []jmap.ID
		for n := start; n <= end; n++ {
			key := jmap.ID(fmt.Sprintf("c%02d", n))
			keys = append(keys, key)
			create[key] = &contactcard.ContactCard{
				AddressBookIDs: map[jmap.ID]bool{bookID: true},
				Card: jscontact.Card{
					Type: "Card",
					UID:  fmt.Sprintf("e2e-card-%02d", n),
					Name: &jscontact.Name{Full: fmt.Sprintf("E2E Person %02d", n)},
				},
			}
		}
		call[*contactcard.SetResponse](t, sc, step{
			RFC: "RFC 9610", Method: "ContactCard/set", Account: alice.Name,
			Request: fmt.Sprintf("create cards %02d-%02d", start, end),
		}, alice.Client, []jmap.URI{contacts.URI}, false, &contactcard.Set{
			Account: aliceID,
			Create:  jmap.Some(create),
		}, func(resp *contactcard.SetResponse) (string, error) {
			made, err := createdCardIDs(resp, keys)
			if err != nil {
				return "", err
			}
			for n := start; n <= end; n++ {
				cardIDs[n-1] = made[jmap.ID(fmt.Sprintf("c%02d", n))]
			}
			return "created=" + strconv.Itoa(len(made)), nil
		})
	}
	cardID := cardIDs[0]
	copyID := cardIDs[1]

	call[*contactcard.QueryResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/query", Account: alice.Name, Request: "inAddressBook",
	}, alice.Client, []jmap.URI{contacts.URI}, false, &contactcard.Query{
		Account: aliceID,
		Filter:  &contactcard.FilterCondition{InAddressBook: bookID},
	}, func(resp *contactcard.QueryResponse) (string, error) {
		if resp == nil || len(resp.IDs) != 30 || !containsAll(resp.IDs, cardIDs) {
			return "", errString("expected 30 cards in the book")
		}
		return "ids=30", nil
	})

	call[*contactcard.GetResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/get", Account: alice.Name, Request: "id=" + string(cardID),
	}, alice.Client, []jmap.URI{contacts.URI}, false, &contactcard.Get{
		Account: aliceID,
		IDs:     jmap.Some([]jmap.ID{cardID}),
	}, func(resp *contactcard.GetResponse) (string, error) {
		if resp == nil || len(resp.List) == 0 || resp.List[0].Name == nil || resp.List[0].Name.Full != "E2E Person 01" {
			return "", errString("full name mismatch")
		}
		return "full=" + resp.List[0].Name.Full, nil
	})

	call[*contactcard.ChangesResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/changes", Account: alice.Name, Request: "sinceState=" + stateBefore,
	}, alice.Client, []jmap.URI{contacts.URI}, false, &contactcard.Changes{
		Account:    aliceID,
		SinceState: stateBefore,
	}, func(resp *contactcard.ChangesResponse) (string, error) {
		if resp == nil || !containsAll(resp.Created, cardIDs) {
			more := false
			if resp != nil {
				more = resp.HasMoreChanges
			}
			return "", errString(fmt.Sprintf("created missing cards hasMore=%v", more))
		}
		return "created=30", nil
	})

	call[*contactcard.SetResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/set", Account: alice.Name, Request: "update name/full",
	}, alice.Client, []jmap.URI{contacts.URI}, false, &contactcard.Set{
		Account: aliceID,
		Update:  jmap.Some(map[jmap.ID]jmap.Patch{cardID: {"name/full": "E2E Renamed"}}),
	}, func(resp *contactcard.SetResponse) (string, error) {
		if err := rejectSetErrors(resp, cardID); err != nil {
			return "", err
		}
		return "updated=" + string(cardID), nil
	})

	call[*contactcard.GetResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/get", Account: alice.Name, Request: "id=" + string(cardID),
	}, alice.Client, []jmap.URI{contacts.URI}, false, &contactcard.Get{
		Account: aliceID,
		IDs:     jmap.Some([]jmap.ID{cardID}),
	}, func(resp *contactcard.GetResponse) (string, error) {
		if resp == nil || len(resp.List) == 0 || resp.List[0].Name == nil || resp.List[0].Name.Full != "E2E Renamed" {
			return "", errString("full name mismatch")
		}
		return "full=" + resp.List[0].Name.Full, nil
	})

	call[*contactcard.SetResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/set", Account: alice.Name, Request: "destroy id=" + string(cardID),
	}, alice.Client, []jmap.URI{contacts.URI}, false, &contactcard.Set{
		Account: aliceID,
		Destroy: jmap.Some([]jmap.ID{cardID}),
	}, func(resp *contactcard.SetResponse) (string, error) {
		if err := rejectSetErrors(resp, cardID); err != nil {
			return "", err
		}
		return "destroyed=" + string(cardID), nil
	})

	call[*contactcard.GetResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/get", Account: alice.Name, Request: "id=" + string(cardID),
	}, alice.Client, []jmap.URI{contacts.URI}, false, &contactcard.Get{
		Account: aliceID,
		IDs:     jmap.Some([]jmap.ID{cardID}),
	}, func(resp *contactcard.GetResponse) (string, error) {
		if resp == nil || !containsID(resp.NotFound, cardID) {
			return "", errString("card id not in notFound")
		}
		return "notFound=" + string(cardID), nil
	})

	call[*contactcard.QueryResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/query", Account: alice.Name, Request: "inAddressBook",
	}, alice.Client, []jmap.URI{contacts.URI}, false, &contactcard.Query{
		Account: aliceID,
		Filter:  &contactcard.FilterCondition{InAddressBook: bookID},
	}, func(resp *contactcard.QueryResponse) (string, error) {
		if resp == nil || len(resp.IDs) != 29 || containsID(resp.IDs, cardID) {
			return "", errString("expected 29 cards in the book")
		}
		return "ids=29", nil
	})

	var bobBook jmap.ID
	call[*addressbook.GetResponse](t, sc, step{
		RFC: "RFC 9610", Method: "AddressBook/get", Account: bob.Name, Request: "list",
	}, bob.Client, []jmap.URI{contacts.URI}, false, &addressbook.Get{Account: bobID}, func(resp *addressbook.GetResponse) (string, error) {
		if resp == nil {
			return "", errString("empty address book list")
		}
		if len(resp.List) == 0 {
			return "list=0", nil
		}
		bobBook = resp.List[0].ID
		return "books=" + strconv.Itoa(len(resp.List)), nil
	})
	if bobBook == "" && !sc.failed {
		call[*addressbook.SetResponse](t, sc, step{
			RFC: "RFC 9610", Method: "AddressBook/set", Account: bob.Name, Request: "create name=Bob",
		}, bob.Client, []jmap.URI{contacts.URI}, false, &addressbook.Set{
			Account: bobID,
			Create: jmap.Some(map[jmap.ID]*addressbook.AddressBook{
				"b1": {Name: "Bob"},
			}),
		}, func(resp *addressbook.SetResponse) (string, error) {
			id, err := createdObjectID(resp, "b1")
			if err != nil {
				return "", err
			}
			bobBook = id
			return "id=" + string(bobBook), nil
		})
	}

	remaining := cardIDs[1:]
	call[*contactcard.QueryResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/query", Account: bob.Name, Request: "Bob does not see Alice's cards",
	}, bob.Client, []jmap.URI{contacts.URI}, false, &contactcard.Query{
		Account: bobID,
		Filter:  &contactcard.FilterCondition{InAddressBook: bobBook},
	}, func(resp *contactcard.QueryResponse) (string, error) {
		if resp == nil || sharesID(resp.IDs, remaining) {
			return "", errString("alice card visible to bob")
		}
		return "ids=" + strconv.Itoa(len(resp.IDs)), nil
	})

	// Bob can copy only after he can read Alice's ContactCard collection.
	// Stalwart checks both account ids and rejects the source account until
	// the book is shared.
	call[*addressbook.SetResponse](t, sc, step{
		RFC: "RFC 9610", Method: "AddressBook/set", Account: alice.Name, Request: "shareWith bob",
	}, alice.Client, []jmap.URI{contacts.URI}, false, &addressbook.Set{
		Account: aliceID,
		Update: jmap.Some(map[jmap.ID]jmap.Patch{
			bookID: {
				"shareWith": map[string]any{
					string(bobID): map[string]bool{
						"mayRead":   true,
						"mayWrite":  true,
						"mayShare":  false,
						"mayDelete": false,
					},
				},
			},
		}),
	}, func(resp *addressbook.SetResponse) (string, error) {
		if err := rejectAddressBookSet(resp, bookID); err != nil {
			return "", err
		}
		return "shared=" + string(bobID), nil
	})

	var copiedID jmap.ID
	_, copied := call[*contactcard.CopyResponse](t, sc, step{
		RFC: "RFC 9610", Method: "ContactCard/copy", Account: bob.Name, Request: "copy id=" + string(copyID),
	}, bob.Client, []jmap.URI{contacts.URI}, true, &contactcard.Copy{Copy: jmap.Copy[contactcard.ContactCard]{
		FromAccount: aliceID,
		Account:     bobID,
		Create: jmap.Some(map[jmap.ID]*contactcard.ContactCard{
			copyID: {AddressBookIDs: map[jmap.ID]bool{bobBook: true}},
		}),
	}}, func(resp *contactcard.CopyResponse) (string, error) {
		if resp == nil {
			return "", errString("empty copy response")
		}
		if nc, ok := resp.NotCreated.Value(); ok && len(nc) > 0 {
			return "", errString(setErrorSummary(nc))
		}
		if len(resp.Created) == 0 {
			return "", errString("copy created nothing")
		}
		created, ok := resp.Created[copyID]
		if !ok || created.ID == "" {
			return "", errString("copied card was not created")
		}
		copiedID = created.ID
		return "id=" + string(copiedID), nil
	})
	if copied || sc.failed {
		call[*contactcard.GetResponse](t, sc, step{
			RFC: "RFC 9610", Method: "ContactCard/get", Account: bob.Name, Request: "id=" + string(copiedID),
		}, bob.Client, []jmap.URI{contacts.URI}, false, &contactcard.Get{
			Account: bobID,
			IDs:     jmap.Some([]jmap.ID{copiedID}),
		}, func(resp *contactcard.GetResponse) (string, error) {
			if resp == nil || len(resp.List) == 0 || resp.List[0].ID != copiedID {
				return "", errString("copied card missing")
			}
			if resp.List[0].Name == nil || resp.List[0].Name.Full != "E2E Person 02" {
				return "", errString("full name mismatch")
			}
			return "full=" + resp.List[0].Name.Full, nil
		})
	}

	if sc.failed {
		t.Fail()
	}
}

var contactSteps = []step{
	{RFC: "RFC 9610", Method: "AddressBook/get", Account: "alice", Request: "list"},
	{RFC: "RFC 9610", Method: "AddressBook/set", Account: "alice", Request: "create name=E2E"},
	{RFC: "RFC 9610", Method: "ContactCard/set", Account: "alice", Request: "create 30 cards"},
	{RFC: "RFC 9610", Method: "ContactCard/query", Account: "alice", Request: "inAddressBook"},
	{RFC: "RFC 9610", Method: "ContactCard/get", Account: "alice", Request: "full name"},
	{RFC: "RFC 9610", Method: "ContactCard/changes", Account: "alice", Request: "sinceState"},
	{RFC: "RFC 9610", Method: "ContactCard/set", Account: "alice", Request: "update name/full"},
	{RFC: "RFC 9610", Method: "ContactCard/set", Account: "alice", Request: "destroy"},
	{RFC: "RFC 9610", Method: "ContactCard/copy", Account: "bob", Request: "copy"},
}

func createdObjectID(resp *addressbook.SetResponse, key jmap.ID) (jmap.ID, error) {
	if resp == nil {
		return "", errString("empty set response")
	}
	if nc, ok := resp.NotCreated.Value(); ok {
		if se := nc[key]; se != nil {
			return "", errString(setErrorText(se))
		}
	}
	created, ok := resp.Created[key]
	if !ok || created.ID == "" {
		return "", errString("address book was not created")
	}
	return created.ID, nil
}

func createdCardIDs(resp *contactcard.SetResponse, keys []jmap.ID) (map[jmap.ID]jmap.ID, error) {
	if resp == nil {
		return nil, errString("empty set response")
	}
	if nc, ok := resp.NotCreated.Value(); ok && len(nc) > 0 {
		return nil, errString(setErrorSummary(nc))
	}
	out := make(map[jmap.ID]jmap.ID, len(keys))
	for _, key := range keys {
		created, ok := resp.Created[key]
		if !ok || created.ID == "" {
			return nil, errString("card " + string(key) + " was not created")
		}
		out[key] = created.ID
	}
	return out, nil
}

func rejectAddressBookSet(resp *addressbook.SetResponse, id jmap.ID) error {
	if resp == nil {
		return errString("empty set response")
	}
	if nu, ok := resp.NotUpdated.Value(); ok {
		if se := nu[id]; se != nil {
			return errString(setErrorText(se))
		}
	}
	return nil
}

func rejectSetErrors(resp *contactcard.SetResponse, id jmap.ID) error {
	if resp == nil {
		return errString("empty set response")
	}
	if nu, ok := resp.NotUpdated.Value(); ok {
		if se := nu[id]; se != nil {
			return errString(setErrorText(se))
		}
	}
	if nd, ok := resp.NotDestroyed.Value(); ok {
		if se := nd[id]; se != nil {
			return errString(setErrorText(se))
		}
	}
	return nil
}

func setErrorSummary(errs map[jmap.ID]*jmap.SetError) string {
	parts := make([]string, 0, len(errs))
	for id, se := range errs {
		parts = append(parts, string(id)+": "+setErrorText(se))
	}
	return joinComma(parts)
}

func setErrorText(se *jmap.SetError) string {
	if se == nil {
		return "set error"
	}
	msg := se.Error()
	if props, ok := se.Properties.Value(); ok && len(props) > 0 {
		msg += " (" + joinComma(props) + ")"
	}
	return msg
}

func containsAll(have, want []jmap.ID) bool {
	for _, id := range want {
		if id == "" || !containsID(have, id) {
			return false
		}
	}
	return true
}

func sharesID(have, want []jmap.ID) bool {
	for _, id := range want {
		if id != "" && containsID(have, id) {
			return true
		}
	}
	return false
}
