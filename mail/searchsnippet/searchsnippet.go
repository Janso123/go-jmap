package searchsnippet

import "github.com/Janso123/go-jmap"

func init() {
	jmap.RegisterMethod("SearchSnippet/get", newGetResponse)
}

// Search preview snippet
// https://www.rfc-editor.org/rfc/rfc8621.html#section-5
type SearchSnippet struct {
	Email jmap.ID `json:"emailId,omitzero"`

	Subject *string `json:"subject,omitzero"`

	Preview *string `json:"preview,omitzero"`
}
