package emailsubmission

import "github.com/Janso123/go-jmap"

// Email submission filter criteria
// https://www.rfc-editor.org/rfc/rfc8621.html#section-7.3
type FilterCondition struct {
	jmap.FilterBase `json:"-"`

	IdentityIDs []jmap.ID     `json:"identityIds,omitzero"`
	EmailIDs    []jmap.ID     `json:"emailIds,omitzero"`
	ThreadIDs   []jmap.ID     `json:"threadIds,omitzero"`
	UndoStatus  UndoStatus    `json:"undoStatus,omitzero"`
	Before      *jmap.UTCDate `json:"before,omitzero"`
	After       *jmap.UTCDate `json:"after,omitzero"`
}

func And(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.And(conds...) }
func Or(conds ...jmap.Filter) *jmap.FilterOperator  { return jmap.Or(conds...) }
func Not(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.Not(conds...) }
