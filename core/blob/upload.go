package blob

import (
	"encoding/json"

	"github.com/Janso123/go-jmap"
)

// Upload binary blobs using JMAP method calls.
// https://www.rfc-editor.org/rfc/rfc9404.html#section-4.1
type Upload struct {
	Account jmap.ID `json:"accountId,omitempty"`

	Create map[jmap.ID]*UploadObject `json:"create,omitempty"`
}

func (m *Upload) Name() string { return "Blob/upload" }

func (m *Upload) Requires() []jmap.URI { return []jmap.URI{URI} }

type UploadObject struct {
	Data []DataSource `json:"data,omitempty"`

	Type string `json:"type,omitempty"`
}

// DataSource is a RFC 9404 DataSourceObject.
type DataSource struct {
	AsText   *string `json:"-"`
	AsBase64 *string `json:"-"`

	BlobID jmap.ID `json:"blobId,omitempty"`
	Offset *uint64 `json:"offset,omitempty"`
	Length *uint64 `json:"length,omitempty"`
}

func (s DataSource) MarshalJSON() ([]byte, error) {
	raw := map[string]interface{}{}
	if s.AsText != nil {
		raw["data:asText"] = s.AsText
	}
	if s.AsBase64 != nil {
		raw["data:asBase64"] = s.AsBase64
	}
	if s.BlobID != "" {
		raw["blobId"] = s.BlobID
	}
	if s.Offset != nil {
		raw["offset"] = s.Offset
	}
	if s.Length != nil {
		raw["length"] = s.Length
	}
	return json.Marshal(raw)
}

func (s *DataSource) UnmarshalJSON(data []byte) error {
	*s = DataSource{}

	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if value, ok := raw["data:asText"]; ok {
		if err := json.Unmarshal(value, &s.AsText); err != nil {
			return err
		}
	}
	if value, ok := raw["data:asBase64"]; ok {
		if err := json.Unmarshal(value, &s.AsBase64); err != nil {
			return err
		}
	}
	if value, ok := raw["blobId"]; ok {
		if err := json.Unmarshal(value, &s.BlobID); err != nil {
			return err
		}
	}
	if value, ok := raw["offset"]; ok {
		if err := json.Unmarshal(value, &s.Offset); err != nil {
			return err
		}
	}
	if value, ok := raw["length"]; ok {
		if err := json.Unmarshal(value, &s.Length); err != nil {
			return err
		}
	}
	return nil
}

type UploadResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	OldState string `json:"oldState,omitempty"`

	Created map[jmap.ID]*UploadBlob `json:"created,omitempty"`

	NotCreated map[jmap.ID]*jmap.SetError `json:"notCreated,omitempty"`
}

type UploadBlob struct {
	ID   jmap.ID `json:"id,omitempty"`
	Type string  `json:"type,omitempty"`
	Size uint64  `json:"size,omitempty"`
}

func newUploadResponse() jmap.MethodResponse { return &UploadResponse{} }
