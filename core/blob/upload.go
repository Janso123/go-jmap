package blob

import (
	"fmt"

	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

// Upload binary blobs using JMAP method calls.
// https://www.rfc-editor.org/rfc/rfc9404.html#section-4.1
type Upload struct {
	Account jmap.ID `json:"accountId,omitzero"`

	Create map[jmap.ID]*UploadObject `json:"create,omitzero"`
}

func (m *Upload) Name() string { return "Blob/upload" }

func (m *Upload) Requires() []jmap.URI { return []jmap.URI{URI} }

type UploadObject struct {
	Data []DataSource `json:"data"`

	Type jmap.Optional[string] `json:"type,omitzero"`
}

// DataSource is a RFC 9404 DataSourceObject.
// Marshal order is data:asText, data:asBase64, blobId, offset, length.
type DataSource struct {
	AsText   jmap.Optional[string] `json:"-"`
	AsBase64 jmap.Optional[string] `json:"-"`

	BlobID jmap.ID           `json:"blobId,omitzero"`
	Offset *jmap.UnsignedInt `json:"offset,omitzero"`
	Length *jmap.UnsignedInt `json:"length,omitzero"`
}

func (s DataSource) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	if err := writeOptionalString(enc, "data:asText", s.AsText); err != nil {
		return err
	}
	if err := writeOptionalString(enc, "data:asBase64", s.AsBase64); err != nil {
		return err
	}
	if err := writeID(enc, "blobId", s.BlobID); err != nil {
		return err
	}
	if err := writeUintPtr(enc, "offset", s.Offset); err != nil {
		return err
	}
	if err := writeUintPtr(enc, "length", s.Length); err != nil {
		return err
	}
	return enc.WriteToken(jsontext.EndObject)
}

func (s *DataSource) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	*s = DataSource{}
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '{' {
		return fmt.Errorf("blob: DataSource: expected JSON object, got %v", tok.Kind())
	}
	for dec.PeekKind() != '}' {
		keyTok, err := dec.ReadToken()
		if err != nil {
			return err
		}
		switch keyTok.String() {
		case "data:asText":
			err = jsonv2.UnmarshalDecode(dec, &s.AsText)
		case "data:asBase64":
			err = jsonv2.UnmarshalDecode(dec, &s.AsBase64)
		case "blobId":
			err = jsonv2.UnmarshalDecode(dec, &s.BlobID)
		case "offset":
			var n *jmap.UnsignedInt
			err = jsonv2.UnmarshalDecode(dec, &n)
			s.Offset = n
		case "length":
			var n *jmap.UnsignedInt
			err = jsonv2.UnmarshalDecode(dec, &n)
			s.Length = n
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

type UploadResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	OldState string `json:"oldState,omitzero"`

	Created map[jmap.ID]*UploadBlob `json:"created,omitzero"`

	NotCreated map[jmap.ID]*jmap.SetError `json:"notCreated,omitzero"`
}

type UploadBlob struct {
	ID   jmap.ID               `json:"id,omitzero"`
	Type jmap.Optional[string] `json:"type,omitzero"`
	Size jmap.UnsignedInt      `json:"size"`
}

func newUploadResponse() jmap.MethodResponse { return &UploadResponse{} }
