package blob

import (
	"fmt"
	"sort"
	"strings"

	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

// Get binary blob data and metadata.
// https://www.rfc-editor.org/rfc/rfc9404.html#section-4.2
type Get struct {
	Account jmap.ID `json:"accountId,omitzero"`

	IDs jmap.Optional[[]jmap.ID] `json:"ids,omitzero"`

	Properties []string `json:"properties,omitzero"`

	Offset *jmap.UnsignedInt `json:"offset,omitzero"`
	Length *jmap.UnsignedInt `json:"length,omitzero"`

	ReferenceIDs        *jmap.ResultReference `json:"#ids,omitzero"`
	ReferenceProperties *jmap.ResultReference `json:"#properties,omitzero"`
}

func (m *Get) Name() string { return "Blob/get" }

func (m *Get) Requires() []jmap.URI { return []jmap.URI{URI} }

type GetResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	List []*GetResult `json:"list,omitzero"`

	NotFound []jmap.ID `json:"notFound,omitzero"`
}

type GetResult struct {
	ID jmap.ID `json:"id,omitzero"`

	DataAsText   jmap.Optional[string] `json:"data:asText,omitzero"`
	DataAsBase64 jmap.Optional[string] `json:"data:asBase64,omitzero"`

	// Booleans stay on the wire when false.
	IsEncodingProblem bool `json:"isEncodingProblem"`
	IsTruncated       bool `json:"isTruncated"`

	// Size is server-set. Nil omits the key; a pointer to 0 stays "size":0.
	Size *jmap.UnsignedInt `json:"size,omitzero"`

	Digests map[string]string `json:"-"`
}

// MarshalJSONTo writes fields in a fixed order: id, data:asText, data:asBase64,
// digest:* sorted by algorithm, size, isEncodingProblem, isTruncated.
func (r GetResult) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	if err := writeID(enc, "id", r.ID); err != nil {
		return err
	}
	if err := writeOptionalString(enc, "data:asText", r.DataAsText); err != nil {
		return err
	}
	if err := writeOptionalString(enc, "data:asBase64", r.DataAsBase64); err != nil {
		return err
	}
	if err := writeDigests(enc, r.Digests); err != nil {
		return err
	}
	if err := writeUintPtr(enc, "size", r.Size); err != nil {
		return err
	}
	if err := writeBool(enc, "isEncodingProblem", r.IsEncodingProblem); err != nil {
		return err
	}
	if err := writeBool(enc, "isTruncated", r.IsTruncated); err != nil {
		return err
	}
	return enc.WriteToken(jsontext.EndObject)
}

func (r *GetResult) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	*r = GetResult{}
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '{' {
		return fmt.Errorf("blob: GetResult: expected JSON object, got %v", tok.Kind())
	}
	for dec.PeekKind() != '}' {
		keyTok, err := dec.ReadToken()
		if err != nil {
			return err
		}
		key := keyTok.String()
		switch {
		case key == "id":
			err = jsonv2.UnmarshalDecode(dec, &r.ID)
		case key == "data:asText":
			err = jsonv2.UnmarshalDecode(dec, &r.DataAsText)
		case key == "data:asBase64":
			err = jsonv2.UnmarshalDecode(dec, &r.DataAsBase64)
		case key == "isEncodingProblem":
			err = jsonv2.UnmarshalDecode(dec, &r.IsEncodingProblem)
		case key == "isTruncated":
			err = jsonv2.UnmarshalDecode(dec, &r.IsTruncated)
		case key == "size":
			var n *jmap.UnsignedInt
			err = jsonv2.UnmarshalDecode(dec, &n)
			r.Size = n
		case strings.HasPrefix(key, "digest:"):
			if dec.PeekKind() == 'n' {
				_, err = dec.ReadValue()
				break
			}
			if r.Digests == nil {
				r.Digests = map[string]string{}
			}
			var digest string
			if err = jsonv2.UnmarshalDecode(dec, &digest); err != nil {
				return err
			}
			r.Digests[strings.TrimPrefix(key, "digest:")] = digest
		default:
			_, err = dec.ReadValue()
		}
		if err != nil {
			return err
		}
	}
	_, err = dec.ReadToken()
	return err
}

func writeID(enc *jsontext.Encoder, key string, id jmap.ID) error {
	if id == "" {
		return nil
	}
	if err := enc.WriteToken(jsontext.String(key)); err != nil {
		return err
	}
	return id.MarshalJSONTo(enc)
}

func writeOptionalString(enc *jsontext.Encoder, key string, v jmap.Optional[string]) error {
	if v.IsZero() {
		return nil
	}
	if err := enc.WriteToken(jsontext.String(key)); err != nil {
		return err
	}
	return v.MarshalJSONTo(enc)
}

func writeStringField(enc *jsontext.Encoder, key, value string) error {
	if err := enc.WriteToken(jsontext.String(key)); err != nil {
		return err
	}
	return enc.WriteToken(jsontext.String(value))
}

func writeUintPtr(enc *jsontext.Encoder, key string, v *jmap.UnsignedInt) error {
	if v == nil {
		return nil
	}
	if err := enc.WriteToken(jsontext.String(key)); err != nil {
		return err
	}
	return v.MarshalJSONTo(enc)
}

func writeBool(enc *jsontext.Encoder, key string, v bool) error {
	if err := enc.WriteToken(jsontext.String(key)); err != nil {
		return err
	}
	return enc.WriteToken(jsontext.Bool(v))
}

func writeDigests(enc *jsontext.Encoder, digests map[string]string) error {
	if len(digests) == 0 {
		return nil
	}
	keys := make([]string, 0, len(digests))
	for algo := range digests {
		keys = append(keys, algo)
	}
	sort.Strings(keys)
	for _, algo := range keys {
		if err := writeStringField(enc, "digest:"+algo, digests[algo]); err != nil {
			return err
		}
	}
	return nil
}

func newGetResponse() jmap.MethodResponse { return &GetResponse{} }
