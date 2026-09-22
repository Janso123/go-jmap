package email

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
)

// HeaderForm is a JMAP Email header:* form (RFC 8621 §4.1.3).
type HeaderForm string

const (
	FormRaw              HeaderForm = "Raw"
	FormText             HeaderForm = "Text"
	FormAddresses        HeaderForm = "Addresses"
	FormGroupedAddresses HeaderForm = "GroupedAddresses"
	FormMessageIds       HeaderForm = "MessageIds"
	FormDate             HeaderForm = "Date"
	FormURLs             HeaderForm = "URLs"
)

// HeaderProp builds a header:* property name, e.g. "header:List-Unsubscribe:asURLs:all".
func HeaderProp(name string, form HeaderForm, all bool) string {
	p := "header:" + name
	if form != "" && form != FormRaw {
		p += ":as" + string(form)
	}
	if all {
		p += ":all"
	}
	return p
}

func headerLookup(extra map[string]jsontext.Value, name string, form HeaderForm) (jsontext.Value, bool) {
	if len(extra) == 0 {
		return nil, false
	}
	if v, ok := extra[HeaderProp(name, form, false)]; ok {
		return v, true
	}
	if v, ok := extra[HeaderProp(name, form, true)]; ok {
		return v, true
	}
	return nil, false
}

// headerDecode returns an unset Optional when the property is absent, a null
// Optional when the server sent JSON null (no such header), and the decoded
// value otherwise.
func headerDecode[T any](extra map[string]jsontext.Value, name string, form HeaderForm) jmap.Optional[T] {
	raw, ok := headerLookup(extra, name, form)
	if !ok {
		return jmap.Optional[T]{}
	}
	var out jmap.Optional[T]
	if err := jsonv2.Unmarshal(raw, &out); err != nil {
		return jmap.Optional[T]{}
	}
	return out
}

// HeaderText returns the header:* asText value for name.
func (e *Email) HeaderText(name string) jmap.Optional[string] {
	return headerDecode[string](e.Extra, name, FormText)
}

// HeaderAddresses returns the header:* asAddresses value for name.
func (e *Email) HeaderAddresses(name string) jmap.Optional[[]*mail.Address] {
	return headerDecode[[]*mail.Address](e.Extra, name, FormAddresses)
}

// HeaderGroupedAddresses returns the header:* asGroupedAddresses value for name.
func (e *Email) HeaderGroupedAddresses(name string) jmap.Optional[[]*AddressGroup] {
	return headerDecode[[]*AddressGroup](e.Extra, name, FormGroupedAddresses)
}

// HeaderURLs returns the header:* asURLs value for name.
func (e *Email) HeaderURLs(name string) jmap.Optional[[]string] {
	return headerDecode[[]string](e.Extra, name, FormURLs)
}

// HeaderDate returns the header:* asDate value for name.
func (e *Email) HeaderDate(name string) jmap.Optional[time.Time] {
	return headerDecode[time.Time](e.Extra, name, FormDate)
}

// HeaderMessageIDs returns the header:* asMessageIds value for name.
func (e *Email) HeaderMessageIDs(name string) jmap.Optional[[]string] {
	return headerDecode[[]string](e.Extra, name, FormMessageIds)
}

// HeaderText returns the header:* asText value for name on a body part.
func (p *BodyPart) HeaderText(name string) jmap.Optional[string] {
	return headerDecode[string](p.Extra, name, FormText)
}

// HeaderAddresses returns the header:* asAddresses value for name on a body part.
func (p *BodyPart) HeaderAddresses(name string) jmap.Optional[[]*mail.Address] {
	return headerDecode[[]*mail.Address](p.Extra, name, FormAddresses)
}

// HeaderGroupedAddresses returns the header:* asGroupedAddresses value for name on a body part.
func (p *BodyPart) HeaderGroupedAddresses(name string) jmap.Optional[[]*AddressGroup] {
	return headerDecode[[]*AddressGroup](p.Extra, name, FormGroupedAddresses)
}

// HeaderURLs returns the header:* asURLs value for name on a body part.
func (p *BodyPart) HeaderURLs(name string) jmap.Optional[[]string] {
	return headerDecode[[]string](p.Extra, name, FormURLs)
}

// HeaderDate returns the header:* asDate value for name on a body part.
func (p *BodyPart) HeaderDate(name string) jmap.Optional[time.Time] {
	return headerDecode[time.Time](p.Extra, name, FormDate)
}

// HeaderMessageIDs returns the header:* asMessageIds value for name on a body part.
func (p *BodyPart) HeaderMessageIDs(name string) jmap.Optional[[]string] {
	return headerDecode[[]string](p.Extra, name, FormMessageIds)
}
