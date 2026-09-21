package contactcard

import (
	"encoding/json"

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
	// ContactCard/changes includes updatedProperties; override kit factory.
	jmap.RegisterMethod("ContactCard/changes", func() jmap.MethodResponse { return &ChangesResponse{} })
}

// ContactCard is a JSContact card with JMAP ContactCard metadata.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3
type ContactCard struct {
	ID jmap.ID `json:"id,omitempty"`

	AddressBookIDs map[jmap.ID]bool `json:"addressBookIds,omitempty"`

	jscontact.Card
}

func (ContactCard) JMAPType() string { return "ContactCard" }

func (ContactCard) Requires() []jmap.URI { return []jmap.URI{contacts.URI} }

func (c ContactCard) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(c.Card)
	if err != nil {
		return nil, err
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	if c.ID != "" {
		data, err := json.Marshal(c.ID)
		if err != nil {
			return nil, err
		}
		obj["id"] = data
	}
	if len(c.AddressBookIDs) > 0 {
		data, err := json.Marshal(c.AddressBookIDs)
		if err != nil {
			return nil, err
		}
		obj["addressBookIds"] = data
	}

	return json.Marshal(obj)
}

func (c *ContactCard) UnmarshalJSON(data []byte) error {
	type metadata struct {
		ID             jmap.ID          `json:"id,omitempty"`
		AddressBookIDs map[jmap.ID]bool `json:"addressBookIds,omitempty"`
	}

	var meta metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}

	var card jscontact.Card
	if err := json.Unmarshal(data, &card); err != nil {
		return err
	}

	c.ID = meta.ID
	c.AddressBookIDs = meta.AddressBookIDs
	c.Card = card

	return nil
}
