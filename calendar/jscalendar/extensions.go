package jscalendar

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
	if tag == "-" || tag == "" {
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

func (e Event) MarshalJSON() ([]byte, error) {
	type plain Event
	return marshalObject(plain(e), e.extensions)
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type plain Event
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*e = Event(decoded)
	e.extensions = extensions
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

func (l Location) MarshalJSON() ([]byte, error) {
	type plain Location
	return marshalObject(plain(l), l.extensions)
}

func (l *Location) UnmarshalJSON(data []byte) error {
	type plain Location
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*l = Location(decoded)
	l.extensions = extensions
	return nil
}

func (v VirtualLocation) MarshalJSON() ([]byte, error) {
	type plain VirtualLocation
	return marshalObject(plain(v), v.extensions)
}

func (v *VirtualLocation) UnmarshalJSON(data []byte) error {
	type plain VirtualLocation
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*v = VirtualLocation(decoded)
	v.extensions = extensions
	return nil
}

func (n NDay) MarshalJSON() ([]byte, error) {
	type plain NDay
	return marshalObject(plain(n), n.extensions)
}

func (n *NDay) UnmarshalJSON(data []byte) error {
	type plain NDay
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*n = NDay(decoded)
	n.extensions = extensions
	return nil
}

func (r RecurrenceRule) MarshalJSON() ([]byte, error) {
	type plain RecurrenceRule
	return marshalObject(plain(r), r.extensions)
}

func (r *RecurrenceRule) UnmarshalJSON(data []byte) error {
	type plain RecurrenceRule
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*r = RecurrenceRule(decoded)
	r.extensions = extensions
	return nil
}

func (p Participant) MarshalJSON() ([]byte, error) {
	type plain Participant
	return marshalObject(plain(p), p.extensions)
}

func (p *Participant) UnmarshalJSON(data []byte) error {
	type plain Participant
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*p = Participant(decoded)
	p.extensions = extensions
	return nil
}

func (a Alert) MarshalJSON() ([]byte, error) {
	type plain Alert
	return marshalObject(plain(a), a.extensions)
}

func (a *Alert) UnmarshalJSON(data []byte) error {
	type plain Alert
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*a = Alert(decoded)
	a.extensions = extensions
	return nil
}

func (o OffsetTrigger) MarshalJSON() ([]byte, error) {
	type plain OffsetTrigger
	return marshalObject(plain(o), o.extensions)
}

func (o *OffsetTrigger) UnmarshalJSON(data []byte) error {
	type plain OffsetTrigger
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*o = OffsetTrigger(decoded)
	o.extensions = extensions
	return nil
}

func (a AbsoluteTrigger) MarshalJSON() ([]byte, error) {
	type plain AbsoluteTrigger
	return marshalObject(plain(a), a.extensions)
}

func (a *AbsoluteTrigger) UnmarshalJSON(data []byte) error {
	type plain AbsoluteTrigger
	var decoded plain
	extensions, err := unmarshalObject(data, &decoded)
	if err != nil {
		return err
	}
	*a = AbsoluteTrigger(decoded)
	a.extensions = extensions
	return nil
}

func (t Trigger) MarshalJSON() ([]byte, error) {
	switch {
	case len(t.Unknown) > 0:
		return t.Unknown, nil
	case t.AbsoluteTrigger != nil:
		return json.Marshal(t.AbsoluteTrigger)
	case t.OffsetTrigger != nil:
		return json.Marshal(t.OffsetTrigger)
	default:
		return []byte("null"), nil
	}
}

func (t *Trigger) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = Trigger{}
		return nil
	}

	var probe struct {
		Type string `json:"@type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}

	switch probe.Type {
	case "", "OffsetTrigger":
		var offset OffsetTrigger
		if err := json.Unmarshal(data, &offset); err != nil {
			return err
		}
		t.OffsetTrigger = &offset
		t.AbsoluteTrigger = nil
		t.Unknown = nil
	case "AbsoluteTrigger":
		var absolute AbsoluteTrigger
		if err := json.Unmarshal(data, &absolute); err != nil {
			return err
		}
		t.AbsoluteTrigger = &absolute
		t.OffsetTrigger = nil
		t.Unknown = nil
	default:
		t.Unknown = append(t.Unknown[:0], data...)
		t.OffsetTrigger = nil
		t.AbsoluteTrigger = nil
	}

	return nil
}
