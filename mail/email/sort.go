package email

import "github.com/Janso123/go-jmap"

// Email sort property constants (RFC 8621 §4.4.2).
const (
	SortReceivedAt              = "receivedAt"
	SortSize                    = "size"
	SortFrom                    = "from"
	SortTo                      = "to"
	SortSubject                 = "subject"
	SortSentAt                  = "sentAt"
	SortHasKeyword              = "hasKeyword"
	SortAllInThreadHaveKeyword  = "allInThreadHaveKeyword"
	SortSomeInThreadHaveKeyword = "someInThreadHaveKeyword"
)

// Email sort criteria
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.4.2
type SortComparator struct {
	Property string `json:"property,omitzero"`

	Keyword string `json:"keyword,omitzero"`

	// IsAscending has no omitzero: RFC default is true, so Desc must emit false.
	IsAscending bool `json:"isAscending"`

	Collation jmap.CollationAlgo `json:"collation,omitzero"`
}

// Asc returns a SortComparator sorted ascending on prop.
func Asc(prop string) *SortComparator {
	return &SortComparator{Property: prop, IsAscending: true}
}

// Desc returns a SortComparator sorted descending on prop.
func Desc(prop string) *SortComparator {
	return &SortComparator{Property: prop, IsAscending: false}
}

// ByKeyword returns a SortComparator for hasKeyword / thread-keyword sorts.
func ByKeyword(kw string, ascending bool) *SortComparator {
	return &SortComparator{Property: SortHasKeyword, Keyword: kw, IsAscending: ascending}
}
