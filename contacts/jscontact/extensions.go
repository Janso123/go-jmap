package jscontact

import (
	"encoding/json"
	"reflect"
	"strings"
)

func marshalObject(v any, extensions map[string]json.RawMessage) ([]byte, error) {
	if len(extensions) == 0 {
		return json.Marshal(v)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	for key, value := range extensions {
		if _, exists := raw[key]; exists {
			continue
		}
		raw[key] = value
	}
	return json.Marshal(raw)
}

func unmarshalObject[T any](data []byte, dst *T) (map[string]json.RawMessage, error) {
	if err := json.Unmarshal(data, dst); err != nil {
		return nil, err
	}

	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	deleteKnownJSONFields(raw, reflect.TypeOf(*dst))
	if len(raw) == 0 {
		return nil, nil
	}
	return raw, nil
}

func deleteKnownJSONFields(raw map[string]json.RawMessage, typ reflect.Type) {
	typ = dereferenceType(typ)
	if typ.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			if name, ok := taggedJSONFieldName(field); ok {
				delete(raw, name)
				continue
			}
			deleteKnownJSONFields(raw, field.Type)
			continue
		}
		if !field.IsExported() {
			continue
		}
		if name, ok := taggedJSONFieldName(field); ok {
			delete(raw, name)
			continue
		}
		delete(raw, field.Name)
	}
}

func taggedJSONFieldName(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", false
	}
	if tag == "" {
		return "", false
	}
	return strings.Split(tag, ",")[0], true
}

func dereferenceType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

func (c Card) MarshalJSON() ([]byte, error) {
	type plain Card
	return marshalObject(plain(c), c.extensions)
}

func (c *Card) UnmarshalJSON(data []byte) error {
	type plain Card
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*c = Card(decoded)
	c.extensions = extensions
	return nil
}

func (r Relation) MarshalJSON() ([]byte, error) {
	type plain Relation
	return marshalObject(plain(r), r.extensions)
}

func (r *Relation) UnmarshalJSON(data []byte) error {
	type plain Relation
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*r = Relation(decoded)
	r.extensions = extensions
	return nil
}

func (n Name) MarshalJSON() ([]byte, error) {
	type plain Name
	return marshalObject(plain(n), n.extensions)
}

func (n *Name) UnmarshalJSON(data []byte) error {
	type plain Name
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*n = Name(decoded)
	n.extensions = extensions
	return nil
}

func (n NameComponent) MarshalJSON() ([]byte, error) {
	type plain NameComponent
	return marshalObject(plain(n), n.extensions)
}

func (n *NameComponent) UnmarshalJSON(data []byte) error {
	type plain NameComponent
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*n = NameComponent(decoded)
	n.extensions = extensions
	return nil
}

func (n Nickname) MarshalJSON() ([]byte, error) {
	type plain Nickname
	return marshalObject(plain(n), n.extensions)
}

func (n *Nickname) UnmarshalJSON(data []byte) error {
	type plain Nickname
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*n = Nickname(decoded)
	n.extensions = extensions
	return nil
}

func (o Organization) MarshalJSON() ([]byte, error) {
	type plain Organization
	return marshalObject(plain(o), o.extensions)
}

func (o *Organization) UnmarshalJSON(data []byte) error {
	type plain Organization
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*o = Organization(decoded)
	o.extensions = extensions
	return nil
}

func (o OrgUnit) MarshalJSON() ([]byte, error) {
	type plain OrgUnit
	return marshalObject(plain(o), o.extensions)
}

func (o *OrgUnit) UnmarshalJSON(data []byte) error {
	type plain OrgUnit
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*o = OrgUnit(decoded)
	o.extensions = extensions
	return nil
}

func (s SpeakToAs) MarshalJSON() ([]byte, error) {
	type plain SpeakToAs
	return marshalObject(plain(s), s.extensions)
}

func (s *SpeakToAs) UnmarshalJSON(data []byte) error {
	type plain SpeakToAs
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*s = SpeakToAs(decoded)
	s.extensions = extensions
	return nil
}

func (p Pronouns) MarshalJSON() ([]byte, error) {
	type plain Pronouns
	return marshalObject(plain(p), p.extensions)
}

func (p *Pronouns) UnmarshalJSON(data []byte) error {
	type plain Pronouns
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*p = Pronouns(decoded)
	p.extensions = extensions
	return nil
}

func (t Title) MarshalJSON() ([]byte, error) {
	type plain Title
	return marshalObject(plain(t), t.extensions)
}

func (t *Title) UnmarshalJSON(data []byte) error {
	type plain Title
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*t = Title(decoded)
	t.extensions = extensions
	return nil
}

func (e EmailAddress) MarshalJSON() ([]byte, error) {
	type plain EmailAddress
	return marshalObject(plain(e), e.extensions)
}

func (e *EmailAddress) UnmarshalJSON(data []byte) error {
	type plain EmailAddress
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*e = EmailAddress(decoded)
	e.extensions = extensions
	return nil
}

func (o OnlineService) MarshalJSON() ([]byte, error) {
	type plain OnlineService
	return marshalObject(plain(o), o.extensions)
}

func (o *OnlineService) UnmarshalJSON(data []byte) error {
	type plain OnlineService
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*o = OnlineService(decoded)
	o.extensions = extensions
	return nil
}

