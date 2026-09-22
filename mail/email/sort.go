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

func Asc(prop string) *jmap.Comparator  { return jmap.Asc(prop) }
func Desc(prop string) *jmap.Comparator { return jmap.Desc(prop) }

func ByKeyword(kw string, ascending bool) *jmap.Comparator {
	return &jmap.Comparator{Property: SortHasKeyword, Keyword: kw, IsAscending: new(ascending)}
}
