package blob

import (
	"encoding/json"
	"strings"

	"git.sr.ht/~rockorager/go-jmap"
)

// Get binary blob data and metadata.
// https://www.rfc-editor.org/rfc/rfc9404.html#section-4.2
type Get struct {
	Account jmap.ID `json:"accountId,omitempty"`

	IDs []jmap.ID `json:"ids,omitempty"`

	Properties []string `json:"properties,omitempty"`

	Offset *uint64 `json:"offset,omitempty"`
	Length *uint64 `json:"length,omitempty"`
}

func (m *Get) Name() string { return "Blob/get" }

func (m *Get) Requires() []jmap.URI { return []jmap.URI{URI} }

type GetResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	List []*GetResult `json:"list,omitempty"`

	NotFound []jmap.ID `json:"notFound,omitempty"`
}

type GetResult struct {
	ID jmap.ID `json:"id,omitempty"`

	DataAsText   *string `json:"data:asText,omitempty"`
	DataAsBase64 *string `json:"data:asBase64,omitempty"`

	IsEncodingProblem bool `json:"isEncodingProblem,omitempty"`
	IsTruncated       bool `json:"isTruncated,omitempty"`

	Size uint64 `json:"size,omitempty"`

	Digests map[string]string `json:"-"`
}

func (r GetResult) MarshalJSON() ([]byte, error) {
	raw := map[string]interface{}{}
	if r.ID != "" {
		raw["id"] = r.ID
	}
	if r.DataAsText != nil {
		raw["data:asText"] = r.DataAsText
	}
	if r.DataAsBase64 != nil {
		raw["data:asBase64"] = r.DataAsBase64
	}
	if r.IsEncodingProblem {
		raw["isEncodingProblem"] = r.IsEncodingProblem
	}
	if r.IsTruncated {
		raw["isTruncated"] = r.IsTruncated
	}
	if r.Size != 0 {
		raw["size"] = r.Size
	}
	for algorithm, digest := range r.Digests {
		raw["digest:"+algorithm] = digest
	}
	return json.Marshal(raw)
}

func (r *GetResult) UnmarshalJSON(data []byte) error {
	*r = GetResult{}

	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key, value := range raw {
		switch {
		case key == "id":
			if err := json.Unmarshal(value, &r.ID); err != nil {
				return err
			}
		case key == "data:asText":
			if err := json.Unmarshal(value, &r.DataAsText); err != nil {
				return err
			}
		case key == "data:asBase64":
			if err := json.Unmarshal(value, &r.DataAsBase64); err != nil {
				return err
			}
		case key == "isEncodingProblem":
			if err := json.Unmarshal(value, &r.IsEncodingProblem); err != nil {
				return err
			}
		case key == "isTruncated":
			if err := json.Unmarshal(value, &r.IsTruncated); err != nil {
				return err
			}
		case key == "size":
			if err := json.Unmarshal(value, &r.Size); err != nil {
				return err
			}
		case strings.HasPrefix(key, "digest:"):
			if r.Digests == nil {
				r.Digests = map[string]string{}
			}
			var digest string
			if err := json.Unmarshal(value, &digest); err != nil {
				return err
			}
			r.Digests[strings.TrimPrefix(key, "digest:")] = digest
		}
	}
	return nil
}

func newGetResponse() jmap.MethodResponse { return &GetResponse{} }