func (p Phone) MarshalJSON() ([]byte, error) {
	type plain Phone
	return marshalObject(plain(p), p.extensions)
}

func (p *Phone) UnmarshalJSON(data []byte) error {
	type plain Phone
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*p = Phone(decoded)
	p.extensions = extensions
	return nil
}

func (l LanguagePref) MarshalJSON() ([]byte, error) {
	type plain LanguagePref
	return marshalObject(plain(l), l.extensions)
}

func (l *LanguagePref) UnmarshalJSON(data []byte) error {
	type plain LanguagePref
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*l = LanguagePref(decoded)
	l.extensions = extensions
	return nil
}

func (r Resource) MarshalJSON() ([]byte, error) {
	type plain Resource
	return marshalObject(plain(r), r.extensions)
}

func (r *Resource) UnmarshalJSON(data []byte) error {
	type plain Resource
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*r = Resource(decoded)
	r.extensions = extensions
	return nil
}

func (c Calendar) MarshalJSON() ([]byte, error) {
	type plain Calendar
	return marshalObject(plain(c), c.extensions)
}

func (c *Calendar) UnmarshalJSON(data []byte) error {
	type plain Calendar
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*c = Calendar(decoded)
	c.extensions = extensions
	return nil
}

func (s SchedulingAddress) MarshalJSON() ([]byte, error) {
	type plain SchedulingAddress
	return marshalObject(plain(s), s.extensions)
}

func (s *SchedulingAddress) UnmarshalJSON(data []byte) error {
	type plain SchedulingAddress
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*s = SchedulingAddress(decoded)
	s.extensions = extensions
	return nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	type plain Address
	return marshalObject(plain(a), a.extensions)
}

func (a *Address) UnmarshalJSON(data []byte) error {
	type plain Address
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*a = Address(decoded)
	a.extensions = extensions
	return nil
}

func (a AddressComponent) MarshalJSON() ([]byte, error) {
	type plain AddressComponent
	return marshalObject(plain(a), a.extensions)
}

func (a *AddressComponent) UnmarshalJSON(data []byte) error {
	type plain AddressComponent
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*a = AddressComponent(decoded)
	a.extensions = extensions
	return nil
}

func (c CryptoKey) MarshalJSON() ([]byte, error) {
	type plain CryptoKey
	return marshalObject(plain(c), c.extensions)
}

func (c *CryptoKey) UnmarshalJSON(data []byte) error {
	type plain CryptoKey
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*c = CryptoKey(decoded)
	c.extensions = extensions
	return nil
}

func (d Directory) MarshalJSON() ([]byte, error) {
	type plain Directory
	return marshalObject(plain(d), d.extensions)
}

func (d *Directory) UnmarshalJSON(data []byte) error {
	type plain Directory
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*d = Directory(decoded)
	d.extensions = extensions
	return nil
}

func (l Link) MarshalJSON() ([]byte, error) {
	type plain Link
	return marshalObject(plain(l), l.extensions)
}

func (l *Link) UnmarshalJSON(data []byte) error {
	type plain Link
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*l = Link(decoded)
	l.extensions = extensions
	return nil
}

func (m Media) MarshalJSON() ([]byte, error) {
	type plain Media
	return marshalObject(plain(m), m.extensions)
}

func (m *Media) UnmarshalJSON(data []byte) error {
	type plain Media
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*m = Media(decoded)
	m.extensions = extensions
	return nil
}

func (a Anniversary) MarshalJSON() ([]byte, error) {
	type plain Anniversary
	return marshalObject(plain(a), a.extensions)
}

func (a *Anniversary) UnmarshalJSON(data []byte) error {
	type plain Anniversary
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*a = Anniversary(decoded)
	a.extensions = extensions
	return nil
}

func (p PartialDate) MarshalJSON() ([]byte, error) {
	type plain PartialDate
	return marshalObject(plain(p), p.extensions)
}

func (p *PartialDate) UnmarshalJSON(data []byte) error {
	type plain PartialDate
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*p = PartialDate(decoded)
	p.extensions = extensions
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	type plain Timestamp
	return marshalObject(plain(t), t.extensions)
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	type plain Timestamp
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*t = Timestamp(decoded)
	t.extensions = extensions
	return nil
}

func (n Note) MarshalJSON() ([]byte, error) {
	type plain Note
	return marshalObject(plain(n), n.extensions)
}

func (n *Note) UnmarshalJSON(data []byte) error {
	type plain Note
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*n = Note(decoded)
	n.extensions = extensions
	return nil
}

func (a Author) MarshalJSON() ([]byte, error) {
	type plain Author
	return marshalObject(plain(a), a.extensions)
}

func (a *Author) UnmarshalJSON(data []byte) error {
	type plain Author
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*a = Author(decoded)
	a.extensions = extensions
	return nil
}

func (p PersonalInfo) MarshalJSON() ([]byte, error) {
	type plain PersonalInfo
	return marshalObject(plain(p), p.extensions)
}

func (p *PersonalInfo) UnmarshalJSON(data []byte) error {
	type plain PersonalInfo
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*p = PersonalInfo(decoded)
	p.extensions = extensions
	return nil
}
