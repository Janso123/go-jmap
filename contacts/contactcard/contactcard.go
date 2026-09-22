package contactcard

import (
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts"
	"github.com/Janso123/go-jmap/contacts/jscontact"
)

func init() {
	jmap.RegisterObject[ContactCard](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodQuery |
			jmap.MethodQueryChanges |
			jmap.MethodSet |
			jmap.MethodCopy,
	)
}

// ContactCard is a JSContact card with JMAP ContactCard metadata.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3
type ContactCard struct {
	ID jmap.ID `json:"id,omitzero"`

	AddressBookIDs map[jmap.ID]bool `json:"addressBookIds,omitzero"`

	jscontact.Card
}

func (ContactCard) JMAPType() string { return "ContactCard" }

func (ContactCard) JMAPCreatable() {}

func (ContactCard) Requires() []jmap.URI { return []jmap.URI{contacts.URI} }

// contactCardJSON is ContactCard without custom marshalers so json/v2 can
// apply Card's Extra embed together with the JMAP metadata keys.
type contactCardJSON struct {
	ID             jmap.ID          `json:"id,omitzero"`
	AddressBookIDs map[jmap.ID]bool `json:"addressBookIds,omitzero"`
	jscontact.Card
}

func (c ContactCard) MarshalJSON() ([]byte, error) {
	return jsonv2.Marshal(contactCardJSON{
		ID:             c.ID,
		AddressBookIDs: c.AddressBookIDs,
		Card:           c.Card,
	})
}

func (c *ContactCard) UnmarshalJSON(data []byte) error {
	var aux contactCardJSON
	if err := jsonv2.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.ID = aux.ID
	c.AddressBookIDs = aux.AddressBookIDs
	c.Card = aux.Card
	return nil
}
